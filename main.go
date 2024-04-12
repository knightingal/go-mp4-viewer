package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

var db *sql.DB

func main() {
	fmt.Println("hello")

	initDB()
	router := gin.Default()
	router.GET("/hello/:fileName/:time", helloHanlder)
	router.GET("/mp4-dir/:baseIndex/*subDir", mp4DirHanlder)
	router.GET("/video-info/:baseIndex/*subDir", videoInfoHandler)
	router.POST("/video-info/:videoId", postVideoInfoHandler)
	router.GET("/init-video/:baseIndex/*subDir", initVideoInfoHandler)
	router.GET("/mount-config", mountConfigHanlder)

	s8082 := &http.Server{
		Addr:    ":8082",
		Handler: router,
	}

	s8082.ListenAndServe()
}

var imgFile = "/tmp/mp4fifo"

func helloHanlder(context *gin.Context) {
	fileName := context.Param("fileName")
	timeStamp := context.Param("time")

	var wg sync.WaitGroup
	wg.Add(1)

	go context.Stream(func(w io.Writer) bool {
		defer wg.Done()
		f, err := os.Open(imgFile)
		if err != nil {
			w.Write([]byte(err.Error()))
			return false
		}
		buf := make([]byte, 1024)

		for {
			readLen, _ := f.Read(buf)
			if readLen <= 0 {
				break
			}

			w.Write(buf[:readLen])
		}
		return false
	})
	log.Default().Println(fileName)
	cmd := exec.Command("/usr/bin/ffmpeg", "-ss", timeStamp, "-i",
		"/home/knightingal/"+fileName, "-frames:v", "1", "-f", "image2", "-y",
		"/tmp/mp4fifo")
	cmd.Run()
	wg.Wait()
}

func queryBaseDirByDirIndex(indexNumber int) (string, error) {
	var baseDir string
	rows, err := db.Query("select dir_path from mp4_base_dir where id=?", indexNumber)
	if err != nil {
		return baseDir, err
	}

	rows.Next()
	rows.Scan(&baseDir)
	rows.Close()
	return baseDir, nil
}

func listVideoFile(subDir string, indexNumber int) []string {
	baseDir, err := queryBaseDirByDirIndex(indexNumber)
	var dirList []string
	if err != nil {
		log.Fatal(err)
	}
	if strings.EqualFold(subDir, "/") {
		dirList = scanBaseDir(baseDir)
	} else {
		dirList = scanFileInDir(baseDir + subDir)
	}
	return dirList
}

func initVideoInfoHandler(context *gin.Context) {
	subDir := context.Param("subDir")
	baseIndex := context.Param("baseIndex")
	forceStr := context.Query("force")
	force, error := strconv.ParseBool(forceStr)
	if error != nil {
		force = false
	}
	indexNumber, _ := strconv.Atoi(baseIndex)
	if force {
		db.Exec("delete from video_info where dir_path=? and base_index=?", subDir, indexNumber)
	}
	dirList := listVideoFile(subDir, indexNumber)
	videoCoverList, missMatchedList := parseVideoCover(dirList)
	for _, videoCover := range videoCoverList {
		var existCount int
		rows, error := db.Query("select count(video_file_name) from video_info where dir_path=? and base_index=? and video_file_name=?",
			subDir, indexNumber, videoCover.videoFileName)

		if error != nil {
			log.Fatal(error)
		}
		rows.Next()
		rows.Scan(&existCount)
		rows.Close()
		if existCount > 0 {
			log.Printf("%s exist, skip insert", videoCover.videoFileName)
			continue
		}

		result, error := db.Exec("insert into video_info("+
			"dir_path, base_index, video_file_name, cover_file_name) values (?,?,?,?)",
			subDir, indexNumber, videoCover.videoFileName, videoCover.coverFileName)
		if error != nil {
			log.Fatal(error)
		}
		insertId, _ := result.LastInsertId()
		fmt.Printf("insert %d\n", insertId)
	}
	for _, unMathed := range missMatchedList {
		log.Printf("%s not matched", unMathed)
		result, error := db.Exec("insert into miss_match_video_record("+
			"dir_path, base_index, video_file_name) values (?,?,?)",
			subDir, indexNumber, unMathed)
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

func parseVideoCover(dirList []string) ([]VideoCover, []string) {
	videoCoverList := make([]VideoCover, 0)
	missMatched := make([]string, 0)

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
		} else {
			missMatched = append(missMatched, videoFileName)
		}
	}

	return videoCoverList, missMatched
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

	parsePureNameFromFileName := func(fileName string) string {
		lastIndex := strings.LastIndex(fileName, ".")
		pureName := fileName[0:lastIndex]
		return pureName
	}

	pureName := parsePureNameFromFileName(videoFileName)
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
	subDir = strings.TrimSuffix(subDir, "/")
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

func postVideoInfoHandler(context *gin.Context) {
	baseIndex := context.Param("videoId")
	videoId, _ := strconv.Atoi(baseIndex)
	log.Default().Println(videoId)
	json := make(map[string]interface{})
	context.BindJSON(&json)
	rate, _ := json["rate"].(float64)
	log.Default().Println(rate)
	db.Exec("update video_info set rate = ? where id = ?", rate, videoId)

	context.Header("Access-Control-Allow-Origin", "*")
	context.Status(http.StatusOK)
}

func mountConfigHanlder(context *gin.Context) {
	rows, err := db.Query("select id, dir_path, url_prefix, api_version from mp4_base_dir ")
	if err != nil {
		context.Header("Access-Control-Allow-Origin", "*")
		context.String(http.StatusInternalServerError, err.Error())
	}

	data := make([]any, 0)
	for rows.Next() {
		var baseDir string
		var id int
		var urlPrefix string
		var apiVersion int
		rows.Scan(&id, &baseDir, &urlPrefix, &apiVersion)
		data = append(data, map[string]interface{}{
			"apiVersion": apiVersion,
			"id":         id,
			"baseDir":    baseDir,
			"urlPrefix":  urlPrefix,
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
