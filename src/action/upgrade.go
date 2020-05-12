package action

import (
	"database/sql"
	"fmt"
	constant "github.com/easysoft/zendata/src/utils/const"
	logUtils "github.com/easysoft/zendata/src/utils/log"

	_ "github.com/mattn/go-sqlite3"
)

func Upgrade() {
	db, err := sql.Open(constant.SqliteDriver, constant.SqliteSource)
	if err != nil {
		logUtils.Screen("fail to open " + constant.SqliteSource)
		return
	}

	sql := "SELECT * FROM cn_city"
	rows, err := db.Query(sql)
	if err != nil {
		logUtils.Screen("fail to exec query " + sql)
		return
	}

	for rows.Next() {
		var id int
		var level int
		var parentCode int
		var areaCode int
		var zipCode int
		var cityCode int
		var name string

		err = rows.Scan(&id, &level, &parentCode, &areaCode, zipCode, cityCode, name)
		if err != nil {
			logUtils.Screen("fail to get sqlite3 row")
			return
		}
		fmt.Println(id, name, level, parentCode, areaCode, zipCode, cityCode)
	}
}