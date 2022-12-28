package service

import (
	"encoding/json"
	"fmt"
	genHelper "github.com/easysoft/zendata/internal/pkg/gen/helper"
	"github.com/easysoft/zendata/internal/pkg/helper"
	"strconv"
	"strings"

	consts "github.com/easysoft/zendata/internal/pkg/const"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
)

type PrintService struct {
	PlaceholderService *PlaceholderService `inject:""`
}

func (s *PrintService) PrintLines() (lines []interface{}) {
	var sqlHeader string

	if vari.GlobalVars.OutputFormat == consts.FormatText {
		s.PrintTextHeader()

	} else if vari.GlobalVars.OutputFormat == consts.FormatSql {
		sqlHeader = s.getInsertSqlHeader()
		if vari.GlobalVars.DBDsn != "" {
			lines = append(lines, sqlHeader)
		}

	} else if vari.GlobalVars.OutputFormat == consts.FormatJson {
		//s.PrintJsonHeader()

	} else if vari.GlobalVars.OutputFormat == consts.FormatXml {
		//s.PrintXmlHeader()
	}

	//for i, cols := range rows {
	//row := make([]string, 0)
	//rowMap := map[string]string{}
	//valuesForSql := make([]string, 0)
	//lineForText := ""
	//
	//for j, col := range cols {
	//	// 3. random replacement
	//	col = s.replacePlaceholder(col)
	//	field := vari.TopFieldMap[fields[j]]
	//	//if field.Length > runewidth.StringWidth(col) {
	//	//col = stringUtils.AddPad(col, field)
	//	//}
	//
	//	if j > 0 && vari.GlobalVars.Human { // use a tab
	//		lineForText = strings.TrimRight(lineForText, "\t")
	//		col = strings.TrimLeft(col, "\t")
	//
	//		lineForText = lineForText + "\t" + col
	//	} else {
	//		lineForText = lineForText + col
	//	}
	//
	//	row = append(row, col)
	//	rowMap[field.Field] = col
	//
	//	colVal := col
	//	if !colIsNumArr[j] {
	//		switch vari.GenVars.DBType {
	//		case consts.DBTypeSqlServer:
	//			colVal = "'" + stringUtils.EscapeValueOfSqlServer(colVal) + "'"
	//		case consts.DBTypeOracle:
	//			colVal = "'" + stringUtils.EscapeValueOfOracle(colVal) + "'"
	//		// case consts.DBTypeMysql:
	//		default:
	//			colVal = "'" + stringUtils.EscapeValueOfMysql(colVal) + "'"
	//		}
	//	}
	//	valuesForSql = append(valuesForSql, colVal)
	//} // for cols
	//
	//if format == consts.FormatText && vari.GlobalVars.DefData.Type == consts.DefTypeArticle { // article need to write to more than one files
	//	lines = append(lines, lineForText)
	//
	//} else if format == consts.FormatText && vari.GlobalVars.DefData.Type != consts.DefTypeArticle {
	//	logUtils.PrintLine(lineForText)
	//
	//} else if format == consts.FormatSql {
	//	if vari.DBDsn != "" { // add to return array for exec sql
	//		sql := strings.Join(valuesForSql, ", ")
	//		lines = append(lines, sql)
	//	} else {
	//
	//		sql := s.genSqlLine(sqlHeader, valuesForSql, vari.GenVars.DBType)
	//		logUtils.PrintLine(sql)
	//	}
	//} else if format == consts.FormatJson {
	//	logUtils.PrintLine(s.genJsonLine(i, row, len(rows), fields))
	//} else if format == consts.FormatXml {
	//	logUtils.PrintLine(s.getXmlLine(i, rowMap, len(rows)))
	//} else if format == consts.FormatData {
	//	lines = append(lines, lineForText)
	//}
	//}

	return
}

func (s *PrintService) PrintTextHeader() {
	if !vari.GlobalVars.Human {
		return
	}
	headerLine := ""
	for idx, field := range vari.GlobalVars.ExportFields {
		headerLine += field
		if idx < len(vari.GlobalVars.ExportFields)-1 {
			headerLine += "\t"
		}
	}

	logUtils.PrintLine(headerLine + "\n")
}

// return "Table> (<column1, column2,...)"
func (s *PrintService) getInsertSqlHeader() string {
	fieldNames := make([]string, 0)
	for _, f := range vari.GlobalVars.ExportFields {
		if vari.GlobalVars.DBType == consts.DBTypeSqlServer {
			f = "[" + helper.EscapeColumnOfSqlServer(f) + "]"
		} else if vari.GlobalVars.DBType == consts.DBTypeOracle {
			f = `"` + f + `"`
		} else {
			f = "`" + helper.EscapeColumnOfMysql(f) + "`"
			//vari.GenVars.DBType == consts.DBTypeMysql {
		}

		fieldNames = append(fieldNames, f)
	}

	var ret string
	switch vari.GlobalVars.DBType {
	case consts.DBTypeSqlServer:
		ret = fmt.Sprintf("[%s] (%s)", vari.GlobalVars.Table, strings.Join(fieldNames, ", "))
	case consts.DBTypeOracle:
		ret = fmt.Sprintf(`"%s" (%s)`, vari.GlobalVars.Table, strings.Join(fieldNames, ", "))
	// case consts.DBTypeMysql:
	default:
		ret = fmt.Sprintf("`%s` (%s)", vari.GlobalVars.Table, strings.Join(fieldNames, ", "))
	}

	return ret
}

func (s *PrintService) rowToJson(cols []string, fieldsToExport []string) string {
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
func (s *PrintService) genSqlLine(sqlheader string, values []string, dbtype string) string {
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

func (s *PrintService) genJsonLine(i int, row []string, length int, fields []string) string {
	temp := s.rowToJson(row, fields)
	if i < length-1 {
		temp = temp + ", "
	} else {
		temp = temp + "\n]"
	}

	return temp
}

func (s *PrintService) getXmlLine(i int, mp map[string]string, length int) string {
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

func (s *PrintService) getValForPlaceholder(placeholderStr string, count int) []string {
	placeholderInt, _ := strconv.Atoi(placeholderStr)
	mp := vari.GlobalVars.RandFieldSectionPathToValuesMap[placeholderInt]

	tp := mp["type"].(string)
	repeatObj := mp["repeat"]

	repeat := 1
	if repeatObj != nil {
		repeat = repeatObj.(int)
	}

	strArr := make([]string, 0)
	repeatTag := mp["repeatTag"].(string)
	if tp == "int" {
		start := mp["start"].(string)
		end := mp["end"].(string)
		precision := mp["precision"].(string)
		format := mp["format"].(string)

		strArr = genHelper.GetRandFromRange("int", start, end, "1",
			repeat, repeatTag, precision, format, count)

	} else if tp == "float" {
		start := mp["start"].(string)
		end := mp["end"].(string)
		stepStr := fmt.Sprintf("%v", mp["step"])

		precision := mp["precision"].(string)
		format := mp["format"].(string)

		strArr = genHelper.GetRandFromRange("float", start, end, stepStr,
			repeat, repeatTag, precision, format, count)

	} else if tp == "char" {
		start := mp["start"].(string)
		end := mp["end"].(string)
		precision := mp["precision"].(string)
		format := mp["format"].(string)

		strArr = genHelper.GetRandFromRange("char", start, end, "1",
			repeat, repeatTag, precision, format, count)

	} else if tp == "list" {
		list := mp["list"].([]string)
		strArr = genHelper.GetRandFromList(list, repeat, count)

	}

	strArr = strArr[:count]
	return strArr
}
