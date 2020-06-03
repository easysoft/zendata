package action

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/easysoft/zendata/src/gen"
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"github.com/easysoft/zendata/src/utils/vari"
	"os"
	"path/filepath"
	"strings"
)

func Generate(deflt string, yml string, total int, fieldsToExportStr string, out string, format string, table string) {
	//startTime := time.Now().Unix()

	if deflt != "" && yml == "" {
		yml = deflt
		deflt = ""
	}

	fieldsToExport := make([]string, 0)
	if fieldsToExportStr != "" {
		fieldsToExport = strings.Split(fieldsToExportStr, ",")
	}

	vari.InputDir = filepath.Dir(yml) + string(os.PathSeparator)
	constant.Total = total

	rows, colTypes := gen.GenerateForDefinition(deflt, yml, fieldsToExport, total)
	content := Print(rows, format, table, colTypes, fieldsToExport)

	if out != "" {
		WriteToFile(out, content)
	}

	//entTime := time.Now().Unix()
	//logUtils.Screen(i118Utils.I118Prt.Sprintf("generate_records", len(rows), out, entTime - startTime ))
}

func Print(rows [][]string, format string, table string, colTypes []bool, fields []string) string {
	content := ""
	sql := ""

	testData := model.TestData{}
	testData.Title = "测试数据"

	for i, cols := range rows {
		line := ""
		row := model.Row{}
		valueList := ""

		for j, col := range cols {
			if j >0 && format == "sql" {
				line = line + ","
				valueList = valueList + ","
			}
			line = line + col

			row.Cols = append(row.Cols, col)

			colVal := col
			if !colTypes[j] { colVal = "'" + colVal + "'" }
			valueList = valueList + colVal
		}

		if format == "text" && i < len(rows) {
			content = content + line + "\n"
		}

		logUtils.Screen(fmt.Sprintf("%s", line))

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