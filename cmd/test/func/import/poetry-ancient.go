package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/easysoft/zendata/cmd/test/func/comm"
	"github.com/easysoft/zendata/cmd/test/func/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"github.com/snowlyg/helper/dir"
)

func main() {
	filePath := "/Users/aaron/work/zentao/product/zd/行业数据/古诗词(分类)"

	var files []string
	filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	tableName := model.PoetryAncient{}.TableName()
	db := comm.GetDB()
	err := db.Exec(fmt.Sprintf(comm.TruncateTable, tableName)).Error
	if err != nil {
		panic(err)
	}

	tableName2 := model.PoetryCategory{}.TableName()
	err = db.Exec(fmt.Sprintf(comm.TruncateTable, tableName2)).Error
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

			name := strings.TrimRight(filepath.Base(file), ".txt")
			cate := model.PoetryCategory{
				EnName: name,
				Name:   name,
			}

			cateOld := model.PoetryCategory{}
			db.Where("en_name = ?", name).First(&cateOld)
			id := cateOld.Id
			if id == 0 {
				db.Save(&cate)
				id = cate.Id
			}

			po := model.PoetryAncient{
				CategoryId: id,
				Content:    col1,
			}

			if po.Content != "" {
				db.Save(&po)
			}
		}
	}
}
