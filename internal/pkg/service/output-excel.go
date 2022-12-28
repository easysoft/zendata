package service

import (
	"encoding/csv"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
)

const (
	sheetName = "Sheet1"
)

var (
	csvWriter *csv.Writer
)

func (s *OutputService) GenExcel() {
	records := s.GenObjs()

	var f *excelize.File
	csvData := make([][]string, 0)

	if vari.GlobalVars.OutputFormat == consts.FormatExcel {
		f = excelize.NewFile()
		index := f.NewSheet(sheetName)
		f.SetActiveSheet(index)

		s.printExcelHeader(f)
	} else if vari.GlobalVars.OutputFormat == consts.FormatCsv {
		s.printCsvHeader(&csvData)
		csvWriter = csv.NewWriter(logUtils.OutputFileWriter)
	}

	for i, record := range records {
		csvRow := make([]string, 0)

		j := 0
		for _, field := range vari.GlobalVars.ExportFields {
			val := record[field]

			if vari.GlobalVars.OutputFormat == consts.FormatExcel {
				colName, _ := excelize.CoordinatesToCellName(j+1, i+2)
				f.SetCellValue(sheetName, colName, val)

			} else if vari.GlobalVars.OutputFormat == consts.FormatCsv {
				csvRow = append(csvRow, val.(string))
			}

			j++
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

func (s *OutputService) printExcelHeader(f *excelize.File) {
	headerLine := ""
	for idx, field := range vari.GlobalVars.ExportFields {
		colName, _ := excelize.CoordinatesToCellName(idx+1, 1)
		f.SetCellValue(sheetName, colName, field)
	}

	logUtils.PrintLine(headerLine)
}

func (s *OutputService) printCsvHeader(csvData *[][]string) {
	csvRow := make([]string, 0)

	for _, field := range vari.GlobalVars.ExportFields {
		csvRow = append(csvRow, field)
	}

	*csvData = append(*csvData, csvRow)
}
