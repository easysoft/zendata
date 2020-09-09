package main

import (
	"database/sql"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/Chain-Zhang/pinyin"
	commonUtils "github.com/easysoft/zendata/src/utils/common"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestImportSqlite(t *testing.T) {
	files := make([]string, 0)
	getExcelFilesInDir("xdoc/words-9.3", &files)

	for _, file := range files {
		importExcel(file, t)
	}
}

func importExcel(filePath string, t *testing.T) {
	excel, err := excelize.OpenFile(filePath)
	if err != nil {
		t.Error("fail to read file " + filePath + ", error: " + err.Error())
		return
	}

	fileName := path.Base(filePath)
	fileName = strings.TrimSuffix(fileName, path.Ext(filePath))
	fileName = strings.TrimSuffix(fileName, "词库")
	tableName := getPinyin(fileName)

	for _, sheet := range excel.GetSheetList() {
		rows, err := excel.GetRows(sheet)

		if len(rows) == 0 {
			continue
		}

		dropTemplate := `DROP TABLE IF EXISTS %s;`
		ddlTemplate := `CREATE TABLE %s (
						seq CHAR (5) PRIMARY KEY ASC UNIQUE,
						%s
					);`
		insertTemplate := "INSERT INTO %s (%s) VALUES %s"

		colDefine := ""
		colList := "seq"
		colCount := 0
		index := 0
		for colIndex, col := range rows[0] {
			colCount++

			val := strings.TrimSpace(col)
			colName := getPinyin(val)
			if colIndex == 0 && colName != "ci" {
				colName = "ci"
			}

			if index > 0 {
				colDefine = colDefine + ",\n"
			}
			colList = colList + ", "

			colType := "VARCHAR"
			if colIndex > 0 {
				colType = "BLOB"
			}
			colDefine = "    " + colDefine + "`" + colName + "` " + colType

			colList = colList + "`" + colName + "`"
			index++
		}

		valList := ""
		for rowIndex, row := range rows {
			if rowIndex == 0 {
				continue
			}

			if rowIndex > 1 {
				valList = valList + ", "
			}
			valList = valList + "(" + strconv.Itoa(rowIndex) + ", "

			for i := 0; i < colCount; i++ {
				if i > 0 {
					valList = valList + ", "
				}

				if i == 0 { // word
					val := strings.TrimSpace(row[i])
					valList = valList + "'" + val + "'"
				} else if i <= len(row) - 1 { // excel value
					str := strings.ToLower(strings.TrimSpace(row[i]))
					if str == "y" {
						val := "TRUE"
						valList = valList + val
					} else {
						valList = valList + "NULL"
					}
				} else {
					valList = valList + "NULL"
				}

			}

			valList = valList + ")"
		}

		dropSql := fmt.Sprintf(dropTemplate, tableName)
		ddl := fmt.Sprintf(ddlTemplate, tableName, colDefine)
		insertSql := fmt.Sprintf(insertTemplate, tableName, colList, valList)

		db, err := sql.Open("sqlite3", constant.SqliteWords)
		defer db.Close()
		_, err = db.Exec(dropSql)
		if err != nil {
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_drop_table", tableName, err.Error()))
			return
		}

		_, err = db.Exec(ddl)
		if err != nil {
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_create_table", tableName, err.Error()))
			return
		}

		_, err = db.Exec(insertSql)
		if err != nil {
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_exec_query", insertSql, err.Error()))
			return
		}
	}
}

func getExcelFilesInDir(folder string, files *[]string) {
	 folder, _ = filepath.Abs(folder)

	if !fileUtils.IsDir(folder) {
		if path.Ext(folder) == ".xlsx" {
			*files = append(*files, folder)
		}

		return
	}

	dir, err := ioutil.ReadDir(folder)
	if err != nil {
		return
	}

	for _, fi := range dir {
		name := fi.Name()
		if commonUtils.IngoreFile(name) {
			continue
		}

		filePath := fileUtils.AddSepIfNeeded(folder) + name
		if fi.IsDir() {
			getExcelFilesInDir(filePath, files)
		} else if strings.Index(name, "~") != 0 && path.Ext(filePath) == ".xlsx" {
			*files = append(*files, filePath)
		}
	}
}

func getPinyin(word string) string {
	p, _ := pinyin.New(word).Split("").Mode(pinyin.WithoutTone).Convert()

	return p
}