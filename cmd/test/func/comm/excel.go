package comm

import (
	"fmt"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
)

func GetExcelTable(filePath, sheetName string) (records []map[string]interface{}) {

	excel, err := excelize.OpenFile(filePath)
	if err != nil {
		fmt.Printf("fail to read file %s, error: %s", filePath, err.Error())
		return
	}

	fileName := fileUtils.GetFileName(filePath)
	fileName = strings.TrimSuffix(fileName, "词库")

	for _, sheet := range excel.GetSheetList() {
		if sheetName != sheet {
			continue
		}

		rows, _ := excel.GetRows(sheet)
		if len(rows) == 0 {
			continue
		}

		colMap := map[int]string{}
		colCount := 0
		for colIndex, colVal := range rows[0] {
			colMap[colIndex] = strings.TrimSpace(colVal)
			colCount++
		}

		for rowIndex, row := range rows {
			if rowIndex == 0 { // ignore header
				continue
			}

			record := map[string]interface{}{}

			for colIndex, col := range row {
				if colIndex >= colCount {
					break
				}

				colName := colMap[colIndex]
				colVal := strings.TrimSpace(col)
				record[colName] = colVal
			}

			records = append(records, record)
		}
	}

	return
}

func GetExcelFirstSheet(filePath string) (sheetName string, rows [][]string) {
	excel, err := excelize.OpenFile(filePath)
	if err != nil {
		fmt.Printf("fail to read file %s, error: %s", filePath, err.Error())
		return
	}

	sheetName = excel.GetSheetList()[0]
	rows, err = excel.GetRows(sheetName)

	if err != nil {
		fmt.Println(err)
	}

	return
}
