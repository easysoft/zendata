package main

import (
	"fmt"
	"github.com/easysoft/zendata/cmd/test/import/comm"
	"github.com/easysoft/zendata/cmd/test/import/model"
)

func main() {
	filePath := "data/name/en.family.v1.xlsx"
	sheetName := "英文姓"

	db := comm.GetDB()
	db.Exec(fmt.Sprintf(comm.TruncateTable, model.DataEnglishFamily{}.TableName()))
	db.AutoMigrate(
		&model.DataEnglishFamily{},
	)

	records := comm.GetExcelTable(filePath, sheetName)

	for _, record := range records {
		po := model.DataEnglishFamily{
			Name:  record["name"].(string),
			Index: record["index"].(string),
		}

		db.Save(&po)
	}
}
