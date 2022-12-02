package helper

import (
	"fmt"
	"github.com/dzwvip/oracle"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	stringUtils "github.com/easysoft/zendata/pkg/utils/string"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var (
	IgnoreWords      = []string{"了", "的"}
	IgnoreCategories = []string{"姓", "名字", "介词"}
)

func ExecSqlInUserDB(lines []interface{}) (count int) {

	//typ, user, password, host, port, db, code := parserDsn(vari.DBDsn)
	db, _ := parserDsnAndConnByGorm(vari.DBDsn)

	if vari.DBClear && vari.Table != "" {
		var deleteSql string
		switch vari.DBDsnParsing.Driver {
		case "mysql":
			deleteSql = fmt.Sprintf("delete from `%s` where 1=1", vari.Table)
		case "sqlserver":
			deleteSql = fmt.Sprintf("delete from [%s] where 1=1", vari.Table)
		case "oracle":
			deleteSql = fmt.Sprintf(`delete from "%s" where 1=1`, vari.Table)
		}

		err := db.Exec(deleteSql).Error
		if err != nil {
			logUtils.PrintErrMsg(err.Error())
		}
	}

	var insertSql = ""
	switch vari.DBDsnParsing.Driver {
	case "mysql", "sqlserver":
		insertSql = "INSERT INTO " + lines[0].(string) + " VALUES (" + lines[1].(string) + ")"

		for _, line := range lines[2:] {
			insertSql += ", (" + line.(string) + ")"
		}
		insertSql += ";"

	case "oracle":
		insertSql = "INSERT ALL "
		headerStr := " INTO " + lines[0].(string) + " VALUES ( "

		for _, line := range lines[1:] {
			insertSql += headerStr + line.(string) + ") "
		}

		insertSql += " SELECT 1 FROM DUAL"
	}

	err := db.Exec(insertSql).Error
	if err != nil {
		logUtils.PrintErrMsg(err.Error())
	}

	return
}

func LoadAllWords() (ret map[string]string) {
	ret = map[string]string{}

	rows, _ := vari.DB.Table("words_v1").Where("true").Select("*").Rows()
	defer rows.Close()

	columns, err := rows.Columns()
	colNum := len(columns)

	colIndexToCategoryName := map[int]string{}
	for index, col := range columns {
		colIndexToCategoryName[index] = col
	}

	// build an empty string array to retrieve row
	var record = make([]interface{}, colNum)
	for i, _ := range record {
		var itf string
		record[i] = &itf
	}

	for rows.Next() {
		err = rows.Scan(record...)
		if err != nil {
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_parse_row", err.Error()))
			return
		}

		for index := len(record) - 1; index >= 0; index-- {
			word := record[1].(*string)
			category := colIndexToCategoryName[index]
			isBelowToCategory := record[index].(*string)

			if *isBelowToCategory == "y" {
				if !stringUtils.StrInArr(category, IgnoreCategories) &&
					!stringUtils.StrInArr(*word, IgnoreWords) {

					ret[*word] = category
				}

				break
			}
		}
	}

	return
}

func parserDsnAndConnByGorm(dsn string) (db *gorm.DB, err error) {
	if vari.DBDsnParsing.Driver == "mysql" {
		dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
			vari.DBDsnParsing.User,
			vari.DBDsnParsing.Password,
			vari.DBDsnParsing.Host,
			vari.DBDsnParsing.Port,
			vari.DBDsnParsing.DbName,
			vari.DBDsnParsing.Code)
		db, err = gorm.Open(mysql.Open(dsn))
		if err != nil { // make sure database is accessible
			logUtils.PrintErrMsg(
				fmt.Sprintf("Error on opening db %s, error is %s", vari.DBDsnParsing.DbName, err.Error()))
		}
	} else if vari.DBDsnParsing.Driver == "sqlserver" {
		//str := "sqlserver://sa:12345678Abc@192.168.198.128:1433?database=TestDB"
		// sqlserver 忽略 code （字符编码），该编码暂无发现通过dsn设置，而且由数据库端设置
		dsn := fmt.Sprintf("%s://%s:%s@%s:%s?database=%s",
			vari.DBDsnParsing.Driver,
			vari.DBDsnParsing.User,
			vari.DBDsnParsing.Password,
			vari.DBDsnParsing.Host,
			vari.DBDsnParsing.Port,
			vari.DBDsnParsing.DbName)
		db, err = gorm.Open(sqlserver.Open(dsn))
		if err != nil { // make sure database is accessible
			logUtils.PrintErrMsg(
				fmt.Sprintf("Error on opening db %s, error is %s", vari.DBDsnParsing.DbName, err.Error()))
		}
	} else if vari.DBDsnParsing.Driver == "oracle" {
		// orcale 目前使用的gorm orcale适配驱动库不能应用于生产环境，
		// 但作为简单的数据导入需求足以
		dsn := fmt.Sprintf("%s/%s@%s:%s/%s",
			vari.DBDsnParsing.User,
			vari.DBDsnParsing.Password,
			vari.DBDsnParsing.Host,
			vari.DBDsnParsing.Port,
			vari.DBDsnParsing.DbName)
		db, err = gorm.Open(oracle.Open(dsn), &gorm.Config{})
		if err != nil { // make sure database is accessible
			logUtils.PrintErrMsg(
				fmt.Sprintf("Error on opening db %s, error is %s", vari.DBDsnParsing.DbName, err.Error()))
		}
	}

	if vari.Verbose {
		db = db.Debug()
	}

	return
}
