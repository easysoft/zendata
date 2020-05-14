package gen

import (
	"database/sql"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"regexp"
	"strconv"
	"strings"
)

func GenerateFieldValuesFromExcel(field *model.Field, fieldValue *model.FieldValue, level int) {
	// get file and step string
	rang := strings.TrimSpace(field.Range)
	sectionArr := strings.Split(rang, ":")
	file := sectionArr[0]
	stepStr := "1"
	if len(sectionArr) == 2 {
		stepStr = sectionArr[1]
	}

	list := make([]string, 0)
	path := constant.DataDir + file
	ConvertExcelToSQLite(*field, path)

	list = ReadDataSQLite(*field)

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
	numbs := GenerateIntItems(0, (int64)(len(list)-1), step, rand)
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

	for _, sheet := range excel.GetSheetList() {
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
				if colIndex > 0 {
					valList = valList + ", "
				}
				valList = valList + "'" + colCell + "'"
			}
			valList = valList + ")"
		}

		tableName := field.Name + "_" + sheet
		dropSql := fmt.Sprintf(dropTemplate, tableName)
		ddl := fmt.Sprintf(ddlTemplate, tableName, colDefine)
		insertSql := fmt.Sprintf(insertTemplate, tableName, colList, valList)

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
}

func ReadDataSQLite(field model.Field) []string {
	list := make([]string, 0)

	db, err := sql.Open(constant.SqliteDriver, constant.SqliteSource)
	if err != nil {
		logUtils.Screen("fail to open " + constant.SqliteSource + ": " + err.Error())
		return list
	}
	field.Filter = replaceDotInTableNameLimit(field.Filter)

	rows, err := db.Query(field.Filter)
	if err != nil {
		logUtils.Screen("fail to exec query " + err.Error())
		return list
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
			logUtils.Screen("fail to get sqlite3 row: " + err.Error())
			return list
		}

		rowMap := map[string]string{}
		for index, v := range values {
			item := v.(*string)

			rowMap[colIndexToName[index]] = *item
		}

		valMapArr = append(valMapArr, rowMap)
	}

	format := field.Format
	for _, item := range valMapArr {
		line := replacePlaceholderWithValue(format,item)
		list = append(list, line)
	}

	return list
}

func replacePlaceholderWithValue(format string, valMap map[string]string) string {
	// ${user_name}_${numb}@${domain}
	regx := regexp.MustCompile(`\$\{([a-zA-z0-9_]+)\}`)
	arrOfName := regx.FindAllStringSubmatch(format, -1)

	ret := ""
	if len(arrOfName) > 0 {
		strLeft := format
		for index, a := range arrOfName {
			found := a[0]
			name := a[1]

			arr := strings.Split(strLeft, found)

			// add string constant
			if arr[0] != "" {
				ret = ret + arr[0]
			}

			ret = ret + valMap[name]

			arr = arr[1:]
			strLeft = strings.Join(arr, "")

			if index == len(arrOfName) - 1 && strLeft != "" { // add last item in arr
				ret = ret + strLeft
			}
		}
	}

	return ret
}

func replaceDotInTableNameLimit(str string) string {
	ret := ""

	str = strings.Replace(str," from ", " FROM ", -1)
	str = strings.Replace(str," where ", " WHERE ", -1)
	str = strings.Replace(str," limit ", " LIMIT ", -1)

	arr1 := strings.Split(str, " FROM ")
	arr2 := strings.Split(arr1[1], " WHERE ")

	ret = arr1[0] + " FROM " + strings.Replace(arr2[0],".", "_", -1) + " WHERE " + arr2[1]
	if !strings.Contains(ret, "LIMIT") {
		ret = ret + " LIMIT " + strconv.Itoa(constant.MaxNumb)
	}

	return ret
}