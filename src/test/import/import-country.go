package main

import (
	"github.com/easysoft/zendata/src/test/import/comm"
	"github.com/easysoft/zendata/src/test/import/model"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
)

func main() {
	filePath := "data/country/v1.xlsx"
	sheetName := "country"

	db := comm.GetDB()

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
