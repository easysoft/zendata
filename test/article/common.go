package main

import (
	"github.com/Chain-Zhang/pinyin"
	commonUtils "github.com/easysoft/zendata/src/utils/common"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
)

func getFilesInDir(folder, ext string, files *[]string) {
	 folder, _ = filepath.Abs(folder)

	if !fileUtils.IsDir(folder) {
		if path.Ext(folder) == ext {
			*files = append(*files, folder)
		}

		return
	}

	dir, err := ioutil.ReadDir(folder)
	if err != nil {
		return
	}

	for _, fi := range dir {
		name := fi.Name()
		if commonUtils.IngoreFile(name) {
			continue
		}

		filePath := fileUtils.AddSepIfNeeded(folder) + name
		if fi.IsDir() {
			getFilesInDir(filePath, ext, files)
		} else if strings.Index(name, "~") != 0 && path.Ext(filePath) == ext {
			*files = append(*files, filePath)
		}
	}
}

func getFileName(filePath string) string {
	fileName := path.Base(filePath)
	fileName = strings.TrimSuffix(fileName, path.Ext(filePath))

	return fileName
}

func changeFileExt(filePath, ext string) string {
	ret := strings.TrimSuffix(filePath, path.Ext(filePath))
	ret += ext

	return ret
}

func getPinyin(word string) string {
	p, _ := pinyin.New(word).Split("").Mode(pinyin.WithoutTone).Convert()

	return p
}