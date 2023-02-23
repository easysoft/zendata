package service

import (
	"fmt"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
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

	s.SqlParseService.genKeysYaml(pkMap)

	s.SqlParseService.genTablesYaml(statementMap, pkMap, fkMap)

	entTime := time.Now().Unix()
	logUtils.PrintTo(i118Utils.I118Prt.Sprintf("generate_yaml", len(statementMap), vari.GlobalVars.Output, entTime-startTime))
}
