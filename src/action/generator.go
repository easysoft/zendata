package action

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/easysoft/zendata/src/gen"
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"github.com/easysoft/zendata/src/utils/vari"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func Generate(def string, total int, fieldsToExportStr string, out string, format string, table string) {
	startTime := time.Now().Unix()
	vari.InputDir = filepath.Dir(def) + string(os.PathSeparator)
	constant.Total = total

	fieldsToExport := strings.Split(fieldsToExportStr, ",")

	referFieldValueMap := gen.LoadClsDef(def, fieldsToExport)
	referFieldValueMap = referFieldValueMap

	rows, colTypes := gen.GenerateForDefinition(total, fieldsToExport)
	content := Print(rows, format, table, colTypes, fieldsToExport)

	WriteToFile(out, content)

	entTime := time.Now().Unix()
	logUtils.Screen(fmt.Sprintf("Genereate %d records, spend %d secs",
		len(rows), entTime - startTime ))
}

func Print(rows [][]string, format string, table string, colTypes []bool, fields []string) string {
	width := stringUtils.GetNumbWidth(len(rows))

	content := ""
	sql := ""

	testData := model.TestData{}
	testData.Title = "测试数据"

	for i, cols := range rows {
		line := ""
		row := model.Row{}
		valueList := ""

		for j, col := range cols {
			if j >0 {
				line = line + ", "
				valueList = valueList + ", "
			}
			line = line + col

			row.Cols = append(row.Cols, col)

			colVal := col
			if !colTypes[j] { colVal = "'" + colVal + "'" }
			valueList = valueList + colVal
		}

		if format == "text" && i < len(rows) - 1 { content = content + line + "\n" }

		idStr := fmt.Sprintf("%" + strconv.Itoa(width) + "d", i+1)
		logUtils.Screen(fmt.Sprintf("%s: %s", idStr, line))

		testData.Table.Rows = append(testData.Table.Rows, row)

		if format == "sql" {
			sent := fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", table, strings.Join(fields, ","), valueList)
			sql = sql + sent + ";\n"
		}
	}

	if format == "json" {
		json, _ := json.Marshal(rows)
		content = string(json)
	} else if format == "xml" {
		xml, _ := xml.Marshal(testData)
		content = string(xml)
	} else if format == "sql" {
		content = sql
	}

	return content
}