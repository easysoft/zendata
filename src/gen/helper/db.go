package helper

import (
	"database/sql"
	"fmt"
	constant "github.com/easysoft/zendata/src/utils/const"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"github.com/easysoft/zendata/src/utils/vari"
	_ "github.com/go-sql-driver/mysql"
	"regexp"
	"strings"
)

func ExecSql(lines []interface{}) (count int) {
	sql := ""

	for _, line := range lines {
		sql += line.(string) + " "
	}

	//typ, user, password, host, port, db, code := parserDsn(vari.DBDsn)
	conn, _ := parserDsnAndConn(vari.DBDsn)
	defer conn.Close()

	if vari.DBClear && vari.Table != "" {
		deleteSql := fmt.Sprintf("delete from %s where 1=1", vari.Table)
		conn.Exec(deleteSql)
	}
	conn.Exec(sql)

	return
}

func parserDsnAndConn(dsn string) (conn *sql.DB, err error) {
	var (
		driver, user, password, host, port, db, code string
	)

	// mysql://root:1234@localhost:3306/dbname#utf8
	reg := regexp.MustCompile(`([a-z,A-Z]+)://(.+):(.+)@(.+):(\d+)/(.+)#(.+)`)
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
				fmt.Sprintf("Error on opening db %s, error is %s", constant.SqliteData, err.Error()))
		}
	}

	return
}
