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

	sql := "SELECT city.id, city.name, parent.name state, city.zipCode, city.cityCode " +
			"FROM cn_city city JOIN cn_city parent ON city.parentCode = parent.areaCode"
	rows, err := db.Query(sql)
	if err != nil {
		logUtils.Screen("fail to exec query " + sql)
		return
	}

	for rows.Next() {
		var id int
		var name string
		var state string
		var zipCode int
		var cityCode int

		err = rows.Scan(&id, &name, &state, &zipCode, &cityCode)
		if err != nil {
			logUtils.Screen("fail to get sqlite3 row")
			return
		}
		fmt.Println(id, name, state, zipCode, cityCode)
	}
}