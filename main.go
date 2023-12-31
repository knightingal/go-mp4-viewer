package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var db *sql.DB

func main() {
	fmt.Println("hello")

	initDB()
	router := gin.Default()
	router.GET("/hello", helloHanlder)
	router.GET("/mp4-dir/:baseIndex/*subDir", mp4DirHanlder)
	router.GET("/video-info/:baseIndex/*subDir", videoInfoHandler)
	router.GET("/init-video/:baseIndex/*subDir", initVideoInfoHandler)
	router.GET("/mount-config", mountConfigHanlder)

	s8082 := &http.Server{
		Addr:    ":8082",
		Handler: router,
	}

	s8082.ListenAndServe()
}

func helloHanlder(context *gin.Context) {
	context.String(http.StatusOK, "hellp")
}

func listVideoFile(subDir string, indexNumber int) []string {
	rows, err := db.Query("select dir_path from mp4_base_dir where id=?", indexNumber)
	var dirList []string
	if strings.EqualFold(subDir, "/") {
		if err != nil {
			log.Fatal(err)
		} else {
			var baseDir string
			rows.Next()
			rows.Scan(&baseDir)
			rows.Close()
			dirList = scanBaseDir(baseDir)
		}
	} else {
		var baseDir string
		rows.Next()
		rows.Scan(&baseDir)
		rows.Close()
		dirList = scanFileInDir(baseDir + subDir)
	}
	return dirList
}

func initVideoInfoHandler(context *gin.Context) {
	subDir := context.Param("subDir")
	baseIndex := context.Param("baseIndex")
	indexNumber, _ := strconv.Atoi(baseIndex)
	dirList := listVideoFile(subDir, indexNumber)
	videoCoverList := parseVideoCover(dirList)
	for _, videoCover := range videoCoverList {
		result, error := db.Exec("insert into video_info("+
			"dir_path, base_index, video_file_name, cover_file_name) values (?,?,?,?)",
			subDir, indexNumber, videoCover.videoFileName, videoCover.coverFileName)
		if error != nil {
			log.Fatal(error)
		}
		insertId, _ := result.LastInsertId()
		fmt.Printf("insert %d\n", insertId)
	}
	context.String(http.StatusOK, "succ")
}

type VideoCover struct {
	videoFileName string
	coverFileName string
}

func filter[T any](src *[]T, fn func(T) bool) *[]T {
	ret := make([]T, 0)
	for _, item := range *src {
		if fn(item) {
			ret = append(ret, item)
		}
	}
	return &ret
}

func parseVideoCover(dirList []string) []VideoCover {
	videoCoverList := make([]VideoCover, 0)

	videoFileNameList := make([]string, 0)
	imgFileNameList := make([]string, 0)
	for _, fileName := range dirList {
		if strings.HasSuffix(fileName, ".mp4") {
			videoFileNameList = append(videoFileNameList, fileName)
		} else if strings.HasSuffix(fileName, ".jpg") || strings.HasSuffix(fileName, ".png") {
			imgFileNameList = append(imgFileNameList, fileName)
		}
	}

	for _, videoFileName := range videoFileNameList {
		videoCover, succ := videoMatchToCover(videoFileName, imgFileNameList)
		if succ {
			videoCoverList = append(videoCoverList, videoCover)
		}
	}

	return videoCoverList
}

func videoMatchToCover(videoFileName string, imgFileNameList []string) (VideoCover, bool) {
	match := func(src string) (string, bool) {

		filterRet := filter(&imgFileNameList, func(dirName string) bool {
			return strings.Contains(dirName, src)
		})

		if len(*filterRet) == 1 {
			fmt.Println("====matched====")
			fmt.Println((*filterRet)[0])
			return (*filterRet)[0], true
		}

		return "", false
	}

	pureName := strings.Split(videoFileName, ".")[0]
	srcArray := []rune(pureName)
	size := len(srcArray)
	for i := 0; i < size; i++ {
		for j := 0; j <= i; j++ {
			sub1 := srcArray[j : j+size-i]
			fmt.Println(string(sub1))
			realName, matched := match(string(sub1))
			if matched {
				return VideoCover{videoFileName: videoFileName, coverFileName: realName}, true
			}
		}
	}
	return VideoCover{}, false
}

func videoInfoHandler(context *gin.Context) {
	subDir := context.Param("subDir")
	baseIndex := context.Param("baseIndex")
	indexNumber, _ := strconv.Atoi(baseIndex)
	rows, err := db.Query("select id, video_file_name, cover_file_name from video_info where dir_path = ? and base_index=?", subDir, indexNumber)
	if err != nil {
		context.Header("Access-Control-Allow-Origin", "*")
		context.JSONP(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	data := make([]any, 0)
	for rows.Next() {
		var videoFileName string
		var coverFileName string
		var id int
		rows.Scan(&id, &videoFileName, &coverFileName)
		data = append(data, map[string]interface{}{
			"id":            id,
			"videoFileName": videoFileName,
			"coverFileName": coverFileName,
		})
	}
	rows.Close()
	context.Header("Access-Control-Allow-Origin", "*")
	context.JSONP(http.StatusOK, data)
}

func mountConfigHanlder(context *gin.Context) {
	rows, err := db.Query("select id, dir_path, url_prefix from mp4_base_dir ")
	if err != nil {
		context.Header("Access-Control-Allow-Origin", "*")
		context.String(http.StatusInternalServerError, err.Error())
	}

	data := make([]any, 0)
	for rows.Next() {
		var baseDir string
		var id int
		var urlPrefix string
		rows.Scan(&id, &baseDir, &urlPrefix)
		data = append(data, map[string]interface{}{
			"id":        id,
			"baseDir":   baseDir,
			"urlPrefix": urlPrefix,
		})
	}
	rows.Close()
	context.Header("Access-Control-Allow-Origin", "*")
	context.JSONP(http.StatusOK, data)
}

func mp4DirHanlder(context *gin.Context) {

	subDir := context.Param("subDir")
	baseIndex := context.Param("baseIndex")
	indexNumber, _ := strconv.Atoi(baseIndex)
	//	data2 := [3]any{"hello", "world", map[string]interface{}{
	//		"year1": 2023,
	//		"year2": 2024,
	//	}}
	dirList := listVideoFile(subDir, indexNumber)
	context.Header("Access-Control-Allow-Origin", "*")
	context.JSONP(http.StatusOK, dirList)
}
