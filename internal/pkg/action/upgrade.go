package action

import (
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	_ "github.com/mattn/go-sqlite3"
	"hash/crc32"
	"strconv"
)

func AddExcelRow(excel *excelize.File, sheet string, rowIndex int, cols []interface{}) string {
	start := byte('A')
	var numb string
	line := ""
	for index, col := range cols {
		numb = string(start + byte(index))

		colNumb := numb + strconv.Itoa(rowIndex)
		excel.SetCellValue(sheet, colNumb, col)
		line = line + col.(string)
	}

	// 校验位
	if rowIndex > 1 {
		numb = string(start + byte(len(cols)))
		crcField := crc32.ChecksumIEEE([]byte(line))
		excel.SetCellValue(sheet, numb+strconv.Itoa(rowIndex), crcField)
	}

	return numb
}
