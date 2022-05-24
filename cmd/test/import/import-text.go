package main

import (
	"flag"
	"fmt"
	"github.com/easysoft/zendata/internal/pkg/test/import/comm"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"strings"
)

var (
	tableName string
	filePath  string
	colNum    int
)

func main() {
	flag.StringVar(&tableName, "t", "", "")
	flag.StringVar(&filePath, "f", "", "")
	flag.IntVar(&colNum, "c", 0, "")

	flag.Parse()

	tableName = "biz_" + strings.TrimLeft(tableName, "biz_")
	db := comm.GetDB()
	db.Exec(fmt.Sprintf(comm.TruncateTable, tableName))

	createTableSql := fmt.Sprintf(comm.CreateTableTempl, tableName)
	err := db.Exec(createTableSql).Error
	if err != nil {
		fmt.Printf("create table %s failed, err %s", tableName, err.Error())
		return
	}

	content := fileUtils.ReadFile(filePath)
	insertSqlArr := make([]string, 0)

	for _, line := range strings.Split(content, "\n") {
		arr := strings.Split(strings.TrimSpace(line), " ")

		if colNum >= len(arr) {
			continue
		}

		content := strings.TrimSpace(arr[colNum])
		if content == "" {
			continue
		}

		insert := fmt.Sprintf("('%s')", content)
		insertSqlArr = append(insertSqlArr, insert)
	}

	for i := 0; i < 1000; i++ {
		start := i * 10000
		end := (i + 1) * 10000

		if end > len(insertSqlArr) {
			end = len(insertSqlArr)
		}

		arr := insertSqlArr[start:end]

		sql := fmt.Sprintf(comm.InsertTemplate, tableName, strings.Join(arr, ","))
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
