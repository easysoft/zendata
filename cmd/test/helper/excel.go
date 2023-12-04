package testHelper

import (
	"log"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func GetExcelData(pth string, sheetIndex, rowIndex, colIndex int) (ret interface{}) {
	excel, err := excelize.OpenFile(pth)
	if err != nil {
		log.Println("fail to read file " + pth + ", error: " + err.Error())
		return
	}

	for sheetIdx, sheet := range excel.GetSheetList() {
		if sheetIdx != sheetIndex {
			continue
		}
		if sheetIdx > sheetIndex {
			break
		}

		rows, _ := excel.GetRows(sheet)
		if len(rows) == 0 {
			continue
		}

		for rowIdx, row := range rows {
			if rowIdx != rowIndex {
				continue
			}
			if rowIdx > rowIndex {
				break
			}

			for colIdx, col := range row {
				if colIdx != colIndex {
					continue
				}
				if colIdx > colIndex {
					break
				}

				ret = col
			}
		}
	}

	return
}
