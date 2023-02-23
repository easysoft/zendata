package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/easysoft/zendata/cmd/test/comm"
	"github.com/easysoft/zendata/cmd/test/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"path/filepath"
	"reflect"
)

func main() {
	filePath := "data/company/v1.xlsx"
	sheetName1 := "company"

	fileUtils.MkDirIfNeeded(filepath.Dir(filePath))

	db := comm.GetDB()
	db.AutoMigrate(
		&model.DataCompany{},
	)

	pos1 := make([]model.DataCompany, 0)
	db.Order("id ASC").Find(&pos1)

	f := excelize.NewFile()
	index1 := f.NewSheet(sheetName1)
	f.SetActiveSheet(index1)

	sheet1 := f.GetSheetName(0)
	f.DeleteSheet(sheet1)

	var infos1 []model.TableInfo
	db.Raw("desc " + model.DataCompany{}.TableName()).Scan(&infos1)

	excelColNameArr1, excelColNameHeader1 := comm.GetExcelColsByTableDef(infos1)
	fieldNames1 := comm.GetStructFields(model.DataCompany{})

	// gen headers
	for index, name := range excelColNameHeader1 {
		excelColName := excelColNameArr1[index]
		excelColId := fmt.Sprintf("%s%d", excelColName, 1)

		f.SetCellValue(sheetName1, excelColId, name)
	}

	// gen rows
	for rowIndex, po := range pos1 {
		for fieldIndex, fieldName := range fieldNames1 {
			val := ""

			if fieldName == "Id" {
				val = fmt.Sprintf("%d", reflect.ValueOf(po).FieldByName(fieldName).Uint())
			} else {
				val = reflect.ValueOf(po).FieldByName(fieldName).String()
			}

			excelColName := excelColNameArr1[fieldIndex]
			excelColId := fmt.Sprintf("%s%d", excelColName, rowIndex+2)

			f.SetCellValue(sheetName1, excelColId, val)
		}
	}

	sheetName2 := "company_abbreviation"
	pos := make([]model.DataCompanyAbbreviation, 0)
	db.Order("id ASC").Find(&pos)

	index2 := f.NewSheet(sheetName2)
	f.SetActiveSheet(index2)

	var infos []model.TableInfo
	db.Raw("desc " + model.DataCompanyAbbreviation{}.TableName()).Scan(&infos)

	excelColNameArr2, excelColNameHeader2 := comm.GetExcelColsByTableDef(infos)
	fieldNames2 := comm.GetStructFields(model.DataCompanyAbbreviation{})

	// gen headers
	for index, name := range excelColNameHeader2 {
		excelColName := excelColNameArr2[index]
		excelColId := fmt.Sprintf("%s%d", excelColName, 1)

		f.SetCellValue(sheetName2, excelColId, name)
	}

	// gen rows
	for rowIndex, po := range pos {
		for fieldIndex, fieldName := range fieldNames2 {
			val := ""

			if fieldName == "Id" {
				val = fmt.Sprintf("%d", reflect.ValueOf(po).FieldByName(fieldName).Uint())
			} else {
				val = reflect.ValueOf(po).FieldByName(fieldName).String()
			}

			excelColName := excelColNameArr2[fieldIndex]
			excelColId := fmt.Sprintf("%s%d", excelColName, rowIndex+2)

			f.SetCellValue(sheetName2, excelColId, val)
		}
	}

	f.SetActiveSheet(index1)
	f.SaveAs(filePath)
}
