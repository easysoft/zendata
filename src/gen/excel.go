package gen

import (
	"database/sql"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"github.com/easysoft/zendata/src/utils/vari"
	"log"
	"os"
	"strings"
	"time"
)

func GenerateFieldValuesFromExcel(path, sheet string, field *model.DefField) (map[string][]string, string) {
	values := map[string][]string{}

	dbName := getDbName(path)

	list := make([]string, 0)
	selectCol := ""
	firstSheet := ConvertExcelToSQLiteIfNeeded(dbName, path)
	if sheet == "" {
		sheet = firstSheet
	}

	list, selectCol = ReadDataFromSQLite(*field, dbName, sheet)
	// get index for data retrieve
	numbs := GenerateIntItems(0, (int64)(len(list)-1), 1, false, 1)
	// get data by index
	index := 0
	for _, numb := range numbs {
		item := list[numb.(int64)]

		if index >= constant.MaxNumb {
			break
		}
		if strings.TrimSpace(item) == "" {
			continue
		}

		values[selectCol] = append(values[selectCol], item)
		index = index + 1
	}

	return values, dbName
}

func getDbName(path string) (dbName string) {
	dbName = strings.Replace(path, vari.WorkDir + constant.ResDirData, "", -1)
	dbName = strings.Replace(dbName, constant.PthSep, "_", -1)
	dbName = strings.Replace(dbName, ".", "_", -1)

	return
}

func ConvertExcelToSQLiteIfNeeded(dbName string, path string) (firstSheet string) {
	excel, err := excelize.OpenFile(path)
	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_read_file", path))
		return
	}

	firstSheet = excel.GetSheetList()[0]
	if !isExcelChanged(path) {
		return
	}

	for _, sheet := range excel.GetSheetList() {
		rows, err := excel.GetRows(sheet)

		if len(rows) == 0 {
			continue
		}

		dropTemplate := `DROP TABLE IF EXISTS %s;`
		ddlTemplate := `CREATE TABLE %s (
						%s
					);`
		insertTemplate := "INSERT INTO %s (%s) VALUES %s"

		colDefine := ""
		colList := ""
		colCount := 0
		index := 0
		for _, col := range rows[0] {
			colCount++

			val := strings.TrimSpace(col)
			if index > 0 {
				colDefine = colDefine + ",\n"
				colList = colList + ", "
			}

			colProp := ""
			if val == "seq" {
				colProp = "CHAR (5) PRIMARY KEY ASC UNIQUE"
			} else {
				colProp = "VARCHAR"
			}
			colDefine = "    " + colDefine + val + " " + colProp

			colList = colList + val
			index++
		}

		valList := ""
		for rowIndex, row := range rows {
			if rowIndex == 0 {
				continue
			}

			if rowIndex > 1 {
				valList = valList + ", "
			}
			valList = valList + "("

			for colIndex, colCell := range row {
				if colIndex >= colCount {
					break
				}

				if colIndex > 0 {
					valList = valList + ", "
				}
				valList = valList + "'" + colCell + "'"
			}
			valList = valList + ")"
		}

		tableName := dbName + "_" + sheet
		dropSql := fmt.Sprintf(dropTemplate, tableName)
		ddl := fmt.Sprintf(ddlTemplate, tableName, colDefine)
		insertSql := fmt.Sprintf(insertTemplate, tableName, colList, valList)

		db, err := sql.Open("sqlite3", constant.SqliteSource)
		defer db.Close()
		_, err = db.Exec(dropSql)
		if err != nil {
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_drop_table", tableName, err.Error()))
			return
		}

		_, err = db.Exec(ddl)
		if err != nil {
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_create_table", tableName, err.Error()))
			return
		}

		_, err = db.Exec(insertSql)
		if err != nil {
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_exec_query", insertSql, err.Error()))
			return
		}
	}

	return
}

func ReadDataFromSQLite(field model.DefField, dbName string, tableName string) ([]string, string) {
	list := make([]string, 0)

	db, err := sql.Open(constant.SqliteDriver, constant.SqliteSource)
	defer db.Close()
	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_connect_sqlite", constant.SqliteSource, err.Error()))
		return list, ""
	}

	selectCol := field.Select
	from := dbName + "_" + tableName
	where := field.Where
	if !strings.Contains(where, "LIMIT") {
		total := vari.Total
		if total > constant.MaxNumb { total = constant.MaxNumb }
		if field.Limit > 0 && total > field.Limit { total = field.Limit }

		where = where + fmt.Sprintf(" LIMIT %d", total)
	}

	sqlStr := fmt.Sprintf("SELECT %s FROM %s WHERE %s", selectCol, from, where)
	rows, err := db.Query(sqlStr)
	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_exec_query", sqlStr, err.Error()))
		return list, ""
	}

	valMapArr := make([]map[string]string, 0)
	columns, err := rows.Columns()
	colNum := len(columns)

	colIndexToName := map[int]string{}
	for index, col := range columns {
		colIndexToName[index] = col
	}

	var values = make([]interface{}, colNum)
	for i, _ := range values {
		var itf string
		values[i] = &itf
	}

	for rows.Next() {
		err = rows.Scan(values...)
		if err != nil {
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_parse_row", err.Error()))
			return list, ""
		}

		rowMap := map[string]string{}
		for index, v := range values {
			item := v.(*string)

			rowMap[colIndexToName[index]] = *item
		}

		valMapArr = append(valMapArr, rowMap)
	}

	for _, item := range valMapArr {
		idx := 0
		for _, val := range item {
			if idx > 0 { break }
			list = append(list, val)
			idx++
		}
	}

	return list, selectCol
}

func isExcelChanged(path string) bool {
	db, err := sql.Open(constant.SqliteDriver, constant.SqliteSource)
	defer db.Close()
	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_connect_sqlite", constant.SqliteSource, err.Error()))
		return true
	}

	sqlStr := fmt.Sprintf("SELECT id, name, changeTime FROM %s WHERE name = '%s';",
		constant.SqliteTrackTable, path)
	rows, err := db.Query(sqlStr)
	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_exec_query", sqlStr, err.Error()))
		return true
	}

	fileChangeTime := getFileModTime(path).Unix()

	found := false
	changed := false
	for rows.Next() {
		found = true
		var id int
		var name string
		var changeTime int64

		err = rows.Scan(&id, &name, &changeTime)
		if err != nil {
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_parse_row", err.Error()))
			changed = true
			break
		}

		if path == name && changeTime < fileChangeTime {
			changed = true
			break
		}
	}
	rows.Close()

	if !found { changed = true }

	if changed {
		if !found {
			sqlStr = fmt.Sprintf("INSERT INTO %s(name, changeTime) VALUES('%s', %d)",
				constant.SqliteTrackTable, path, fileChangeTime)
		} else {
			sqlStr = fmt.Sprintf("UPDATE %s SET changeTime = %d WHERE name = '%s'",
				constant.SqliteTrackTable, fileChangeTime, path)
		}

		_, err = db.Exec(sqlStr)
		if err != nil {
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_exec_query", sqlStr, err.Error()))
		}
	}

	return changed
}

func getFileModTime(path string) time.Time {
	f, err := os.Open(path)
	if err != nil {
		log.Println("open file error:" + path)
		return time.Now()
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Println("stat fileinfo error")
		return time.Now()
	}

	fileChangeTime := fi.ModTime()

	//timeStr := fileChangeTime.Format("2006-01-02 15:04:05")
	//logUtils.Screen(i118Utils.I118Prt.Sprintf("file_change_time", timeStr))

	return fileChangeTime
}