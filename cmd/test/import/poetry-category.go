package main

import (
	"fmt"
	"github.com/easysoft/zendata/cmd/test/comm"
	"github.com/easysoft/zendata/cmd/test/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"github.com/snowlyg/helper/dir"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	filePath := "/Users/aaron/work/zentao/product/zd/行业数据/古诗词(分类)"

	var files []string
	filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	tableName := model.PoetryCategory{}.TableName()
	db := comm.GetDB()
	err := db.Exec(fmt.Sprintf(comm.TruncateTable, tableName)).Error
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if !dir.IsFile(file) || strings.Index(file, ".txt") < 0 {
			continue
		}

		content := fileUtils.ReadFile(file)

		for _, line := range strings.Split(content, "\n") {
			col1 := strings.TrimSpace(line)

			po := model.PoetryCategory{
				Name:    strings.TrimRight(filepath.Base(file), ".txt"),
				Content: col1,
			}

			if po.Content != "" {
				db.Save(&po)
			}
		}
	}
}
