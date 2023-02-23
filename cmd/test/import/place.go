package main

import (
	"fmt"
	"github.com/easysoft/zendata/cmd/test/comm"
	"github.com/easysoft/zendata/cmd/test/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"strings"
)

func main() {
	filePath := "/Users/aaron/work/zentao/product/zd/行业数据/场所.txt"

	tableName := model.Place{}.TableName()

	db := comm.GetDB()
	err := db.Exec(fmt.Sprintf(comm.TruncateTable, tableName)).Error
	if err != nil {
		panic(err)
	}

	content := fileUtils.ReadFile(filePath)

	for _, line := range strings.Split(content, "\r") {
		arr := strings.Split(strings.TrimSpace(line), " ")
		col1 := arr[0]

		po := model.Place{
			Name: col1,
		}

		if po.Name != "" {
			db.Save(&po)
		}
	}
}
