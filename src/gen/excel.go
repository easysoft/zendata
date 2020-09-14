package gen

import (
	"database/sql"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"github.com/easysoft/zendata/src/utils/vari"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func GenerateFieldValuesFromExcel(filePath, sheet string, field *model.DefField) (map[string][]string, string) {
	values := map[string][]string{}

	dbName := getDbName(filePath)

	if !fileUtils.IsDir(filePath) { // file
		firstSheet := ConvertSingleExcelToSQLiteIfNeeded(dbName, filePath)
		if sheet == "" {
			sheet = firstSheet
		}
	} else { // dir
		ConvertExcelsToSQLiteIfNeeded(dbName, filePath)
	}

	list, selectCol := ReadDataFromSQLite(*field, dbName, sheet)
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
	dbName = strings.Replace(path, vari.WorkDir + constant.ResDirData + constant.PthSep, "", -1)
	dbName = strings.Replace(dbName, constant.PthSep, "_", -1)
	dbName = strings.Replace(dbName, ".", "_", -1)

	return
}

func ConvertSingleExcelToSQLiteIfNeeded(dbName string, path string) (firstSheet string) {
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
			colDefine = "    " + colDefine + "`" + val + "` " + colProp

			colList = colList + "`" + val + "`"
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
				colCell = strings.Replace(colCell, "'", "''", -1)
				valList = valList + "'" + colCell + "'"
			}
			valList = valList + ")"
		}

		tableName := dbName + "_" + sheet
		dropSql := fmt.Sprintf(dropTemplate, tableName)
		ddl := fmt.Sprintf(ddlTemplate, tableName, colDefine)
		insertSql := fmt.Sprintf(insertTemplate, tableName, colList, valList)

		db, err := sql.Open(constant.SqliteDriver, constant.SqliteData)
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

func ConvertExcelsToSQLiteIfNeeded(tableName string, dir string) {
	if !isExcelChanged(dir) {
		return
	}

	files := make([]string, 0)
	fileUtils.GetFilesByExtInDir(dir, ".xlsx", &files)

	seq := 1
	ddlFields := make([]string, 0)
	insertSqls := make([]string, 0)

	colMap := map[string]bool{}
	for _, file := range files {
		importExcel(file, tableName, &seq, &ddlFields, &insertSqls, &colMap)
	}

	db, err := sql.Open(constant.SqliteDriver, constant.SqliteData)
	defer db.Close()

	if err != nil {
		log.Println(i118Utils.I118Prt.Sprintf("fail_to_connect_sqlite", constant.SqliteData, err.Error()))
		return
	}

	dropSql := `DROP TABLE IF EXISTS ` + tableName + `;`
	_, err = db.Exec(dropSql)
	if err != nil {
		log.Println(i118Utils.I118Prt.Sprintf("fail_to_drop_table", tableName, err.Error()))
		return
	}

	ddlTemplate := "CREATE TABLE " + tableName + "(\n" +
		"\t`seq` CHAR (5) PRIMARY KEY ASC UNIQUE,\n" +
		"%s" +
		"\n);"
	ddlSql := fmt.Sprintf(ddlTemplate, strings.Join(ddlFields, ", \n"))
	_, err = db.Exec(ddlSql)
	if err != nil {
		log.Println(i118Utils.I118Prt.Sprintf("fail_to_create_table", tableName, err.Error()))
		return
	}

	sql := strings.Join(insertSqls, "\n")
	_, err = db.Exec(sql)
	if err != nil {
		log.Println(i118Utils.I118Prt.Sprintf("fail_to_exec_query", sql, err.Error()))
		return
	}

	return
}

func ReadDataFromSQLite(field model.DefField, dbName string, tableName string) ([]string, string) {
	list := make([]string, 0)

	db, err := sql.Open(constant.SqliteDriver, constant.SqliteData)
	defer db.Close()
	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_connect_sqlite", constant.SqliteData, err.Error()))
		return list, ""
	}

	selectCol := field.Select
	from := dbName
	if tableName != "" {
		from += "_" + tableName
	}

	where := field.Where
	if where == "" {
		where = "1=1"
	}
	if vari.Def.Type == constant.ConfigTypeArticle && strings.ToLower(where) == "true" {
		where = selectCol + " = 'y'"
	}

	where = strings.Replace(where, "-", "_", -1)

	if field.Rand {
		where += " ORDER BY RANDOM() "
	}

	if !strings.Contains(where, "LIMIT") {
		total := vari.Total
		if total > constant.MaxNumb { total = constant.MaxNumb }
		if field.Limit > 0 && total > field.Limit { total = field.Limit }

		where = where + fmt.Sprintf(" LIMIT %d", total)
	}

	colStr := selectCol
	if vari.Def.Type == constant.ConfigTypeArticle {
		colStr = "ci AS " + selectCol
	}

	sqlStr := fmt.Sprintf("SELECT %s FROM %s WHERE %s", colStr, from, where)
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

	if field.Select == "xingrongci_waimao_nvxing" {
		log.Println(field.Select)
	}

	return list, selectCol
}

func isExcelChanged(path string) bool {
	if !fileUtils.FileExist(path) {
		return false
	}

	fileChangeTime := time.Time{}.Unix()
	if !fileUtils.IsDir(path) {
		fileChangeTime = getFileModTime(path).Unix()
	} else {
		fileChangeTime = getDirModTime(path).Unix()
	}

	db, err := sql.Open(constant.SqliteDriver, constant.SqliteData)
	defer db.Close()
	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_connect_sqlite", constant.SqliteData, err.Error()))
		return true
	}

	sqlStr := fmt.Sprintf("SELECT id, name, changeTime FROM %s WHERE name = '%s';",
		constant.SqliteTrackTable, path)
	rows, err := db.Query(sqlStr)
	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_exec_query", sqlStr, err.Error()))
		return true
	}

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

func getDirModTime(path string) (dirChangeTime time.Time) {
	files := make([]string, 0)
	fileUtils.GetFilesByExtInDir(path, "", &files)

	for _, file := range files {
		time := getFileModTime(file)
		if dirChangeTime.Unix() < time.Unix() {
			dirChangeTime = time
		}
	}

	return
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

func importExcel(filePath, tableName string, seq *int, ddlFields, insertSqls *[]string, colMap *map[string]bool) {
	excel, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Println("fail to read file " + filePath + ", error: " + err.Error())
		return
	}

	fileName := fileUtils.GetFileName(filePath)
	fileName = strings.TrimSuffix(fileName, "词库")
	colPrefix := stringUtils.GetPinyin(fileName)

	for rowIndex, sheet := range excel.GetSheetList() {
		rows, _ := excel.GetRows(sheet)
		if len(rows) == 0 {
			continue
		}

		colDefine := ""
		colList := make([]string, 0)

		colCount := 0
		index := 0
		for colIndex, col := range rows[0] {
			val := strings.TrimSpace(col)
			if rowIndex == 0 && val == "" {
				break
			}
			colCount++

			colName := stringUtils.GetPinyin(val)

			if colIndex == 0 && colName != "ci" {
				colName = "ci"
			}
			if colName != "ci" {
				colName = colPrefix + "_" + colName
			}

			if (*colMap)[colName] == false {
				colType := "VARCHAR"
				colDefine = "    " + "`" + colName + "` " + colType + " DEFAULT ''"
				*ddlFields = append(*ddlFields, colDefine)

				(*colMap)[colName] = true
			}
			colList = append(colList, "`" + colName + "`")

			index++
		}

		valList := make([]string, 0)
		for rowIndex, row := range rows {
			if rowIndex == 0 {
				continue
			}

			valListItem := make([]string, 0)
			valListItem = append(valListItem, strconv.Itoa(*seq))
			*seq += 1

			for i := 0; i < colCount; i++ {
				val := ""
				if i == 0 { // word
					val = strings.TrimSpace(row[i])
				} else if i <= len(row) - 1 { // excel value
					val = strings.ToLower(strings.TrimSpace(row[i]))
					if val != "y" && val != "b" && val != "f" && val != "m" {
						val = ""
					}
				} else {
					val = ""
				}
				valListItem = append(valListItem,"'" + val + "'")
			}
			valList = append(valList, "(" + strings.Join(valListItem, ", ") + ")")
		}

		insertTemplate := "INSERT INTO `" + tableName + "` (`seq`, %s) VALUES %s;"
		insertSql := fmt.Sprintf(insertTemplate,
			strings.Join(colList, ", "),
			strings.Join(valList, ", "),
		)
		*insertSqls = append(*insertSqls, insertSql)
	}
}