package main

import (
	"fmt"
	"github.com/easysoft/zendata/src/test/import/comm"
	"github.com/easysoft/zendata/src/test/import/model"
)

func main() {
	tableName := model.DataCity{}.TableName()
	db := comm.GetDB()

	truncateTableSql := fmt.Sprintf(comm.TruncateTable, tableName)
	db.Exec(truncateTableSql)

	createTableSql := fmt.Sprintf(comm.CreateTableTempl, tableName)
	err := db.Exec(createTableSql).Error
	if err != nil {
		fmt.Printf("create table %s failed, err %s", tableName, err.Error())
		return
	}

	records := comm.GetExcelTable()

	for _, record := range records {
		po := model.DataCity{
			Name:    record["city"].(string),
			Code:    record["cityCode"].(string),
			ZipCode: record["zipCode"].(string),
			State:   record["state"].(string),
		}

		db.Save(&po)
	}
}
