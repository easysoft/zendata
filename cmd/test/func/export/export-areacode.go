package main

import (
	"fmt"
	"path/filepath"
	"reflect"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/easysoft/zendata/cmd/test/func/comm"
	"github.com/easysoft/zendata/cmd/test/func/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
)

func main() {
	filePath := "data/areacode/v1.xlsx"
	sheetNameCity := "city"
	sheetNameCountry := "country"

	fileUtils.RmFile(filePath)

	fileUtils.MkDirIfNeeded(filepath.Dir(filePath))

	db := comm.GetDB()

	pos := make([]model.AreaCodeCity, 0)
	db.Where("NOT deleted").Find(&pos)

	f := excelize.NewFile()

	index := f.NewSheet(sheetNameCity)
	f.SetActiveSheet(index)

	sheet1 := f.GetSheetName(0)
	f.DeleteSheet(sheet1)

	var infos []model.TableInfo
	db.Raw("desc " + model.AreaCodeCity{}.TableName()).Scan(&infos)

	excelColNameArr, excelColNameHeader := comm.GetExcelColsByTableDef(infos)
	fieldNames := comm.GetStructFields(model.AreaCodeCity{})

	// gen headers
	for index, name := range excelColNameHeader {
		excelColName := excelColNameArr[index]
		excelColId := fmt.Sprintf("%s%d", excelColName, 1)

		f.SetCellValue(sheetNameCity, excelColId, name)
	}

	// gen rows
	for rowIndex, po := range pos {
		for fieldIndex, fieldName := range fieldNames {
			val := ""

			if fieldName == "Id" {
				val = fmt.Sprintf("%d", reflect.ValueOf(po).FieldByName(fieldName).Uint())
			} else {
				val = reflect.ValueOf(po).FieldByName(fieldName).String()
			}

			excelColName := excelColNameArr[fieldIndex]
			excelColId := fmt.Sprintf("%s%d", excelColName, rowIndex+2)

			f.SetCellValue(sheetNameCity, excelColId, val)
		}
	}

	//
	pos2 := make([]model.AreaCodeCountry, 0)
	db.Where("NOT deleted").Find(&pos2)

	index2 := f.NewSheet(sheetNameCountry)
	f.SetActiveSheet(index2)

	var infos2 []model.TableInfo
	db.Raw("desc " + model.AreaCodeCountry{}.TableName()).Scan(&infos2)

	excelColNameArr2, excelColNameHeader2 := comm.GetExcelColsByTableDef(infos2)
	fieldNames2 := comm.GetStructFields(model.AreaCodeCountry{})

	// gen headers
	for index, name := range excelColNameHeader2 {
		excelColName := excelColNameArr2[index]
		excelColId := fmt.Sprintf("%s%d", excelColName, 1)

		f.SetCellValue(sheetNameCountry, excelColId, name)
	}

	// gen rows
	for rowIndex, po := range pos2 {
		for fieldIndex, fieldName := range fieldNames2 {
			val := ""

			if fieldName == "Id" {
				val = fmt.Sprintf("%d", reflect.ValueOf(po).FieldByName(fieldName).Uint())
			} else {
				val = reflect.ValueOf(po).FieldByName(fieldName).String()
			}

			excelColName := excelColNameArr2[fieldIndex]
			excelColId := fmt.Sprintf("%s%d", excelColName, rowIndex+2)

			f.SetCellValue(sheetNameCountry, excelColId, val)
		}
	}

	f.SaveAs(filePath)
}
