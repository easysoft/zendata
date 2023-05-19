package main

import (
	"fmt"

	"github.com/easysoft/zendata/cmd/test/func/comm"
	"github.com/easysoft/zendata/cmd/test/func/model"
)

func main() {
	filePath := "data/name/en.given.v1.xlsx"
	sheetName := "英文名"

	db := comm.GetDB()
	db.Exec(fmt.Sprintf(comm.TruncateTable, model.DataEnglishGiven{}.TableName()))
	db.AutoMigrate(
		&model.DataEnglishGiven{},
	)

	records := comm.GetExcelTable(filePath, sheetName)

	for _, record := range records {
		po := model.DataEnglishGiven{
			Name:  record["name"].(string),
			Index: record["index"].(string),
			Sex:   record["sex"].(string),
		}

		db.Save(&po)
	}
}
