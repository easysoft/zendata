package main

import (
	"fmt"
	"strings"

	"github.com/easysoft/zendata/src/test/import/comm"
	"github.com/easysoft/zendata/src/test/import/model"
	"gorm.io/gorm"
)

func main() {
	filePathArr := []string{
		"./data/words/v1/互联网黑话.xlsx",
		"./data/words/v1/介词词库.xlsx",
		"./data/words/v1/代词词库.xlsx",
		"./data/words/v1/副词词库.xlsx",
		"./data/words/v1/动词词库.xlsx",
		"./data/words/v1/助词词库.xlsx",
		"./data/words/v1/名词词库.xlsx",
		"./data/words/v1/形容词做谓语.xlsx",
		"./data/words/v1/形容词词库.xlsx",
		"./data/words/v1/数词词库.xlsx",
		"./data/words/v1/连词词库.xlsx",
		"./data/words/v1/量词词库.xlsx",
	}

	db := comm.GetDB()
	if err := db.AutoMigrate(
		&model.DataWordTagGroup{}, &model.DataWordTag{}, &model.DataWord{},
	); err != nil {
		panic(err)
	}

	logs := []string{}
	for i, path := range filePathArr {
		sheetName, count := ImoprtWordsFromExcel(path, db)
		str := fmt.Sprintf("%d) Path:[%s] SheetName:[%s], count:[%d]\n", i, path, sheetName, count)
		fmt.Print(str)
		logs = append(logs, str)
	}

	for _, l := range logs {
		fmt.Print(l)
	}
}

func ImoprtWordsFromExcel(filePath string, db *gorm.DB) (sheetName string, count int) {
	sheetName, rows := comm.GetExcel1stSheet(filePath)
	fmt.Printf("importing ... : [%s] \n", sheetName)
	// fmt.Print(rows)

	if len(rows) == 0 {
		return
	}

	headers := rows[0]
	if len(headers) < 1 {
		return
	}

	headers = headers[1:]

	// create word-tag-group
	group := model.DataWordTagGroup{Name: strings.TrimSpace(sheetName)}
	if err := db.Save(&group).Error; err != nil {
		fmt.Errorf("creating word-tag-group : \n%v\n", err)
		return
	}

	// create word-tags
	tags := make([]model.DataWordTag, 0, len(headers))
	for _, v := range headers {
		tag := model.DataWordTag{Name: strings.TrimSpace(v)}

		ret := db.First(&tag, "name = ?", strings.TrimSpace(v))
		if ret.RowsAffected == 0 {
			tag.Groups = []*model.DataWordTagGroup{&group}
			if err := db.Save(&tag).Error; err != nil {
				fmt.Errorf("sheetName:[%s], Tag[%s]", sheetName, v)
				return
			}
		} else {
			if err := db.Model(&tag).Association("Groups").Append(&group); err != nil {
				fmt.Errorf("sheetName:[%s], Tag[%s]", sheetName, v)
				return
			}
		}

		tags = append(tags, tag)
	}

	// create word

	words := make([]model.DataWord, 0, len(rows))

	for i, r := range rows[1:] {
		if len(r) == 0 {
			fmt.Printf("shettName:[%s}, row index:[%d]", sheetName, i+1)
			continue
		}

		word := model.DataWord{Word: strings.TrimSpace(r[0])}

		flag := false
		for j, v := range r[1:] {
			if strings.TrimSpace(v) != "" {
				word.Tags = []*model.DataWordTag{&tags[j]}
				flag = true
				break
			}
		}

		if !flag {
			fmt.Errorf("word (%v) is not tag.", word)
		}

		words = append(words, word)
	}

	tx := db.CreateInBatches(&words, 1000)
	if tx.Error != nil {
		fmt.Errorf("%v", tx.Error)
	}

	var count64 int64
	tx.Count(&count64)
	count = int(count64)

	return
}
