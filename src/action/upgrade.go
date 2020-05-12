package action

import (
	"database/sql"
	"fmt"
	constant "github.com/easysoft/zendata/src/utils/const"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	numbUtils "github.com/easysoft/zendata/src/utils/numb"

	_ "github.com/mattn/go-sqlite3"
)

func Upgrade() {
	// TODO: update date from remote server

	// read data from db
	db, err := sql.Open(constant.SqliteDriver, constant.SqliteSource)
	if err != nil {
		logUtils.Screen("fail to open " + constant.SqliteSource + ": " + err.Error())
		return
	}

	sql := "SELECT id, 'name', state, zipCode, cityCode FROM city"
	rows, err := db.Query(sql)
	if err != nil {
		logUtils.Screen("fail to exec query " + sql + ": " + err.Error())
		return
	}

	sqls := make([]string, 0)
	for rows.Next() {
		var id int
		var name string
		var state string
		var zipCode string
		var cityCode string

		err = rows.Scan(&id, &name, &state, &zipCode, &cityCode)
		if err != nil {
			logUtils.Screen("fail to get sqlite3 row: " + err.Error())
			return
		}
		fmt.Println(numbUtils.NumToBHex(id), name, state, zipCode, cityCode)

		sql := fmt.Sprintf("UPDATE city SET seq = '%s' where id = %d;", numbUtils.NumToBHex(id), id)
		sqls = append(sqls, sql)
	}

	// update seq column
	for _, sql := range sqls {
		_, err := db.Exec(sql)
		if err != nil {
			logUtils.Screen("fail to update row: " + err.Error())
		}
	}
}