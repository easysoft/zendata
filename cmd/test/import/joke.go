package main

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zendata/cmd/test/comm"
	"github.com/easysoft/zendata/cmd/test/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"strings"
)

func main() {
	filePath := "/Users/aaron/work/zentao/product/zd/行业数据/笑话.json"

	tableName := model.Joke{}.TableName()
	db := comm.GetDB()
	err := db.Exec(fmt.Sprintf(comm.TruncateTable, tableName)).Error
	if err != nil {
		panic(err)
	}

	content := fileUtils.ReadFileBuf(filePath)

	arr := make([]string, 0)

	json.Unmarshal(content, &arr)

	for _, item := range arr {
		item = strings.TrimSpace(item)

		cate := model.Joke{
			Content: item,
		}

		if cate.Content != "" {
			db.Save(&cate)
		}
	}
}
