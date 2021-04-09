package gen

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zendata/src/gen/helper"
	constant "github.com/easysoft/zendata/src/utils/const"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/mattn/go-runewidth"
	"regexp"
	"strings"
)

func Print(rows [][]string, format string, table string, colIsNumArr []bool,
	fields []string) (lines []interface{}) {

	if format == constant.FormatText {
		printTextHeader(fields)
	} else if format == constant.FormatSql {
		sqlHeader := getInsertSqlHeader(fields, table)
		if vari.DBDsn != "" {
			lines = append(lines, sqlHeader)
		} else {
			logUtils.PrintLine(sqlHeader)
		}
	} else if format == constant.FormatJson {
		printJsonHeader()
	} else if format == constant.FormatXml {
		printXmlHeader(fields, table)
	}

	for i, cols := range rows {
		row := make([]string, 0)
		rowMap := map[string]string{}
		valuesForSql := make([]string, 0)
		lineForText := ""

		for j, col := range cols {
			// 3. random replacement
			col = replacePlaceholder(col)
			field := vari.TopFieldMap[fields[j]]
			if field.Length > runewidth.StringWidth(col) {
				//col = stringUtils.AddPad(col, field)
			}

			if j > 0 && vari.Human { // use a tab
				lineForText = strings.TrimRight(lineForText, "\t")
				col = strings.TrimLeft(col, "\t")

				lineForText = lineForText + "\t" + col
			} else {
				lineForText = lineForText + col
			}

			row = append(row, col)
			rowMap[field.Field] = col

			colVal := stringUtils.ConvertForSql(col)
			if !colIsNumArr[j] {
				colVal = "'" + colVal + "'"
			}
			valuesForSql = append(valuesForSql, colVal)
		}

		if format == constant.FormatText && vari.Def.Type == constant.ConfigTypeArticle { // article need to write to more than one files
			lines = append(lines, lineForText)

		} else if format == constant.FormatText && vari.Def.Type != constant.ConfigTypeArticle {
			logUtils.PrintLine(lineForText)

		} else if format == constant.FormatSql {
			sql := genSqlLine(strings.Join(valuesForSql, ", "), i, len(rows))

			if vari.DBDsn != "" { // add to return array for exec sql
				lines = append(lines, sql)
			} else {
				logUtils.PrintLine(sql)
			}

		} else if format == constant.FormatJson {
			logUtils.PrintLine(genJsonLine(i, row, len(rows), fields))
		} else if format == constant.FormatXml {
			logUtils.PrintLine(getXmlLine(i, rowMap, len(rows)))
		} else if format == constant.FormatData {
			lines = append(lines, lineForText)
		}
	}

	return
}

func printTextHeader(fields []string) {
	if !vari.WithHead {
		return
	}
	headerLine := ""
	for idx, field := range fields {
		headerLine += field
		if idx < len(fields)-1 {
			headerLine += "\t"
		}
	}

	logUtils.PrintLine(headerLine)
}

func getInsertSqlHeader(fields []string, table string) string {
	fieldNames := make([]string, 0)
	for _, f := range fields {
		if vari.Server == "mysql" {
			f = "`" + f + "`"
		}
		fieldNames = append(fieldNames, f)
	}
	ret := fmt.Sprintf("INSERT INTO %s(%s)", table, strings.Join(fieldNames, ", "))
	return ret
}

func printJsonHeader() {
	logUtils.PrintLine("[")
}

func printXmlHeader(fields []string, table string) {
	logUtils.PrintLine("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<testdata>\n  <title>Test Data</title>")
}

func RowToJson(cols []string, fieldsToExport []string) string {
	rowMap := map[string]string{}
	for j, col := range cols {
		rowMap[fieldsToExport[j]] = col
	}

	jsonObj, _ := json.Marshal(rowMap)
	respJson := string(jsonObj)

	return respJson
}

func genSqlLine(valuesForSql string, i int, length int) string {

	temp := ""
	if i == 0 {
		temp = fmt.Sprintf("  VALUES (%s)", valuesForSql)
	} else {
		temp = fmt.Sprintf("         (%s)", valuesForSql)
	}

	if i < length-1 {
		temp = temp + ", "
	} else {
		temp = temp + "; "
	}

	return temp
}

func genJsonLine(i int, row []string, length int, fields []string) string {
	temp := "  " + RowToJson(row, fields)
	if i < length-1 {
		temp = temp + ", "
	} else {
		temp = temp + "\n]"
	}

	return temp
}

func getXmlLine(i int, mp map[string]string, length int) string {
	str := ""
	j := 0
	for key, val := range mp {
		str += fmt.Sprintf("    <%s>%s</%s>", key, val, key)
		if j != len(mp)-1 {
			str = str + "\n"
		}

		j++
	}

	text := fmt.Sprintf("  <row>\n%s\n  </row>", str)
	if i == length-1 {
		text = text + "\n</testdata>"
	}
	return text
}

func replacePlaceholder(col string) string {
	ret := col

	re := regexp.MustCompile("(?siU)\\${(.*)}")
	matchResultArr := re.FindAllStringSubmatch(col, -1)
	matchTimes := len(matchResultArr)

	for _, childArr := range matchResultArr {
		placeholderStr := childArr[1]
		values := getValForPlaceholder(placeholderStr, matchTimes)

		for _, str := range values {
			temp := Placeholder(placeholderStr)
			ret = strings.Replace(ret, temp, str, 1)
		}
	}

	return ret
}

func getValForPlaceholder(placeholderStr string, count int) []string {
	mp := vari.RandFieldNameToValuesMap[placeholderStr]

	tp := mp["type"].(string)
	repeatObj := mp["repeat"]

	repeat := 1
	if repeatObj != nil {
		repeat = repeatObj.(int)
	}

	strs := make([]string, 0)
	repeatTag := mp["repeatTag"].(string)
	if tp == "int" {
		start := mp["start"].(string)
		end := mp["end"].(string)
		precision := mp["precision"].(string)
		format := mp["format"].(string)

		strs = helper.GetRandFromRange("int", start, end, "1", repeat, repeatTag, precision, count, format)
	} else if tp == "float" {
		start := mp["start"].(string)
		end := mp["end"].(string)
		precision := mp["precision"].(string)
		format := mp["format"].(string)

		strs = helper.GetRandFromRange("float", start, end, "1", repeat, repeatTag, precision, count, format)
	} else if tp == "char" {
		start := mp["start"].(string)
		end := mp["end"].(string)
		precision := mp["precision"].(string)
		format := mp["format"].(string)

		strs = helper.GetRandFromRange("char", start, end, "1", repeat, repeatTag, precision, count, format)
	} else if tp == "list" {
		list := mp["list"].([]string)
		strs = helper.GetRandFromList(list, repeat, count)
	}

	strs = strs[:count]
	return strs
}
