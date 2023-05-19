package main

import (
	"fmt"

	"github.com/easysoft/zendata/cmd/test/func/comm"
	"github.com/easysoft/zendata/cmd/test/func/model"
	stringUtils "github.com/easysoft/zendata/pkg/utils/string"
)

func main() {
	filePath := "data/country/v1.xlsx"
	sheetName := "country"

	db := comm.GetDB()
	db.Exec(fmt.Sprintf(comm.TruncateTable, model.DataCountry{}.TableName()))
	db.AutoMigrate(
		&model.DataCountry{},
	)

	records := comm.GetExcelTable(filePath, sheetName)

	for _, record := range records {
		po := model.DataCountry{
			ContinentId:  stringUtils.ParseInt(record["continent_id"].(string)),
			Continent:    record["continent"].(string),
			AreaCode:     record["areacode"].(string),
			EnglishShort: record["enshort"].(string),
			EnglishFull:  record["enfull"].(string),
			ChineseShort: record["cnshort"].(string),
			ChineseFull:  record["cnfull"].(string),
		}

		db.Save(&po)
	}
}
