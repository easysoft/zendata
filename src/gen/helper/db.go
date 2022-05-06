package helper

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	constant "github.com/easysoft/zendata/src/utils/const"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"github.com/easysoft/zendata/src/utils/vari"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var (
	IgnoreWords      = []string{"了", "的"}
	IgnoreCategories = []string{"姓", "名字", "介词"}
)

func ExecSqlInUserDB(lines []interface{}) (count int) {
	sql := ""

	for _, line := range lines {
		sql += line.(string) + " "
	}

	//typ, user, password, host, port, db, code := parserDsn(vari.DBDsn)
	db, _ := parserDsnAndConnByGorm(vari.DBDsn)

	if vari.DBClear && vari.Table != "" {
		deleteSql := fmt.Sprintf("delete from %s where 1=1", vari.Table)
		err := db.Exec(deleteSql).Error
		if err != nil {
			logUtils.PrintErrMsg(err.Error())
		}
	}
	err := db.Exec(sql).Error
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
		conn, err = sql.Open(driver, str)
		err = conn.Ping() // make sure database is accessible
		if err != nil {
			logUtils.PrintErrMsg(
				fmt.Sprintf("Error on opening db %s, error is %s", constant.SqliteFile, err.Error()))
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
	}

	if vari.Verbose {
		db = db.Debug()
	}

	return
}
