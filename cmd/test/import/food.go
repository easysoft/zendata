package main

import (
	"fmt"
	"github.com/easysoft/zendata/cmd/test/comm"
	"github.com/easysoft/zendata/cmd/test/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"strings"
)

func main() {
	tableName := "food"
	filePath := "/Users/aaron/work/zentao/product/zd/行业数据/食物.txt"

	tableName = "biz_data_" + tableName
	db := comm.GetDB()
	err := db.Exec(fmt.Sprintf(comm.TruncateTable, tableName)).Error
	if err != nil {
		panic(err)
	}

	content := fileUtils.ReadFile(filePath)

	for _, line := range strings.Split(content, "\n") {
		arr := strings.Split(strings.TrimSpace(line), " ")
		col1 := arr[0]

		po := model.DataFood{
			DataComm: model.DataComm{
				Name: col1,
			},
		}

		db.Save(&po)
	}
}
