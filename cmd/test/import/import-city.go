package main

import (
	"fmt"
	"github.com/easysoft/zendata/cmd/test/comm"
	"github.com/easysoft/zendata/cmd/test/model"
)

func main() {
	filePath := "data/city/cn.v1.xlsx"
	sheetName := "city"

	db := comm.GetDB()
	db.Exec(fmt.Sprintf(comm.TruncateTable, model.DataCity{}.TableName()))
	db.AutoMigrate(
		&model.DataCity{},
	)

	records := comm.GetExcelTable(filePath, sheetName)

	for _, record := range records {
		po := model.DataCity{
			Name:        record["city"].(string),
			Code:        record["cityCode"].(string),
			ZipCode:     record["zipCode"].(string),
			State:       record["state"].(string),
			StateShort:  record["stateShort"].(string),
			StateShort2: record["stateShort2"].(string),
		}

		db.Save(&po)
	}
}
