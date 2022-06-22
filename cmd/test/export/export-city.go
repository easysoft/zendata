package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/easysoft/zendata/cmd/test/import/comm"
	"github.com/easysoft/zendata/cmd/test/import/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"path/filepath"
	"reflect"
)

func main() {
	filePath := "data/city/v1.xlsx"
	sheetName := "city"

	fileUtils.MkDirIfNeeded(filepath.Dir(filePath))

	db := comm.GetDB()
	db.AutoMigrate(
		&model.DataCity{},
	)

	pos := make([]model.DataCity, 0)
	db.Where("NOT deleted").Find(&pos)

	f := excelize.NewFile()
	index := f.NewSheet(sheetName)

	var infos []model.TableInfo
	db.Raw("desc " + model.DataCity{}.TableName()).Scan(&infos)

	excelColNameArr, excelColNameHeader := comm.GetExcelColsByTableDef(infos)
	fieldNames := comm.GetStructFields(model.DataCity{})

	// gen headers
	for index, name := range excelColNameHeader {
		excelColName := excelColNameArr[index]
		excelColId := fmt.Sprintf("%s%d", excelColName, 1)

		f.SetCellValue(sheetName, excelColId, name)
	}

	// gen rows
	for rowIndex, po := range pos {
		for fieldIndex, fieldName := range fieldNames {
			val := ""

			if fieldName == "Id" {
				val = fmt.Sprintf("%d", reflect.ValueOf(po).FieldByName(fieldName).Int())
			} else {
				val = reflect.ValueOf(po).FieldByName(fieldName).String()
			}

			excelColName := excelColNameArr[fieldIndex]
			excelColId := fmt.Sprintf("%s%d", excelColName, rowIndex+2)

			f.SetCellValue(sheetName, excelColId, val)
		}
	}

	f.SetActiveSheet(index)
	f.SaveAs(filePath)
}
