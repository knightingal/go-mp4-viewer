package main

import (
	"io/fs"
	"os"
	"sort"
	"strings"
)

func scanBaseDir(baseDir string) []string {
	dir := os.DirFS(baseDir)
	dirEntityList, _ := fs.ReadDir(dir, ".")
	fileNames := make([]string, 0)
	dirFileEntiryList := make([]fs.DirEntry, 0)

	for _, file := range dirEntityList {
		if file.IsDir() {
			dirFileEntiryList = append(dirFileEntiryList, file)
		}
	}

	sort.Slice(dirFileEntiryList, func(i, j int) bool {
		infoI, _ := dirFileEntiryList[i].Info()
		infoJ, _ := dirFileEntiryList[j].Info()
		return infoI.ModTime().Compare(infoJ.ModTime()) > 0
	})

	for _, file := range dirFileEntiryList {
		fileNames = append(fileNames, file.Name())
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
