package gen

import (
	"encoding/csv"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/mattn/go-runewidth"
)

const (
	sheetName = "Sheet1"
)

var (
	csvWriter *csv.Writer
)

func Write(rows [][]string, table string, colIsNumArr []bool,
	fields []string) (lines []interface{}) {

	f := excelize.NewFile()
	index := f.NewSheet(sheetName)
	f.SetActiveSheet(index)

	if vari.GlobalVars.OutputFormat == consts.FormatExcel {
		printExcelHeader(fields, f)
	} else if vari.GlobalVars.OutputFormat == consts.FormatCsv {
		csvWriter = csv.NewWriter(logUtils.OutputFileWriter)
	}

	csvData := make([][]string, 0)
	for i, cols := range rows {
		csvRow := make([]string, 0)

		for j, col := range cols {
			col = replacePlaceholder(col)
			field := vari.GlobalVars.TopFieldMap[fields[j]]
			if field.Length > runewidth.StringWidth(col) {
				//col = stringUtils.AddPad(col, field)
			}

			if vari.GlobalVars.OutputFormat == consts.FormatExcel {
				colName, _ := excelize.CoordinatesToCellName(j+1, i+2)
				f.SetCellValue(sheetName, colName, col)

			} else if vari.GlobalVars.OutputFormat == consts.FormatCsv {
				csvRow = append(csvRow, col)
			}
		}
		csvData = append(csvData, csvRow)
	}

	var err error
	if vari.GlobalVars.OutputFormat == consts.FormatExcel {
		err = f.SaveAs(logUtils.OutputFilePath)
	} else if vari.GlobalVars.OutputFormat == consts.FormatCsv {
		err = csvWriter.WriteAll(csvData)
		csvWriter.Flush()
	}
	if err != nil {
		logUtils.PrintErrMsg(err.Error())
	}

	return
}

func printExcelHeader(fields []string, f *excelize.File) {
	headerLine := ""
	for idx, field := range fields {
		colName, _ := excelize.CoordinatesToCellName(idx+1, 1)
		f.SetCellValue(sheetName, colName, field)
	}

	logUtils.PrintLine(headerLine)
}