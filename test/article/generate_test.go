package main

import (
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func TestGenerate(t *testing.T) {
	files := make([]string, 0)
	getFilesInDir("xdoc/words-9.3", ".txt", &files)

	for _, filePath := range files {
		article := fileUtils.ReadFile(filePath)
		content := convertToYaml(article)

		newPath := changeFileExt(filePath, ".yaml")
		fileUtils.WriteFile(newPath, content)
	}
}

func convertToYaml(article string) (content string) {


	return
}