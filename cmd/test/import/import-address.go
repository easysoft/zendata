package main

import (
	"fmt"
	"github.com/easysoft/zendata/cmd/test/import/comm"
	"github.com/easysoft/zendata/cmd/test/import/model"
)

func main() {
	filePath := "data/address/cn.v1.xlsx"
	sheetName := "china"

	db := comm.GetDB()
	db.Exec(fmt.Sprintf(comm.TruncateTable, model.DataCity{}.TableName()))
	db.AutoMigrate(
		&model.DataCity{},
	)

	records := comm.GetExcelTable(filePath, sheetName)

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
