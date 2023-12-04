package main

import (
	"fmt"
	"path/filepath"
	"reflect"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/easysoft/zendata/cmd/test/others/func/comm"
	"github.com/easysoft/zendata/cmd/test/others/func/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
)

func main() {
	filePath := "data/song/v1.xlsx"
	sheetName := "song"

	fileUtils.MkDirIfNeeded(filepath.Dir(filePath))

	db := comm.GetDB()
	db.AutoMigrate(
		&model.Song{},
	)

	pos := make([]model.Song, 0)
	db.Where("NOT deleted").Find(&pos)

	f := excelize.NewFile()
	index := f.NewSheet(sheetName)
	f.SetActiveSheet(index)

	sheet1 := f.GetSheetName(0)
	f.DeleteSheet(sheet1)

	var infos []model.TableInfo
	db.Raw("desc " + model.Song{}.TableName()).Scan(&infos)

	excelColNameArr, excelColNameHeader := comm.GetExcelColsByTableDef(infos)
	fieldNames := comm.GetStructFields(model.Song{})

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
				val = fmt.Sprintf("%d", reflect.ValueOf(po).FieldByName(fieldName).Uint())
			} else if reflect.ValueOf(po).FieldByName(fieldName).Kind() == reflect.Int {
				val = fmt.Sprintf("%d", reflect.ValueOf(po).FieldByName(fieldName).Int())
			} else {
				val = reflect.ValueOf(po).FieldByName(fieldName).String()
			}

			excelColName := excelColNameArr[fieldIndex]
			excelColId := fmt.Sprintf("%s%d", excelColName, rowIndex+2)

			f.SetCellValue(sheetName, excelColId, val)
		}
	}

	f.SaveAs(filePath)
}
