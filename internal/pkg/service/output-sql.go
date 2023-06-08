package service

import (
	"fmt"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/helper"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"strings"
)

func (s *OutputService) GenSql() {
	records := s.GenRecords()

	lines := make([]interface{}, 0)

	sqlHeader := s.getInsertSqlHeader()
	if vari.GlobalVars.DBDsn != "" {
		lines = append(lines, sqlHeader)
	}

	logUtils.PrintLine(sqlHeader)

	for index, record := range records {
		valuesForSql := make([]string, 0)

		for j, colName := range vari.GlobalVars.ExportFields {
			colVal := fmt.Sprintf("%v", record[colName])

			if !vari.GlobalVars.ColIsNumArr[j] {
				switch vari.GlobalVars.DBType {
				case consts.DBTypeMysql:
					colVal = "'" + helper.EscapeValueOfMysql(colVal) + "'"
				case consts.DBTypeOracle:
					colVal = "'" + helper.EscapeValueOfOracle(colVal) + "'"
				case consts.DBTypeSqlServer:
					colVal = "'" + helper.EscapeValueOfSqlServer(colVal) + "'"
				default:
				}
			}

			valuesForSql = append(valuesForSql, colVal)
		}

		if vari.GlobalVars.DBDsn != "" { // add to return array for sql exec
			sql := strings.Join(valuesForSql, ", ")
			lines = append(lines, sql)
		} else {
			sql := s.genSqlLine(valuesForSql)
			sql = strings.Repeat(" ", len("INSERT")+1) + sql

			if index < len(records)-1 {
				sql += ","
				logUtils.PrintLine(sql + "\n")

			} else {
				if vari.GlobalVars.DBType == consts.DBTypeSqlServer {
					logUtils.PrintLine(sql + "; GO")
				} else {
					logUtils.PrintLine(sql + ";")
				}
			}
		}
	}

	logUtils.PrintLine("\n")

	return
}

// return Table (column1, column2, ...)
func (s *OutputService) getInsertSqlHeader() string {
	fieldNames := make([]string, 0)

	for _, f := range vari.GlobalVars.ExportFields {
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
		ret = fmt.Sprintf("`%s` (%s)", vari.GlobalVars.Table, strings.Join(fieldNames, ", "))
	case consts.DBTypeOracle:
		ret = fmt.Sprintf(`"%s" (%s)`, vari.GlobalVars.Table, strings.Join(fieldNames, ", "))
	case consts.DBTypeSqlServer:
		ret = fmt.Sprintf("[%s] (%s)", vari.GlobalVars.Table, strings.Join(fieldNames, ", "))
	default:
	}

	ret = "INSERT INTO " + ret + " VALUES\n"

	return ret
}

func (s *OutputService) genSqlLine(values []string) string {
	return "(" + strings.Join(values, ",") + ")"
}
