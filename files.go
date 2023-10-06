package main

import (
	"io/fs"
	"os"
	"strings"
)

func scanBaseDir(baseDir string) []string {
	dir := os.DirFS(baseDir)
	dirEntityList, _ := fs.ReadDir(dir, ".")
	fileNames := make([]string, 0)
	for _, file := range dirEntityList {
		if file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}
	return fileNames
}

func scanFileInDir(baseDir string) []string {
	dir := os.DirFS(baseDir)
	dirEntityList, _ := fs.ReadDir(dir, ".")
	fileNames := make([]string, 0)
	for _, file := range dirEntityList {
		if !strings.HasSuffix(file.Name(), ".torrent") {
			fileNames = append(fileNames, file.Name())
		}
	}
	return fileNames
}
