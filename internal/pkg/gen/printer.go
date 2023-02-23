package gen

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zendata/internal/pkg/helper"
	"regexp"
	"strconv"
	"strings"

	consts "github.com/easysoft/zendata/internal/pkg/const"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
)

func PrintLines(rows [][]string, format string, table string, colIsNumArr []bool,
	fields []string) (lines []interface{}) {

	var sqlHeader string

	if format == consts.FormatText {
		PrintHumanHeaderIfNeeded(fields)

	} else if format == consts.FormatSql {
		sqlHeader = getInsertSqlHeader(fields, table)
		if vari.GlobalVars.DBDsn != "" {
			lines = append(lines, sqlHeader)
		}

	} else if format == consts.FormatJson {
		printJsonHeader()

	} else if format == consts.FormatXml {
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
			field := vari.GlobalVars.TopFieldMap[fields[j]]
			//if field.Length > runewidth.StringWidth(col) {
			//col = stringUtils.AddPad(col, field)
			//}

			if j > 0 && vari.GlobalVars.Human { // use a tab
				lineForText = strings.TrimRight(lineForText, "\t")
				col = strings.TrimLeft(col, "\t")

				lineForText = lineForText + "\t" + col
			} else {
				lineForText = lineForText + col
			}

			row = append(row, col)
			rowMap[field.Field] = col

			colVal := col
			if !colIsNumArr[j] {
				switch vari.GlobalVars.DBType {
				case consts.DBTypeSqlServer:
					colVal = "'" + helper.EscapeValueOfSqlServer(colVal) + "'"
				case consts.DBTypeOracle:
					colVal = "'" + helper.EscapeValueOfOracle(colVal) + "'"
				// case consts.DBTypeMysql:
				default:
					colVal = "'" + helper.EscapeValueOfMysql(colVal) + "'"
				}
			}
			valuesForSql = append(valuesForSql, colVal)
		} // for cols

		if format == consts.FormatText && vari.GlobalVars.DefData.Type == consts.DefTypeArticle { // article need to write to more than one files
			lines = append(lines, lineForText)

		} else if format == consts.FormatText && vari.GlobalVars.DefData.Type != consts.DefTypeArticle {
			logUtils.PrintLine(lineForText)

		} else if format == consts.FormatSql {
			if vari.GlobalVars.DBDsn != "" { // add to return array for exec sql
				sql := strings.Join(valuesForSql, ", ")
				lines = append(lines, sql)
			} else {

				sql := genSqlLine(sqlHeader, valuesForSql, vari.GlobalVars.DBType)
				logUtils.PrintLine(sql)
			}
		} else if format == consts.FormatJson {
			logUtils.PrintLine(genJsonLine(i, row, len(rows), fields))
		} else if format == consts.FormatXml {
			logUtils.PrintLine(getXmlLine(i, rowMap, len(rows)))
		} else if format == consts.FormatData {
			lines = append(lines, lineForText)
		}
	}

	return
}

func PrintHumanHeaderIfNeeded(fields []string) {
	if !vari.GlobalVars.Human {
		return
	}
	headerLine := ""
	for idx, field := range fields {
		headerLine += field
		if idx < len(fields)-1 {
			headerLine += "\t"
		}
	}

	logUtils.PrintLine(headerLine + "\n")
}

// return Table (column1, column2, ...)
func getInsertSqlHeader(fields []string, table string) string {
	fieldNames := make([]string, 0)

	for _, f := range fields {
		if vari.GlobalVars.DBType == consts.DBTypeMysql {
			f = "`" + helper.EscapeColumnOfMysql(f) + "`"
		} else if vari.GlobalVars.DBType == consts.DBTypeOracle {
			f = `"` + f + `"`
		} else if vari.GlobalVars.DBType == consts.DBTypeSqlServer {
			f = "[" + helper.EscapeColumnOfSqlServer(f) + "]"
		}

		fieldNames = append(fieldNames, f)
	}

	var ret string
	switch vari.GlobalVars.DBType {
	case consts.DBTypeMysql:
		ret = fmt.Sprintf("`%s` (%s)", table, strings.Join(fieldNames, ", "))
	case consts.DBTypeOracle:
		ret = fmt.Sprintf(`"%s" (%s)`, table, strings.Join(fieldNames, ", "))
	case consts.DBTypeSqlServer:
		ret = fmt.Sprintf("[%s] (%s)", table, strings.Join(fieldNames, ", "))
	default:
	}

	return ret
}

func printJsonHeader() {
	logUtils.PrintLine("[")
}

func printXmlHeader(fields []string, table string) {
	logUtils.PrintLine("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<testdata>\n  <title>Test Data</title>")
}

func rowToJson(cols []string, fieldsToExport []string) string {
	rowMap := map[string]string{}
	for j, col := range cols {
		rowMap[fieldsToExport[j]] = col
	}

	jsonObj, _ := json.MarshalIndent(rowMap, "", "\t")
	respJson := string(jsonObj)

	respJson = strings.ReplaceAll("\t"+respJson, "\n", "\n\t")

	return respJson
}

// @return ""
func genSqlLine(sqlheader string, values []string, dbtype string) string {
	var tmp string
	switch dbtype {
	case consts.DBTypeSqlServer:
		tmp = "INSERT INTO " + sqlheader + " VALUES (" + strings.Join(values, ",") + "); GO"
	default:
		// consts.DBTypeMysql
		// consts.DBTypeOracle:
		tmp = "INSERT INTO " + sqlheader + " VALUES (" + strings.Join(values, ",") + ");"
	}

	return tmp
}

func genJsonLine(i int, row []string, length int, fields []string) string {
	temp := rowToJson(row, fields)
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
			key, _ := strconv.Atoi(placeholderStr)
			temp := Placeholder(key)
			ret = strings.Replace(ret, temp, str, 1)
		}
	}

	return ret
}

func getValForPlaceholder(placeholderStr string, count int) (strArr []string) {
	//placeholderInt, _ := strconv.Atoi(placeholderStr)
	//mp := vari.GlobalVars.RandFieldSectionPathToValuesMap[placeholderInt]
	//
	//tp := mp["type"].(string)
	//repeatObj := mp["repeat"]
	//
	//repeat := 1
	//if repeatObj != nil {
	//	repeat = repeatObj.(int)
	//}
	//
	//repeatTag := mp["repeatTag"].(string)
	//if tp == "int" {
	//	start := mp["start"].(string)
	//	end := mp["end"].(string)
	//	precision := mp["precision"].(string)
	//	format := mp["format"].(string)
	//
	//	strArr = genHelper.GetRandValuesFromRange("int", start, end, "1",
	//		repeat, repeatTag, precision, format, count)
	//
	//} else if tp == "float" {
	//	start := mp["start"].(string)
	//	end := mp["end"].(string)
	//	stepStr := fmt.Sprintf("%v", mp["step"])
	//
	//	precision := mp["precision"].(string)
	//	format := mp["format"].(string)
	//
	//	strArr = genHelper.GetRandValuesFromRange("float", start, end, stepStr,
	//		repeat, repeatTag, precision, format, count)
	//
	//} else if tp == "char" {
	//	start := mp["start"].(string)
	//	end := mp["end"].(string)
	//	precision := mp["precision"].(string)
	//	format := mp["format"].(string)
	//
	//	strArr = genHelper.GetRandValuesFromRange("char", start, end, "1",
	//		repeat, repeatTag, precision, format, count)
	//
	//} else if tp == "list" {
	//	list := mp["list"].([]string)
	//	strArr = genHelper.GetRandFromList(list, repeat, count)
	//
	//}
	//
	//strArr = strArr[:count]
	//return strArr

	// method not used
	return
}
