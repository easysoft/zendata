package service

import (
	"fmt"
	"github.com/easysoft/zendata/internal/pkg/helper"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"io/ioutil"
	"strings"
	"time"
)

type MsTableParseService struct {
	SqlParseService   *SqlParseService   `inject:""`
	TableParseService *TableParseService `inject:""`
}

func (s *MsTableParseService) GenYamlFromTable() {
	startTime := time.Now().Unix()

	db, err := gorm.Open(sqlserver.Open(vari.GlobalVars.DBDsn))
	if err != nil {
		logUtils.PrintTo(
			fmt.Sprintf("Error on opening db %s, error is %s", vari.GlobalVars.DBDsnParsing.DbName, err.Error()))
	}

	tableName := vari.GlobalVars.Table
	tableNames := strings.Split(tableName, ".")
	sqlForTableStructure := fmt.Sprintf(`
select
    c.name as columnName,
	ty.name as dataType,
	iif(kc.name is not null, 1, 0) as isPrimaryKey,
	c.max_length as columnLength
from sys.columns c
inner join sys.tables t on t.object_id = c.object_id
inner join sys.types ty on c.user_type_id = ty.user_type_id
left join sys.index_columns ic on c.object_id = ic.object_id and ic.column_id = c.column_id
left join sys.key_constraints kc on kc.parent_object_id = ic.object_id and kc.type = 'PK'
where t.name = '%s';`, tableNames[len(tableNames)-1])
	rows, err := db.Raw(sqlForTableStructure).Rows()
	defer rows.Close()

	recordsMap := map[string]map[string][]interface{}{}
	var records []map[string]interface{}
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	db.Raw(query).Scan(&records)
	recordsMap[tableName] = s.TableParseService.GenColArr(records)

	pkMap := map[string]string{}
	fkMap := map[string][2]string{}
	var columnName, columnType, columnLength string
	var isPrimaryKey int
	columns := []string{}
	types := map[string]helper.FieldTypeInfo{}
	for rows.Next() {
		rows.Scan(&columnName, &columnType, &isPrimaryKey, &columnLength)
		if isPrimaryKey == 1 {
			pkMap[tableName] = columnName
		}
		columns = append(columns, columnName)
		types[columnName] = helper.GenerateFieldDefByMetadata(columnType, columnLength, columnName, recordsMap[tableName][columnName])
	}

	s.SqlParseService.genKeysYaml(pkMap)
	s.SqlParseService.writeColumnToYamlFile(tableName, columns, pkMap, fkMap, types)

	entTime := time.Now().Unix()
	files, err := ioutil.ReadDir(vari.GlobalVars.Output)
	logUtils.PrintTo(i118Utils.I118Prt.Sprintf("generate_yaml", len(files), vari.GlobalVars.Output, entTime-startTime))
}
