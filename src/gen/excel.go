package gen

import (
	"database/sql"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"strconv"
	"strings"
)

func GenerateFieldValuesFromExcel(field *model.Field, fieldValue *model.FieldValue, level int) {
	// get file and step string
	rang := strings.TrimSpace(field.Range)
	sectionArr := strings.Split(rang, ":")
	file := sectionArr[0]
	stepStr := "1"
	if len(sectionArr) == 2 { stepStr = sectionArr[1] }

	list := make([]string, 0)
	path := constant.DataDir + file
	ConvertExcelToSQLite(*field, path)

	list = strings.Split("str" + path, "\n")

	// get step and rand
	rand := false
	step := 1
	if strings.ToLower(strings.TrimSpace(stepStr)) != "r" {
		stepInt, err := strconv.Atoi(stepStr)
		if err == nil {
			step = stepInt
		}
	} else {
		rand = true
	}

	// get index for data retrieve
	numbs := GenerateIntItems(0, (int64)(len(list) - 1), step, rand)
	// get data by index
	index := 0
	for _, numb := range numbs {
		item := list[numb.(int64)]

		if index >= constant.MaxNumb { break }
		if strings.TrimSpace(item) == "" { continue }

		fieldValue.Values = append(fieldValue.Values, item)
		index = index + 1
	}

	if len(fieldValue.Values) == 0 {
		fieldValue.Values = append(fieldValue.Values, "N/A")
	}
}

func ConvertExcelToSQLite(field model.Field, path string) {
	excel, err := excelize.OpenFile(path)
	if err != nil {
		logUtils.Screen("fail to read file: " + err.Error())
		return
	}

	sheet := excel.GetSheetName(0);
	rows, err := excel.GetRows(sheet)

	dropTemplate := `DROP TABLE IF EXISTS %s;`
	ddlTemplate := `CREATE TABLE %s (
						%s
					);`
	insertTemplate := "INSERT INTO %s (%s) VALUES %s"

	colDefine := ""
	colList := ""
	index := 0
	for _, col := range rows[0] {
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
		colDefine = "    " +  colDefine + val + " " + colProp

		colList = colList + val
		index++
	}

	valList := ""
	for rowIndex, row := range rows {
		if rowIndex == 0 { continue }

		if rowIndex > 1 { valList = valList + ", " }
		valList = valList + "("

		for colIndex, colCell := range row {
			if colIndex > 0 { valList = valList + ", " }
			valList = valList + "'" + colCell + "'"
		}
		valList = valList + ")"
	}

	dropSql := fmt.Sprintf(dropTemplate, field.Name)
	ddl := fmt.Sprintf(ddlTemplate, field.Name, colDefine)
	insertSql := fmt.Sprintf(insertTemplate, field.Name, colList, valList)

	db, err := sql.Open("sqlite3", constant.SqliteSource)
	_, err = db.Exec(dropSql)
	_, err = db.Exec(ddl)
	if err != nil {
		logUtils.Screen("fail to create table: " + err.Error())
		return
	} else {
		_, err = db.Exec(insertSql)
		if err != nil {
			logUtils.Screen("fail to insert data: " + err.Error())
			return
		}
	}
}

func ReadDataSQLite(table string) []string {
	list := make([]string, 0)



	return list
}