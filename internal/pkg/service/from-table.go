package service

import (
	"fmt"
	"time"

	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TableParseService struct {
	SqlParseService *SqlParseService `inject:""`
}

func (s *TableParseService) GenYamlFromTable() {
	startTime := time.Now().Unix()

	db, err := gorm.Open(mysql.Open(vari.GlobalVars.DBDsn))
	if err != nil {
		logUtils.PrintTo(
			fmt.Sprintf("Error on opening db %s, error is %s", vari.GlobalVars.DBDsnParsing.DbName, err.Error()))
	}

	var mp map[string]interface{}
	db.Raw("SHOW CREATE TABLE " + vari.GlobalVars.Table).Scan(&mp)
	sql := mp["Create Table"].(string) + ";"
	statementMap, pkMap, fkMap := s.SqlParseService.getCreateStatement(sql)

	recordsMap := map[string]map[string][]interface{}{}
	var records []map[string]interface{}
	query := fmt.Sprintf("SELECT * FROM %s", vari.GlobalVars.Table)
	db.Raw(query).Scan(&records)
	recordsMap[vari.GlobalVars.Table] = s.GenColArr(records)

	s.SqlParseService.genKeysYaml(pkMap)
	s.SqlParseService.genTablesYaml(statementMap, pkMap, fkMap, recordsMap)

	entTime := time.Now().Unix()
	logUtils.PrintTo(i118Utils.I118Prt.Sprintf("generate_yaml", len(statementMap), vari.GlobalVars.Output, entTime-startTime))
}

func (s *TableParseService) GenColArr(records []map[string]interface{}) (ret map[string][]interface{}) {
	ret = map[string][]interface{}{}

	for _, record := range records {
		for key, val := range record {
			if val != nil {
				ret[key] = append(ret[key], val)
			}
		}
	}

	return
}
