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

	s8082 := &http.Server{
		Addr:    ":8082",
		Handler: router,
	}

	s8082.ListenAndServe()
}

func helloHanlder(context *gin.Context) {
	context.String(http.StatusOK, "hellp")
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
			dirList = make([]string, 0)
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
