package service

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/easysoft/zendata/internal/pkg/domain"
	"github.com/easysoft/zendata/internal/pkg/helper"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/fatih/color"
	"gorm.io/gorm"
)

type ExcelService struct {
	ExpressionService *ExpressionService `inject:""`
}

func (s *ExcelService) generateFieldValuesFromExcel(filePath, sheet string, field *domain.DefField, total int) (
	values map[string][]interface{}) {
	values = map[string][]interface{}{}

	// sql has variable expr
	if filePath == "" || helper.IsSelectExcelWithExpr(*field) {
		return
	}

	dbName := s.getDbName(filePath)

	if !fileUtils.IsDir(filePath) { // file
		firstSheet := s.ConvertSingleExcelToSQLiteIfNeeded(dbName, filePath)
		if sheet == "" {
			sheet = firstSheet
		}
	} else { // dir, for article generation only
		s.ConvertWordExcelsToSQLiteIfNeeded(dbName, filePath)
	}

	list, fieldSelect := s.ReadDataFromSQLite(*field, dbName, sheet, total, filePath)
	// get index list for data retrieve
	numbs := helper.GenerateItems(int64(0), int64(len(list)-1), int64(1), 0, false, 1, "", 0)

	// get data by index
	index := 0
	for _, numb := range numbs {
		item := list[numb.(int64)%(int64(len(list)))]

		if index >= consts.MaxNumb {
			break
		}

		values[fieldSelect] = append(values[fieldSelect], item)
		index = index + 1
	}

	return
}

func (s *ExcelService) getDbName(path string) (dbName string) {
	dbName = strings.Replace(path, vari.WorkDir+consts.ResDirData+consts.PthSep, "", -1)
	dbName = strings.Replace(dbName, consts.PthSep, "_", -1)
	dbName = strings.Replace(dbName, ".", "_", -1)

	return
}

func (s *ExcelService) ConvertSingleExcelToSQLiteIfNeeded(dbName string, path string) (firstSheet string) {
	excel, err := excelize.OpenFile(path)
	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_read_file", path))
		return
	}

	firstSheet = excel.GetSheetList()[0]
	changed, sqlBeforeCompleted := s.isExcelChanged(path)
	if !changed {
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

			dataColCount := 0
			for colIndex, colCell := range row {
				if colIndex >= colCount {
					break
				}

				if colIndex > 0 {
					valList = valList + ", "
				}
				colCell = strings.Replace(colCell, "'", "''", -1)
				valList = valList + "'" + colCell + "'"

				dataColCount++
			}

			for dataColCount < colCount {
				valList = valList + ", ''"
				dataColCount++
			}

			valList = valList + ")"
		}

		tableName := dbName + "_" + sheet
		dropSql := fmt.Sprintf(dropTemplate, tableName)
		ddl := fmt.Sprintf(ddlTemplate, tableName, colDefine)
		insertSql := fmt.Sprintf(insertTemplate, tableName, colList, valList)

		err = vari.DB.Exec(dropSql).Error
		if err != nil {
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_drop_table", tableName, err.Error()))
			return
		}

		err = vari.DB.Exec(ddl).Error
		if err != nil {
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_create_table", tableName, err.Error()))
			return
		}

		err = vari.DB.Exec(insertSql).Error
		if err != nil {
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_exec_query", insertSql, err.Error()))
			return
		}

		if changed {
			err = vari.DB.Exec(sqlBeforeCompleted).Error
			if err != nil {
				logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_exec_query", sqlBeforeCompleted, err.Error()))
			}
		}
	}

	return
}

func (s *ExcelService) ConvertWordExcelsToSQLiteIfNeeded(tableName string, dir string) {
	changed, sqlBeforeCompleted := s.isExcelChanged(dir)
	if !changed {
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
		s.genImportExcelSqls(file, tableName, &seq, &ddlFields, &insertSqls, &colMap)
	}

	dropSql := `DROP TABLE IF EXISTS ` + tableName + `;`
	err := vari.DB.Exec(dropSql).Error
	if err != nil {
		log.Println(i118Utils.I118Prt.Sprintf("fail_to_drop_table", tableName, err.Error()))
		return
	}

	ddlTemplate := "CREATE TABLE " + tableName + "(\n" +
		"\t`seq` CHAR (5) PRIMARY KEY ASC UNIQUE,\n" +
		"%s" +
		"\n);"
	ddlSql := fmt.Sprintf(ddlTemplate, strings.Join(ddlFields, ", \n"))
	err = vari.DB.Exec(ddlSql).Error
	if err != nil {
		log.Println(i118Utils.I118Prt.Sprintf("fail_to_create_table", tableName, err.Error()))
		return
	}

	insertSql := strings.Join(insertSqls, "\n")
	err = vari.DB.Exec(insertSql).Error
	if err != nil {
		log.Println(i118Utils.I118Prt.Sprintf("fail_to_exec_query", insertSql, err.Error()))
		return
	}

	if changed {
		err = vari.DB.Exec(sqlBeforeCompleted).Error
		if err != nil {
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_exec_query", sqlBeforeCompleted, err.Error()))
		}
	}

	return
}

func (s *ExcelService) ReadDataFromSQLite(field domain.DefField, dbName string, tableName string, total int, filePath string) (
	[]string, string) {
	list := make([]string, 0)

	fieldSelect := field.Select
	from := dbName
	if tableName != "" {
		from += "_" + tableName
	}

	where := strings.TrimSpace(field.Where)
	if vari.GlobalVars.DefData.Type == consts.DefTypeArticle {
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
		if total > consts.MaxNumb {
			total = consts.MaxNumb
		}
		if field.Limit > 0 && total > field.Limit {
			total = field.Limit
		}

		where = where + fmt.Sprintf(" LIMIT %d", total)
	}

	colStr := fieldSelect
	if vari.GlobalVars.DefData.Type == consts.DefTypeArticle {
		colStr = "`词语` AS `" + fieldSelect + "`"
	}

	sqlStr := fmt.Sprintf("SELECT %s FROM `%s` WHERE %s", colStr, from, where)
	rows, err := vari.DB.Raw(sqlStr).Rows()

	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_exec_query", "", err.Error()))
		logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("pls_check_excel", filePath), color.FgRed)

		return list, ""
	}

	defer rows.Close()
	valMapArr := make([]map[string]string, 0)
	columns, err := rows.Columns()
	colNum := len(columns)

	colIndexToName := map[int]string{}
	for index, col := range columns {
		colIndexToName[index] = col
	}

	var values = make([]interface{}, colNum)
	for i := range values {
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

func (s *ExcelService) genExcelValuesWithExpr(field *domain.DefField, fieldNameToValuesMap map[string][]interface{}) (
	values []interface{}) {
	selects := s.ExpressionService.ReplaceVariableValues(field.Select, fieldNameToValuesMap)
	wheres := s.ExpressionService.ReplaceVariableValues(field.Where, fieldNameToValuesMap)

	childMapValues := make([][]interface{}, 0)
	for index, slct := range selects {
		temp := *field
		temp.Select = slct
		temp.Where = wheres[index%len(wheres)]

		resFile, _, sheet := fileUtils.GetResProp(temp.From, temp.FileDir)

		selectCount := vari.GlobalVars.Total/len(selects) + 1
		mp := s.generateFieldValuesFromExcel(resFile, sheet, &temp, selectCount) // re-generate values
		for _, items := range mp {
			childMapValues = append(childMapValues, items)
		}
	}
	for index := 0; len(values) < vari.GlobalVars.Total; {
		for i := range selects {
			values = append(values, childMapValues[i][index%len(childMapValues[i])])
		}
		index++
	}

	return
}

type ExcelChangedResult struct {
	Id         uint
	Path       string
	ChangeTime int64 `gorm:"column:changeTime"`
}

func (s *ExcelService) isExcelChanged(path string) (changed bool, sqlBeforeCompleted string) {
	if !fileUtils.FileExist(path) {
		return
	}

	fileChangeTime := time.Time{}.Unix()
	if !fileUtils.IsDir(path) {
		fileChangeTime = s.getFileModTime(path).Unix()
	} else {
		fileChangeTime = s.getDirModTime(path).Unix()
	}

	sqlStr := fmt.Sprintf("SELECT id, path, changeTime FROM %s "+
		"WHERE path = '%s' "+
		"ORDER BY changeTime DESC "+
		"LIMIT 1;",
		consts.SqliteTrackTable, path)

	record := ExcelChangedResult{}
	err := vari.DB.Raw(sqlStr).Scan(&record).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_exec_query", sqlStr, err.Error()))

		changed = true
		return
	}

	found := false

	if record.Id > 0 { // found
		found = true

		if path == record.Path && record.ChangeTime < fileChangeTime { // update exist record
			changed = true
		}
	} else { // not found, to add a record
		changed = true
	}

	if changed {
		if !found {
			sqlBeforeCompleted = fmt.Sprintf("INSERT INTO %s(path, changeTime) VALUES('%s', %d)",
				consts.SqliteTrackTable, path, fileChangeTime)
		} else {
			sqlBeforeCompleted = fmt.Sprintf("UPDATE %s SET changeTime = %d WHERE path = '%s'",
				consts.SqliteTrackTable, fileChangeTime, path)
		}
	}

	return
}

func (s *ExcelService) getDirModTime(path string) (dirChangeTime time.Time) {
	files := make([]string, 0)
	fileUtils.GetFilesByExtInDir(path, "", &files)

	for _, file := range files {
		time := s.getFileModTime(file)
		if dirChangeTime.Unix() < time.Unix() {
			dirChangeTime = time
		}
	}

	return
}

func (s *ExcelService) getFileModTime(path string) time.Time {
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

func (s *ExcelService) genImportExcelSqls(filePath, tableName string, seq *int, ddlFields, insertSqls *[]string, colMap *map[string]bool) {
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
