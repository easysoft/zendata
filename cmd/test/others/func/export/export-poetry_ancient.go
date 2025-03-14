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
	filePath := "data/poetry_ancient/v1.xlsx"
	sheetName := "poetry_ancient"

	fileUtils.MkDirIfNeeded(filepath.Dir(filePath))

	db := comm.GetDB()
	db.AutoMigrate(
		&model.PoetryAncient{},
	)

	pos := make([]model.PoetryAncient, 0)
	db.Where("NOT deleted").Find(&pos)

	f := excelize.NewFile()
	index := f.NewSheet(sheetName)
	f.SetActiveSheet(index)

	sheet1 := f.GetSheetName(0)
	f.DeleteSheet(sheet1)

	var infos []model.TableInfo
	db.Raw("desc " + model.PoetryAncient{}.TableName()).Scan(&infos)

	excelColNameArr, excelColNameHeader := comm.GetExcelColsByTableDef(infos)
	fieldNames := comm.GetStructFields(model.PoetryAncient{})

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

			kind := reflect.ValueOf(po).FieldByName(fieldName).Kind()

			if fieldName == "Id" {
				val = fmt.Sprintf("%d", reflect.ValueOf(po).FieldByName(fieldName).Uint())
			} else if kind == reflect.Int {
				val = fmt.Sprintf("%d", reflect.ValueOf(po).FieldByName(fieldName).Int())
			} else if kind == reflect.Uint {
				val = fmt.Sprintf("%d", reflect.ValueOf(po).FieldByName(fieldName).Uint())
			} else {
				val = fmt.Sprintf("%v", reflect.ValueOf(po).FieldByName(fieldName))
			}

			excelColName := excelColNameArr[fieldIndex]
			excelColId := fmt.Sprintf("%s%d", excelColName, rowIndex+2)

			f.SetCellValue(sheetName, excelColId, val)
		}
	}

	f.SaveAs(filePath)
}
