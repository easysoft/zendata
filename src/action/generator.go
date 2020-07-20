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
	"strings"
)

func Generate(defaultFile string, configFile string, total int, fieldsToExportStr string, out string, format string, table string) {
	//startTime := time.Now().Unix()

	if defaultFile != "" && configFile == "" {
		configFile = defaultFile
		defaultFile = ""
	}

	fieldsToExport := make([]string, 0)
	if fieldsToExportStr != "" {
		fieldsToExport = strings.Split(fieldsToExportStr, ",")
	}

	constant.Total = total

	rows, colTypes := gen.GenerateForDefinition(defaultFile, configFile, &fieldsToExport, total)
	var content string
	content, vari.JsonResp = Print(rows, format, table, colTypes, fieldsToExport)

	if out != "" {
		WriteToFile(out, content)
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
			if j >0 && format == constant.FormatSql {
				line = line + ","
				valueList = valueList + ","
			}
			line = line + col

			row.Cols = append(row.Cols, col)

			colVal := col
			//colVal = stringUtils.AddPad(colVal, vari.Def.Fields[j])
			if !colTypes[j] { colVal = "'" + colVal + "'" }
			valueList = valueList + colVal
		}

		if format == constant.FormatText && i < len(rows) {
			content = content + line + "\n"
		}

		logUtils.Screen(fmt.Sprintf("%s", line))

		testData.Table.Rows = append(testData.Table.Rows, row)

		if format == constant.FormatSql {
			fieldNames := make([]string, 0)

			for _, f := range fields {
				fieldNames = append(fieldNames, "`" + f + "`")
			}
			sent := fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", table, strings.Join(fieldNames, ", "), valueList)
			sql = sql + sent + ";\n"
		}
	}

	respJson := "{}"
	if format == constant.FormatJson || vari.RunMode == constant.RunModeServer {
		mapArr := RowsToMap(rows, fields)
		jsonObj, _ := json.Marshal(mapArr)
		respJson = string(jsonObj)
	}

	if format == constant.FormatJson {
		content = respJson
	} else if format == constant.FormatJson {
		xml, _ := xml.Marshal(testData)
		content = string(xml)
	} else if format == constant.FormatSql {
		content = sql
	}

	return content, respJson
}

func RowsToMap(rows [][]string, fieldsToExport []string) (ret []map[string]string) {
	ret = []map[string]string{}

	for _, cols := range rows {
		rowMap := map[string]string{}
		for j, col := range cols {
			rowMap[fieldsToExport[j]] = col
		}

		ret = append(ret, rowMap)
	}
	return
}