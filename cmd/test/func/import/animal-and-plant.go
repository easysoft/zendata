package main

import (
	"fmt"
	"strings"

	"github.com/easysoft/zendata/cmd/test/func/comm"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
)

func main() {
	var insertTemplate = "INSERT INTO %s (name) VALUES %s;"

	tableName := "animal_plant"
	filePath := "/Users/aaron/work/zentao/product/zd/行业数据/动植物.txt"

	tableName = "biz_data_" + tableName
	db := comm.GetDB()
	err := db.Exec(fmt.Sprintf(comm.TruncateTable, tableName)).Error
	if err != nil {
		panic(err)
	}

	content := fileUtils.ReadFile(filePath)

	insertSqlArr := make([]string, 0)

	arr := strings.Split(content, "\n")
	for _, line := range arr {
		col1 := strings.Split(strings.TrimSpace(line), " ")[0]
		col1 = strings.TrimSpace(col1)

		insert := fmt.Sprintf("('%s')", col1)
		insertSqlArr = append(insertSqlArr, insert)
	}

	for i := 0; i < 1000; i++ {
		start := i * 10000
		end := (i + 1) * 10000

		if end > len(insertSqlArr) {
			end = len(insertSqlArr)
		}

		arr := insertSqlArr[start:end]

		sql := fmt.Sprintf(insertTemplate, tableName, strings.Join(arr, ","))
		err = db.Exec(sql).Error
		if err != nil {
			fmt.Printf("insert data failed, err %s", err.Error())
			return
		}

		if end >= len(insertSqlArr) {
			break
		}
	}
}
