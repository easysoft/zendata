package main

import (
	"fmt"
	"github.com/easysoft/zendata/src/test/import/comm"
	"github.com/easysoft/zendata/src/test/import/model"
)

func main() {
	filePath := "data/color/v1.xlsx"
	sheetName := "color"

	db := comm.GetDB()
	db.Exec(fmt.Sprintf(comm.TruncateTable, model.DataColor{}.TableName()))
	db.AutoMigrate(
		&model.DataColor{},
	)

	records := comm.GetExcelTable(filePath, sheetName)

	for _, record := range records {
		po := model.DataColor{
			Chinese: record["chinese"].(string),
			English: record["english"].(string),
			Hex:     record["hex"].(string),
			Rgb:     record["rgb"].(string),
		}

		db.Save(&po)
	}
}
