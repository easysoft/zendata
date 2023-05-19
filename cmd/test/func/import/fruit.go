package main

import (
	"fmt"
	"strings"

	"github.com/easysoft/zendata/cmd/test/func/comm"
	"github.com/easysoft/zendata/cmd/test/func/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
)

func main() {
	tableName := "fruit"
	filePath := "/Users/aaron/work/zentao/product/zd/行业数据/水果.txt"

	tableName = "biz_data_" + tableName
	db := comm.GetDB()
	err := db.Exec(fmt.Sprintf(comm.TruncateTable, tableName)).Error
	if err != nil {
		panic(err)
	}

	content := fileUtils.ReadFile(filePath)

	for _, line := range strings.Split(content, "\n") {
		col1 := strings.Split(strings.TrimSpace(line), " ")[0]

		po := model.DataFruit{
			DataComm: model.DataComm{
				Name: col1,
			},
		}

		db.Save(&po)
	}
}
