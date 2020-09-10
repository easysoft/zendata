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

	tableName := "words"
	seq := 1
	ddlFields := make([]string, 0)
	insertSqls := make([]string, 0)

	colMap := map[string]bool{}
	for _, file := range files {
		importExcel(file, tableName, &seq, &ddlFields, &insertSqls, &colMap)
	}

	db, err := sql.Open("sqlite3", constant.SqliteWords)
	defer db.Close()

	dropSql := `DROP TABLE IF EXISTS ` + tableName + `;`
	_, err = db.Exec(dropSql)
	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_drop_table", tableName, err.Error()))
		return
	}

	ddlTemplate := "CREATE TABLE " + tableName + "(\n" +
		"\t`seq` CHAR (5) PRIMARY KEY ASC UNIQUE,\n" +
		"%s" +
		"\n);"
	ddlSql := fmt.Sprintf(ddlTemplate, strings.Join(ddlFields, ", \n"))
	_, err = db.Exec(ddlSql)
	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_create_table", tableName, err.Error()))
		return
	}

	sql := strings.Join(insertSqls, "\n")
	_, err = db.Exec(sql)
	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_exec_query", sql, err.Error()))
		return
	}
}

func importExcel(filePath, tableName string, seq *int, ddlFields, insertSqls *[]string, colMap *map[string]bool) {
	excel, err := excelize.OpenFile(filePath)
	if err != nil {
		logUtils.PrintTo("fail to read file " + filePath + ", error: " + err.Error())
		return
	}

	fileName := path.Base(filePath)
	fileName = strings.TrimSuffix(fileName, path.Ext(filePath))
	fileName = strings.TrimSuffix(fileName, "词库")
	colPrefix := getPinyin(fileName)

	for rowIndex, sheet := range excel.GetSheetList() {
		rows, _ := excel.GetRows(sheet)
		if len(rows) == 0 {
			continue
		}

		colDefine := ""
		colList := make([]string, 0)

		colCount := 0
		index := 0
		for colIndex, col := range rows[0] {
			val := strings.TrimSpace(col)
			if rowIndex == 0 && val == "" {
				break
			}
			colCount++

			colName := getPinyin(val)
			if colIndex == 0 && colName != "ci" {
				colName = "ci"
			} else {
				colName = colPrefix + ":" + colName
			}

			if (*colMap)[colName] == false {
				colType := "VARCHAR"
				colDefine = "    " + "`" + colName + "` " + colType
				*ddlFields = append(*ddlFields, colDefine)

				(*colMap)[colName] = true
			}
			colList = append(colList, "`" + colName + "`")

			index++
		}

		valList := make([]string, 0)
		for rowIndex, row := range rows {
			if rowIndex == 0 {
				continue
			}

			valListItem := make([]string, 0)
			valListItem = append(valListItem, strconv.Itoa(*seq))
			*seq += 1

			for i := 0; i < colCount; i++ {
				val := ""
				if i == 0 { // word
					val = strings.TrimSpace(row[i])
				} else if i <= len(row) - 1 { // excel value
					val = strings.ToLower(strings.TrimSpace(row[i]))
					if val != "y" && val != "b" && val != "f" && val != "m" {
						val = ""
					}
				} else {
					val = ""
				}
				valListItem = append(valListItem,"'" + val + "'")
			}
			valList = append(valList, "(" + strings.Join(valListItem, ", ") + ")")
		}

		insertTemplate := "INSERT INTO `" + tableName + "` (`seq`, %s) VALUES %s;"
		insertSql := fmt.Sprintf(insertTemplate,
			strings.Join(colList, ", "),
			strings.Join(valList, ", "),
			)
		*insertSqls = append(*insertSqls, insertSql)
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