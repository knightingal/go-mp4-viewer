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
	var dirList []string

	fmt.Println(subDir)
	rows, err := db.Query("select dir_path from mp4_base_dir where id=?", indexNumber)
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
	context.Header("Access-Control-Allow-Origin", "*")
	context.JSONP(http.StatusOK, dirList)
}
