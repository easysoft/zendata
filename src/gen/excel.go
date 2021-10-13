package gen

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/easysoft/zendata/src/gen/helper"
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/fatih/color"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func generateFieldValuesFromExcel(filePath, sheet string, field *model.DefField, total int) (values map[string][]string) {
	values = map[string][]string{}

	// sql has variable expr
	if filePath == "" || helper.SelectExcelWithExpr(*field) {
		return
	}

	dbName := getDbName(filePath)

	if !fileUtils.IsDir(filePath) { // file
		firstSheet := ConvertSingleExcelToSQLiteIfNeeded(dbName, filePath)
		if sheet == "" {
			sheet = firstSheet
		}
	} else { // dir, for article generation only
		ConvertWordExcelsToSQLiteIfNeeded(dbName, filePath)
	}

	list, fieldSelect := ReadDataFromSQLite(*field, dbName, sheet, total, filePath)
	// get index for data retrieve
	numbs := GenerateIntItems(0, (int64)(len(list)-1), 1, false, 1, "")
	// get data by index
	index := 0
	for _, numb := range numbs {
		item := list[numb.(int64)%(int64(len(list)))]

		if index >= constant.MaxNumb {
			break
		}

		values[fieldSelect] = append(values[fieldSelect], item)
		index = index + 1
	}

	return
}

func getDbName(path string) (dbName string) {
	dbName = strings.Replace(path, vari.ZdPath+constant.ResDirData+constant.PthSep, "", -1)
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

		_, err = vari.DB.Exec(dropSql)
		if err != nil {
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_drop_table", tableName, err.Error()))
			return
		}

		_, err = vari.DB.Exec(ddl)
		if err != nil {
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_create_table", tableName, err.Error()))
			return
		}

		_, err = vari.DB.Exec(insertSql)
		if err != nil {
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_exec_query", insertSql, err.Error()))
			return
		}
	}

	return
}

func ConvertWordExcelsToSQLiteIfNeeded(tableName string, dir string) {
	if !isExcelChanged(dir) {
		return
	}

	files := make([]string, 0)
	fileUtils.GetFilesByExtInDir(dir, ".xlsx", &files)

	seq := 1
	insertSqls := make([]string, 0)
	ddlFields := make([]string, 0)
	ddlFields = append(ddlFields, "    `词语` VARCHAR DEFAULT ''")

	colMap := map[string]bool{}
	for _, file := range files {
		importExcel(file, tableName, &seq, &ddlFields, &insertSqls, &colMap)
	}

	dropSql := `DROP TABLE IF EXISTS ` + tableName + `;`
	_, err := vari.DB.Exec(dropSql)
	if err != nil {
		log.Println(i118Utils.I118Prt.Sprintf("fail_to_drop_table", tableName, err.Error()))
		return
	}

	ddlTemplate := "CREATE TABLE " + tableName + "(\n" +
		"\t`seq` CHAR (5) PRIMARY KEY ASC UNIQUE,\n" +
		"%s" +
		"\n);"
	ddlSql := fmt.Sprintf(ddlTemplate, strings.Join(ddlFields, ", \n"))
	_, err = vari.DB.Exec(ddlSql)
	if err != nil {
		log.Println(i118Utils.I118Prt.Sprintf("fail_to_create_table", tableName, err.Error()))
		return
	}

	sql := strings.Join(insertSqls, "\n")
	_, err = vari.DB.Exec(sql)
	if err != nil {
		log.Println(i118Utils.I118Prt.Sprintf("fail_to_exec_query", sql, err.Error()))
		return
	}

	return
}

func ReadDataFromSQLite(field model.DefField, dbName string, tableName string, total int, filePath string) (
	[]string, string) {
	list := make([]string, 0)

	fieldSelect := field.Select
	from := dbName
	if tableName != "" {
		from += "_" + tableName
	}

	where := strings.TrimSpace(field.Where)
	if vari.Def.Type == constant.ConfigTypeArticle {
		if where == "" {
			where = "y"
		}

		cols := strings.Split(fieldSelect, "-")
		wheres := ""
		for index, col := range cols {
			if index == 0 {
				wheres += fmt.Sprintf("`%s` = '%s'", col, "y")
			} else {
				wheres += " AND "
				wheres += fmt.Sprintf("`%s` = '%s'", col, where)
			}
		}

		where = wheres

	} else {
		if where == "" {
			where = "1=1"
		}
	}

	if field.Rand {
		where += " ORDER BY RANDOM() "
	}

	if !strings.Contains(where, "LIMIT") {
		if total > constant.MaxNumb {
			total = constant.MaxNumb
		}
		if field.Limit > 0 && total > field.Limit {
			total = field.Limit
		}

		where = where + fmt.Sprintf(" LIMIT %d", total)
	}

	colStr := fieldSelect
	if vari.Def.Type == constant.ConfigTypeArticle {
		colStr = "`词语` AS `" + fieldSelect + "`"
	}

	sqlStr := fmt.Sprintf("SELECT %s FROM `%s` WHERE %s", colStr, from, where)
	rows, err := vari.DB.Query(sqlStr)
	defer rows.Close()
	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_exec_query", sqlStr, err.Error()))
		logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("pls_check_excel", filePath), color.FgRed)

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
			if idx > 0 {
				break
			}
			list = append(list, val)
			idx++
		}
	}

	return list, fieldSelect
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

	sqlStr := fmt.Sprintf("SELECT id, name, changeTime FROM %s WHERE name = '%s' ORDER BY changeTime DESC LIMIT 1;",
		constant.SqliteTrackTable, path)
	rows, err := vari.DB.Query(sqlStr)
	defer rows.Close()
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

	if !found {
		changed = true
	}

	if changed {
		if !found {
			sqlStr = fmt.Sprintf("INSERT INTO %s(name, changeTime) VALUES('%s', %d)",
				constant.SqliteTrackTable, path, fileChangeTime)
		} else {
			sqlStr = fmt.Sprintf("UPDATE %s SET changeTime = %d WHERE name = '%s'",
				constant.SqliteTrackTable, fileChangeTime, path)
		}

		_, err = vari.DB.Exec(sqlStr)
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

	colPrefix := fileName // stringUtils.GetPinyin(fileName)
	*ddlFields = append(*ddlFields, "    `"+colPrefix+"` VARCHAR DEFAULT ''")

	for sheetIndex, sheet := range excel.GetSheetList() {
		rows, _ := excel.GetRows(sheet)
		if len(rows) == 0 {
			continue
		}

		colDefine := ""
		colList := make([]string, 0)

		colCount := 0
		index := 0
		// gen col list for ddl and insert cols
		for colIndex, col := range rows[0] {
			val := strings.TrimSpace(col)

			if sheetIndex == 0 && val == "" {
				break
			}
			colCount++

			colList = append(colList, val)

			colNames := val
			colNameArr := strings.Split(colNames, "-")
			for _, colName := range colNameArr {
				if (*colMap)[colName] == false {
					colType := "VARCHAR"
					colDefine = "    " + "`" + colName + "` " + colType + " DEFAULT ''"

					if colIndex == 0 {
						colName = "词语"
					} else { // first already added
						*ddlFields = append(*ddlFields, colDefine)
					}

					(*colMap)[colName] = true
				}
			}

			index++
		}

		insertTemplate := "INSERT INTO `" + tableName + "` (%s) VALUES (%s);"
		// gen values for insert
		for rowIndex, row := range rows {
			if rowIndex == 0 { // ignore title line
				continue
			}

			record := map[string]interface{}{}
			record[colPrefix] = "y"
			record["seq"] = *seq
			*seq += 1

			for colIndex, col := range row {
				if colIndex >= len(colList) {
					break
				}

				colNames := colList[colIndex]

				val := strings.ToLower(strings.TrimSpace(col))
				if val == "" {
					continue
				}

				if colIndex == 0 { // word
					record["词语"] = val
				} else {
					if val == "y" || val == "b" || val == "f" || val == "m" {
						val = "y"
					} else {
						val = ""
					}

					colNameArr := strings.Split(colNames, "-")
					for _, colName := range colNameArr {
						record[colName] = val
					}
				}
			}

			cols := make([]string, 0)
			vals := make([]string, 0)

			for key, val := range record {
				cols = append(cols, "`"+key+"`")

				valStr := ""
				switch val.(type) {
				case int:
					valStr = strconv.Itoa(val.(int))
				default:
					valStr = "'" + val.(string) + "'"
				}

				vals = append(vals, valStr)
			}

			insertSql := fmt.Sprintf(insertTemplate, strings.Join(cols, ","), strings.Join(vals, ","))
			*insertSqls = append(*insertSqls, insertSql)
		}
	}
}
