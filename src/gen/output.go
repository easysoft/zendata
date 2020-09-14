package gen

import (
	"encoding/json"
	"fmt"
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
		printSqlHeader(fields, table)
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
			col = replacePlaceholder(col)
			field := vari.TopFieldMap[fields[j]]
			if field.Width > runewidth.StringWidth(col) {
				col = stringUtils.AddPad(col, field)
			}

			if j > 0 && vari.Human {
				lineForText = strings.TrimRight(lineForText, "\t")
				col = strings.TrimLeft(col, "\t")

				lineForText = lineForText + "\t" + col
			} else {
				lineForText = lineForText + col
			}

			row = append(row, col)
			rowMap[field.Field] = col

			colVal := stringUtils.ConvertForSql(col)
			if !colIsNumArr[j] { colVal = "'" + colVal + "'" }
			valuesForSql = append(valuesForSql, colVal)
		}

		if format == constant.FormatText {
			logUtils.PrintLine(lineForText)
		} else if format == constant.FormatSql {
			logUtils.PrintLine(genSqlLine(strings.Join(valuesForSql, ", "), i, len(rows)))
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

func printSqlHeader(fields []string, table string) {
	fieldNames := make([]string, 0)
	for _, f := range fields { fieldNames = append(fieldNames, "`" + f + "`") }
	logUtils.PrintLine(fmt.Sprintf("INSERT INTO %s(%s)", table, strings.Join(fieldNames, ", ")))
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

	if i < length - 1 {
		temp = temp + ", "
	} else {
		temp = temp + "; "
	}

	return temp
}

func genJsonLine(i int, row []string,  length int, fields []string) string {
	temp := "  " + RowToJson(row, fields)
	if i < length - 1 {
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
		if j != len(mp) - 1 {
			str = str + "\n"
		}

		j++
	}

	text := fmt.Sprintf("  <row>\n%s\n  </row>", str)
	if i == length - 1 {
		text = text + "\n</testdata>"
	}
	return text
}

func replacePlaceholder(col string) string {
	ret := col

	re := regexp.MustCompile("(?siU)\\${(.*)}")
	arr := re.FindAllStringSubmatch(col, -1)

	strForReplaceMap := map[string][]string{}
	for _, childArr := range arr {
		placeholderStr := childArr[1]
		strForReplaceMap[placeholderStr] = getValForPlaceholder(placeholderStr, len(childArr))

		for _, str := range strForReplaceMap[placeholderStr] {
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

	repeat := "1"
	if repeatObj != nil {
		repeat = repeatObj.(string)
	}

	strs := make([]string, 0)
	if tp == "int" {
		start := mp["start"].(string)
		end := mp["end"].(string)
		precision := mp["precision"].(string)
		format := mp["format"].(string)

		strs = GetRandFromRange("int", start, end, "1", repeat, precision, count, format)
	} else if tp == "float" {
		start := mp["start"].(string)
		end := mp["end"].(string)
		precision := mp["precision"].(string)
		format := mp["format"].(string)

		strs = GetRandFromRange("float", start, end, "1", repeat, precision, count, format)
	} else if tp == "char" {
		start := mp["start"].(string)
		end := mp["end"].(string)
		precision := mp["precision"].(string)
		format := mp["format"].(string)

		strs = GetRandFromRange("char", start, end, "1", repeat, precision, count, format)
	} else if tp == "list" {
		list := mp["list"].([]string)
		strs = GetRandFromList(list, repeat, count)
	}

	return strs
}