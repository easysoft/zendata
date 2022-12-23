package genHelper

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/dzwvip/oracle"
	constant "github.com/easysoft/zendata/internal/pkg/const"
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
	db, _ := parserDsnAndConnByGorm(vari.GlobalVars.DBDsn)

	if vari.GlobalVars.DBClear && vari.GlobalVars.Table != "" {
		var deleteSql string
		switch vari.GlobalVars.DBDsnParsing.Driver {
		case "mysql":
			deleteSql = fmt.Sprintf("delete from `%s` where 1=1", vari.GlobalVars.Table)
		case "sqlserver":
			deleteSql = fmt.Sprintf("delete from [%s] where 1=1", vari.GlobalVars.Table)
		case "oracle":
			deleteSql = fmt.Sprintf(`delete from "%s" where 1=1`, vari.GlobalVars.Table)
		}

		err := db.Exec(deleteSql).Error
		if err != nil {
			logUtils.PrintErrMsg(err.Error())
		}
	}

	var insertSql = ""
	switch vari.GlobalVars.DBDsnParsing.Driver {
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

func parserDsnAndConn(dsn string) (conn *sql.DB, err error) {
	var (
		driver, user, password, host, port, db, code string
	)

	// mysql://root:1234@localhost:3306/dbname#utf8
	reg := regexp.MustCompile(`([a-z,A-Z]+)://(.+):(.*)@(.+):(\d+)/(.+)#(.+)`)
	arr := reg.FindAllStringSubmatch(dsn, -1)

	if len(arr) == 0 {
		return
	}

	sections := arr[0]

	driver = strings.ToLower(sections[1])
	user = sections[2]
	password = sections[3]
	host = sections[4]
	port = sections[5]
	db = sections[6]
	code = sections[7]

	if driver == "mysql" {
		str := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", user, password, host, port, db, code)
		conn, _ = sql.Open(driver, str)
		err = conn.Ping() // make sure database is accessible
		if err != nil {
			logUtils.PrintErrMsg(
				fmt.Sprintf("Error on opening db %s, error is %s", constant.SqliteFile, err.Error()))
		}
	}

	return
}

func parserDsnAndConnByGorm(dsn string) (db *gorm.DB, err error) {
	if vari.GlobalVars.DBDsnParsing.Driver == "mysql" {
		dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
			vari.GlobalVars.DBDsnParsing.User,
			vari.GlobalVars.DBDsnParsing.Password,
			vari.GlobalVars.DBDsnParsing.Host,
			vari.GlobalVars.DBDsnParsing.Port,
			vari.GlobalVars.DBDsnParsing.DbName,
			vari.GlobalVars.DBDsnParsing.Code)
		db, err = gorm.Open(mysql.Open(dsn))
		if err != nil { // make sure database is accessible
			logUtils.PrintErrMsg(
				fmt.Sprintf("Error on opening db %s, error is %s", vari.GlobalVars.DBDsnParsing.DbName, err.Error()))
		}
	} else if vari.GlobalVars.DBDsnParsing.Driver == "sqlserver" {
		//str := "sqlserver://sa:12345678Abc@192.168.198.128:1433?database=TestDB"
		// sqlserver 忽略 code （字符编码），该编码暂无发现通过dsn设置，而且由数据库端设置
		dsn := fmt.Sprintf("%s://%s:%s@%s:%s?database=%s",
			vari.GlobalVars.DBDsnParsing.Driver,
			vari.GlobalVars.DBDsnParsing.User,
			vari.GlobalVars.DBDsnParsing.Password,
			vari.GlobalVars.DBDsnParsing.Host,
			vari.GlobalVars.DBDsnParsing.Port,
			vari.GlobalVars.DBDsnParsing.DbName)
		db, err = gorm.Open(sqlserver.Open(dsn))
		if err != nil { // make sure database is accessible
			logUtils.PrintErrMsg(
				fmt.Sprintf("Error on opening db %s, error is %s", vari.GlobalVars.DBDsnParsing.DbName, err.Error()))
		}
	} else if vari.GlobalVars.DBDsnParsing.Driver == "oracle" {
		// orcale 目前使用的gorm orcale适配驱动库不能应用于生产环境，
		// 但作为简单的数据导入需求足以
		dsn := fmt.Sprintf("%s/%s@%s:%s/%s",
			vari.GlobalVars.DBDsnParsing.User,
			vari.GlobalVars.DBDsnParsing.Password,
			vari.GlobalVars.DBDsnParsing.Host,
			vari.GlobalVars.DBDsnParsing.Port,
			vari.GlobalVars.DBDsnParsing.DbName)
		db, err = gorm.Open(oracle.Open(dsn), &gorm.Config{})
		if err != nil { // make sure database is accessible
			logUtils.PrintErrMsg(
				fmt.Sprintf("Error on opening db %s, error is %s", vari.GlobalVars.DBDsnParsing.DbName, err.Error()))
		}
	}

	if vari.Verbose {
		db = db.Debug()
	}

	return
}
