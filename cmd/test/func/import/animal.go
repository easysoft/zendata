package main

import (
	"fmt"
	"strings"

	"github.com/easysoft/zendata/cmd/test/func/comm"
	"github.com/easysoft/zendata/cmd/test/func/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
)

type DataAnimal struct {
	model.DataComm
}

func (DataAnimal) TableName() string {
	return "biz_data_animal"
}

func main() {
	tableName := "animal"
	filePath := "/Users/aaron/work/zentao/product/zd/行业数据/动物.txt"

	var insertTemplate = "INSERT INTO %s (name) VALUES %s;"
	var createTableTempl = `CREATE TABLE IF NOT EXISTS %s (
		id bigint auto_increment,
		name varchar(1000) not null,
    
		primary key(id)
	) engine=innodb default charset=utf8 auto_increment=1;`

	tableName = "biz_data_" + tableName
	db := comm.GetDB()
	db.Exec(fmt.Sprintf(comm.TruncateTable, tableName))

	createTableSql := fmt.Sprintf(createTableTempl, tableName)
	err := db.Exec(createTableSql).Error
	if err != nil {
		fmt.Printf("create table %s failed, err %s", tableName, err.Error())
		return
	}

	content := fileUtils.ReadFile(filePath)
	records := make([]DataAnimal, 0)

	for _, line := range strings.Split(content, "\n") {
		col1 := strings.Split(strings.TrimSpace(line), " ")[0]

		po := DataAnimal{
			DataComm: model.DataComm{
				Name: col1,
			},
		}

		records = append(records, po)
	}

	insertSqlArr := make([]string, 0)

	for _, record := range records {
		insert := fmt.Sprintf("('%s')", record.Name)
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
