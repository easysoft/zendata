package action

import (
	"database/sql"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	constant "github.com/easysoft/zendata/src/utils/const"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	numbUtils "github.com/easysoft/zendata/src/utils/numb"
	_ "github.com/mattn/go-sqlite3"
	"hash/crc32"
	"strconv"
)

func Upgrade() {
	// TODO: update date from remote server

	// read data from db
	db, err := sql.Open(constant.SqliteDriver, constant.SqliteSource)
	if err != nil {
		logUtils.Screen("fail to open " + constant.SqliteSource + ": " + err.Error())
		return
	}

	sql := "SELECT id, name, state, zipCode, cityCode FROM city"
	rows, err := db.Query(sql)
	if err != nil {
		logUtils.Screen("fail to exec query " + sql + ": " + err.Error())
		return
	}

	sqls := make([]string, 0)

	sheet := "city"
	excel := excelize.NewFile()
	excel.SetSheetName("Sheet1", sheet)
	headerData := []interface{}{"seq", "name", "state", "zipCode", "cityCode", "crc"}
	colNumb := AddExcelRow(excel, sheet, 1, headerData)

	rowIndex := 1
	for rows.Next() {
		var seq string
		var id int
		var name string
		var state string
		var zipCode string
		var cityCode string

		err = rows.Scan(&id, &name, &state, &zipCode, &cityCode)
		if err != nil {
			logUtils.Screen("fail to get sqlite3 row: " + err.Error())
			return
		}

		seq = numbUtils.NumToBHex(id)
		fmt.Println(seq, name, state, zipCode, cityCode)

		// gen update sql
		sql := fmt.Sprintf("UPDATE city SET seq = '%s' where id = %d;", numbUtils.NumToBHex(id), id)
		sqls = append(sqls, sql)

		// gen excel row
		rowIndex = rowIndex + 1
		rowData := []interface{}{seq, name, state, zipCode, cityCode}
		AddExcelRow(excel, sheet, rowIndex, rowData)
	}

	headerStyle, err := excel.NewStyle(constant.ExcelHeader)
	excel.SetCellStyle(sheet, "A1", colNumb + "1", headerStyle)

	borderStyle, err := excel.NewStyle(constant.ExcelBorder)
	excel.SetCellStyle(sheet, "A1", colNumb + strconv.Itoa(rowIndex), borderStyle)

	// update seq column
	for _, sql := range sqls {
		_, err := db.Exec(sql)
		if err != nil {
			logUtils.Screen("fail to update row: " + err.Error())
		}
	}

	err = excel.SaveAs(constant.ExcelFile)
	if err != nil {
		logUtils.Screen("fail to save excel: " + err.Error())
	}
}

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