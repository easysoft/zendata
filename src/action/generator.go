package action

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/easysoft/zendata/src/gen"
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/fatih/color"
	"net/http"
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

	rows, colTypes := gen.GenerateForDefinition(deflt, yml, &fieldsToExport, total)
	var content string
	content, vari.JsonResp = Print(rows, format, table, colTypes, fieldsToExport)

	if out != "" {
		WriteToFile(out, content)
	}

	if out != "" {
		WriteToFile(out, content)
	}

	if vari.HttpService {
		logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("press_to_exist"), color.FgCyan)

		http.HandleFunc("/", DataHandler)
		http.HandleFunc("/data", DataHandler)
		http.ListenAndServe(":58848", nil)
	}

	//entTime := time.Now().Unix()
	//logUtils.Screen(i118Utils.I118Prt.Sprintf("generate_records", len(rows), out, entTime - startTime ))
}

func Print(rows [][]string, format string, table string, colTypes []bool, fields []string) (string, string) {
	content := ""
	sql := ""

	if vari.WithHead {
		line := ""
		for idx, field := range fields {
			line += field
			if idx < len(fields) - 1 {
				line += vari.HeadSep
			}
		}
		logUtils.Screen(fmt.Sprintf("%s", line))
		content += line + "\n"
	}

	testData := model.TestData{}
	testData.Title = "Test Data"

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
			fieldNames := make([]string, 0)

			for _, f := range fields {
				fieldNames = append(fieldNames, "`" + f + "`")
			}
			sent := fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", table, strings.Join(fieldNames, ", "), valueList)
			sql = sql + sent + ";\n"
		}
	}

	respJson := "[]"
	if format == "json" || vari.HttpService {
		if vari.WithHead {
			mapObj := RowsToMap(rows, fields)
			jsonObj, _ := json.Marshal(mapObj)
			respJson = string(jsonObj)
		} else {
			jsonObj, _ := json.Marshal(rows)
			respJson = string(jsonObj)
		}
	}

	if format == "json" {
		content = respJson
	} else if format == "xml" {
		xml, _ := xml.Marshal(testData)
		content = string(xml)
	} else if format == "sql" {
		content = sql
	}

	return content, respJson
}

func RowsToMap(rows [][]string, fieldsToExport []string) (ret map[string]string) {
	ret = map[string]string{}

	for _, cols := range rows {
		for j, col := range cols {
			ret[fieldsToExport[j]] = col
		}
	}
	return
}

func DataHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, vari.JsonResp)
}