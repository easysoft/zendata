package main

import (
	"fmt"
	"github.com/easysoft/zendata/src/test/import/comm"
	"github.com/easysoft/zendata/src/test/import/model"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
)

func main() {
	filePath := "data/name/cn.family.v1.xlsx"
	sheetName := "中文姓"

	db := comm.GetDB()
	db.Exec(fmt.Sprintf(comm.TruncateTable, model.DataChineseFamily{}.TableName()))
	db.AutoMigrate(
		&model.DataChineseFamily{},
	)

	records := comm.GetExcelTable(filePath, sheetName)

	for _, record := range records {
		po := model.DataChineseFamily{
			Name:   record["name"].(string),
			Pinyin: record["pinyin"].(string),
			Double: stringUtils.ParseBool(record["double"].(string)),
		}

		db.Save(&po)
	}
}
