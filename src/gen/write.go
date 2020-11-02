package gen

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	constant "github.com/easysoft/zendata/src/utils/const"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/mattn/go-runewidth"
)

const (
	sheetName = "Sheet1"
)

func Write(rows [][]string, format string, table string, colIsNumArr []bool,
		fields []string) (lines []interface{}) {

	f := excelize.NewFile()
	index := f.NewSheet(sheetName)
	f.SetActiveSheet(index)

	if format == constant.FormatExcel {
		printExcelHeader(fields, f)
	}

	for i, row := range rows {
		for j, col := range row {
			col = replacePlaceholder(col)
			field := vari.TopFieldMap[fields[j]]
			if field.Length > runewidth.StringWidth(col) {
				col = stringUtils.AddPad(col, field)
			}

			colName, _ := excelize.CoordinatesToCellName(j + 1, i + 2)
			f.SetCellValue(sheetName, colName, col)
		}
	}

	if err := f.SaveAs(logUtils.FilePath); err != nil {
		fmt.Println(err)
	}

	return
}

func printExcelHeader(fields []string, f *excelize.File) {
	headerLine := ""
	for idx, field := range fields {
		colName, _ := excelize.CoordinatesToCellName(idx + 1, 1)
		f.SetCellValue(sheetName, colName, field)
	}

	logUtils.PrintLine(headerLine)
}