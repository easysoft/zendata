package gen

import (
	"fmt"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/kataras/iris/v12"
	"strconv"
	"strings"
)

func ParseCache() (cacheKey, cacheOpt string, batch int, hasCache, isBatch bool) {
	if vari.CacheParam == "" {
		return
	}

	arr := strings.Split(vari.CacheParam, "-")

	if len(arr) > 0 {
		cacheKey = arr[0]
	}

	if len(arr) > 1 {
		var err error
		batch, err = strconv.Atoi(arr[1])

		if err != nil {
			cacheOpt = arr[1]
		}

	}

	hasCache, isBatch = HasCache(cacheKey)

	return
}

func RetrieveCache(cacheKey string, fieldsToExport *[]string) (rows [][]string, colIsNumArr []bool, err error) {
	record := iris.Map{}
	sql2 := "SELECT * FROM `" + getTableNameIsNum(cacheKey) + "`;"
	err = vari.DB.Raw(sql2).
		Scan(&record).
		Error

	colIsNumArr = stringToBoolArr(record["is_nums"].(string))

	selectedCols := ""
	if len(*fieldsToExport) > 0 {
		selectedCols = strings.Join(*fieldsToExport, ",")
	} else {
		selectedCols = record["fields"].(string)
		*fieldsToExport = strings.Split(selectedCols, ",")
	}

	var records []iris.Map
	sql := "SELECT " + selectedCols +
		" FROM " + getTableName(cacheKey) +
		" LIMIT " + strconv.Itoa(vari.GlobalVars.Total) + " ;"
	err = vari.DB.Raw(sql).
		Scan(&records).
		Error

	selectedColArr := strings.Split(selectedCols, ",")
	for _, record := range records {
		row := make([]string, 0)
		for _, col := range selectedColArr {
			row = append(row, record[col].(string))
		}

		rows = append(rows, row)
	}

	return
}

func RetrieveCacheBatch(cacheKey string, fieldsToExport *[]string, batch int) (
	rows [][]string, colIsNumArr []bool, err error) {

	cacheKey = getBatchKey(cacheKey, batch)

	record := iris.Map{}
	baseKey := removeBatchNumInKey(cacheKey)
	sql2 := "SELECT * FROM `" + getTableNameIsNum(baseKey) + "`;"
	err = vari.DB.Raw(sql2).
		Scan(&record).
		Error

	colIsNumArr = stringToBoolArr(record["is_nums"].(string))

	selectedCols := ""
	if len(*fieldsToExport) > 0 {
		selectedCols = strings.Join(*fieldsToExport, ",")
	} else {
		selectedCols = record["fields"].(string)
		*fieldsToExport = strings.Split(selectedCols, ",")
	}

	var records []iris.Map
	sql := "SELECT " + selectedCols +
		" FROM " + getTableName(cacheKey) +
		" LIMIT " + strconv.Itoa(vari.GlobalVars.Total) + " ;"
	err = vari.DB.Raw(sql).
		Scan(&records).
		Error

	selectedColArr := strings.Split(selectedCols, ",")
	for _, record := range records {
		row := make([]string, 0)
		for _, col := range selectedColArr {
			row = append(row, record[col].(string))
		}

		rows = append(rows, row)
	}

	return
}

func CreateCache(cacheKey string, fieldsToExport []string, rows [][]string, colIsNumArr []bool) {
	ClearCache(cacheKey)

	CreateCacheIsNumTable(cacheKey, fieldsToExport)
	CreateCacheDataTable(cacheKey, fieldsToExport)

	CreateCacheData(cacheKey, fieldsToExport, rows)
	CreateCacheIsNum(cacheKey, fieldsToExport, colIsNumArr)

	return
}

func CreateCacheData(cacheKey string, fieldsToExport []string, rows [][]string) (err error) {
	insertTemplate := "INSERT INTO `" + getTableName(cacheKey) + "` (%s) VALUES %s;"

	cols := make([]string, 0)
	vals := make([]string, 0)

	for _, col := range fieldsToExport {
		cols = append(cols, "`"+col+"`")
	}

	for _, row := range rows {
		arr := make([]string, 0)
		for _, c := range row {
			arr = append(arr, fmt.Sprintf("'%s'", c))
		}

		vals = append(vals, fmt.Sprintf("(%s)", strings.Join(arr, ",")))
	}

	insertSql := fmt.Sprintf(insertTemplate, strings.Join(cols, ","), strings.Join(vals, ","))
	err = vari.DB.Exec(insertSql).Error
	if err != nil {
		logUtils.PrintLine(i118Utils.I118Prt.Sprintf("fail_to_exec_query", insertSql, err.Error()))
	}

	return
}

func CreateCacheIsNum(cacheKey string, fieldsToExport []string, colIsNumArr []bool) (err error) {
	insertTemplate2 := "INSERT INTO `" + getTableNameIsNum(cacheKey) +
		"` (is_nums, fields) VALUES ('%s', '%s');"
	insertSql2 := fmt.Sprintf(insertTemplate2, boolArrToString(colIsNumArr), strings.Join(fieldsToExport, ","))
	err = vari.DB.Exec(insertSql2).Error
	if err != nil {
		logUtils.PrintLine(i118Utils.I118Prt.Sprintf("fail_to_exec_query", insertSql2, err.Error()))
	}

	return
}

func CreateCacheDataTable(cacheKey string, fieldsToExport []string) (err error) {
	var ddlFields []string

	for _, colName := range fieldsToExport {
		colDefine := fmt.Sprintf("    `%s` VARCHAR DEFAULT ''", colName)
		ddlFields = append(ddlFields, colDefine)
	}

	ddlTemplate :=
		`CREATE TABLE %s (
	id INTEGER PRIMARY KEY AUTOINCREMENT, 
%s
);`

	ddlSql := fmt.Sprintf(ddlTemplate, getTableName(cacheKey), strings.Join(ddlFields, ", \n"))
	err = vari.DB.Exec(ddlSql).Error
	if err != nil {
		logUtils.PrintLine(i118Utils.I118Prt.Sprintf("fail_to_create_table", cacheKey, err.Error()))
		return
	}

	return
}

func CreateCacheIsNumTable(cacheKey string, fieldsToExport []string) (err error) {
	ddlTemplate2 :=
		`CREATE TABLE %s (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			is_nums VARCHAR DEFAULT '',
			fields VARCHAR DEFAULT ''
		);`
	ddlSql2 := fmt.Sprintf(ddlTemplate2, getTableNameIsNum(cacheKey))
	err = vari.DB.Exec(ddlSql2).Error

	c.Inc("is_nums_table_created")

	return
}

func ClearCache(cacheKey string) (err error) {
	dropSql := fmt.Sprintf("DROP TABLE IF EXISTS %s;", getTableName(cacheKey))
	err = vari.DB.Exec(dropSql).Error
	if err != nil {
		logUtils.PrintLine(i118Utils.I118Prt.Sprintf("fail_to_drop_table", cacheKey, err.Error()))
	}

	dropSql2 := fmt.Sprintf("DROP TABLE IF EXISTS %s;", getTableNameIsNum(cacheKey))
	err = vari.DB.Exec(dropSql2).Error

	return
}

func ClearAllCache() (err error) {
	tables := make([]string, 0)
	vari.DB.Raw("SELECT name fROM sqlite_master").
		Scan(&tables)

	for _, table := range tables {
		if strings.Index(table, consts.CachePrefix) < 0 {
			continue
		}
		dropSql := fmt.Sprintf("DROP TABLE IF EXISTS %s;", table)
		err = vari.DB.Exec(dropSql).Error
	}

	return
}

func ClearBatchCache(key string) (err error) {
	tables := make([]string, 0)
	vari.DB.Raw("SELECT name fROM sqlite_master").
		Scan(&tables)

	for _, table := range tables {
		batchTablePrefix := fmt.Sprintf("%s%s_", consts.CachePrefix, key)
		if strings.Index(table, batchTablePrefix) < 0 {
			continue
		}
		dropSql := fmt.Sprintf("DROP TABLE IF EXISTS %s;", table)
		err = vari.DB.Exec(dropSql).Error
	}

	return
}

func HasCache(key string) (hasCache, isBatch bool) {
	records := make(iris.Map, 0)

	err := vari.DB.Raw(fmt.Sprintf("select id from %s LIMIT 1", getTableName(key))).
		Scan(&records).
		Error
	if err == nil && len(records) > 0 {
		hasCache = true
		return
	}

	err = vari.DB.Raw(fmt.Sprintf("select id from %s LIMIT 1", getTableName(key)+"_0")).
		Scan(&records).
		Error
	if err == nil && len(records) > 0 {
		hasCache = true
		isBatch = true
	}

	return
}

func boolArrToString(colIsNumArr []bool) (ret string) {
	arr := make([]string, 0)

	for _, b := range colIsNumArr {
		val := "0"
		if b {
			val = "1"
		}

		arr = append(arr, val)
	}

	ret = strings.Join(arr, ",")

	return
}

func stringToBoolArr(str string) (ret []bool) {
	arr := strings.Split(str, ",")

	for _, s := range arr {
		val := false
		if s == "1" {
			val = true
		}

		ret = append(ret, val)
	}

	return
}

func getTableName(key string) (ret string) {
	return consts.CachePrefix + key
}

func getTableNameIsNum(key string) (ret string) {
	return consts.CachePrefix + key + consts.CachePostfix
}
