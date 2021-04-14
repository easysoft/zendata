package helper

import (
	"database/sql"
	"fmt"
	constant "github.com/easysoft/zendata/src/utils/const"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"github.com/easysoft/zendata/src/utils/vari"
	_ "github.com/go-sql-driver/mysql"
	"regexp"
	"strings"
)

var (
	IgnoreCategories = []string{"姓", "名字"}
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

func LoadAllWords() (ret map[string]string) {
	ret = map[string]string{}

	sqlStr := fmt.Sprintf("SELECT * FROM words_v1")
	rows, err := vari.DB.Query(sqlStr)
	defer rows.Close()
	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_exec_query", sqlStr, err.Error()))
		return
	}

	columns, err := rows.Columns()
	colNum := len(columns)

	colIndexToName := map[int]string{}
	for index, col := range columns {
		colIndexToName[index] = col
	}

	var values = make([]interface{}, colNum)
	for i, _ := range values {
		var itf string
		values[i] = &itf
	}

	for rows.Next() {
		err = rows.Scan(values...)
		if err != nil {
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_parse_row", err.Error()))
			return
		}

		//for index, v := range values {
		//	item := v.(*string)
		//	if *item == "y" {
		//		key := values[1].(*string)
		//		ret[*key] = colIndexToName[index]
		//		break
		//	}
		//}

		for index := len(values) - 1; index >= 0; index-- {
			val := values[index]
			item := val.(*string)

			if *item == "y" {
				key := values[1].(*string)
				if !stringUtils.StrInArr(colIndexToName[index], IgnoreCategories) {
					ret[*key] = colIndexToName[index]
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
