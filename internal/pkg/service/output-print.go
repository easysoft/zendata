package service

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zendata/internal/pkg/helper"
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
