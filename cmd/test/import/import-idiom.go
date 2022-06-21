package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/easysoft/zendata/cmd/test/import/comm"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"github.com/kataras/iris/v12"
	"strings"
)

func main() {
	var insertTemplate = "INSERT INTO %s (derivation, example, explanation, pinyin, word, abbreviation) VALUES %s;"
	var createTableTempl = `CREATE TABLE IF NOT EXISTS %s (
		id bigint auto_increment,

		derivation varchar(1000) not null,
		example varchar(1000) not null,
		explanation varchar(1000) not null,
		pinyin varchar(1000) not null,
		word varchar(1000) not null,
		abbreviation varchar(1000) not null,

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
	records := make([]iris.Map, 0)
	json.Unmarshal(content, &records)

	insertSqlArr := make([]string, 0)

	for _, record := range records {
		derivation := strings.ReplaceAll(record["derivation"].(string), "'", "''")
		example := strings.ReplaceAll(record["example"].(string), "'", "''")
		explanation := strings.ReplaceAll(record["explanation"].(string), "'", "''")
		pinyin := strings.ReplaceAll(record["pinyin"].(string), "'", "''")
		word := strings.ReplaceAll(record["word"].(string), "'", "''")
		abbreviation := strings.ReplaceAll(record["abbreviation"].(string), "'", "''")

		insert := fmt.Sprintf("('%s', '%s', '%s', '%s', '%s', '%s')",
			derivation, example, explanation, pinyin, word, abbreviation)
		insertSqlArr = append(insertSqlArr, insert)
	}

	for i := 0; i < 10000; i++ {
		start := i * 1000
		end := (i + 1) * 1000

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
