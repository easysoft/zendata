package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/easysoft/zendata/cmd/test/import/comm"
	"github.com/easysoft/zendata/cmd/test/import/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"strings"
)

func main() {
	var insertTemplate = "INSERT INTO %s (riddle, answer) VALUES %s;"
	var createTableTempl = `CREATE TABLE IF NOT EXISTS %s (
		id bigint auto_increment,
		riddle varchar(1000) not null,
		answer varchar(1000) not null,
		tag varchar(50),
		primary key(id)
	) engine=innodb default charset=utf8 auto_increment=1;`

	var tableName string
	var filePath string

	flag.StringVar(&tableName, "t", "", "")
	flag.StringVar(&filePath, "f", "", "")

	flag.Parse()

	tableName = "biz_" + strings.TrimLeft(tableName, "biz_")
	db := comm.GetDB()
	db.Exec(fmt.Sprintf(comm.TruncateTable, tableName))

	createTableSql := fmt.Sprintf(createTableTempl, tableName)
	err := db.Exec(createTableSql).Error
	if err != nil {
		fmt.Printf("create table %s failed, err %s", tableName, err.Error())
		return
	}

	content := fileUtils.ReadFileBuf(filePath)
	records := make([]model.DataXiehouyu, 0)
	json.Unmarshal(content, &records)

	insertSqlArr := make([]string, 0)

	for _, record := range records {
		riddle := record.Riddle
		answer := record.Answer

		insert := fmt.Sprintf("('%s', '%s')", riddle, answer)
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
