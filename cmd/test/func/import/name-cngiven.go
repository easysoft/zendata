package main

import (
	"fmt"

	"github.com/easysoft/zendata/cmd/test/func/comm"
	"github.com/easysoft/zendata/cmd/test/func/model"
)

func main() {
	filePath := "data/name/cn.given.v1.xlsx"
	sheetName := "中文名"

	db := comm.GetDB()
	db.Exec(fmt.Sprintf(comm.TruncateTable, model.DataChineseGiven{}.TableName()))
	db.AutoMigrate(
		&model.DataChineseGiven{},
	)

	records := comm.GetExcelTable(filePath, sheetName)

	for _, record := range records {
		po := model.DataChineseGiven{
			Name:   record["name"].(string),
			Pinyin: record["pinyin"].(string),
			Sex:    record["sex"].(string),
		}

		db.Save(&po)
	}
}
