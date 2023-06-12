package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/easysoft/zendata/cmd/test/others/func/comm"
	"github.com/easysoft/zendata/cmd/test/others/func/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"gorm.io/gorm"
	"path/filepath"
)

func main() {
	db := comm.GetDB()
	db.AutoMigrate(
		&model.DataWord{},
	)

	// load tag groups
	groups := make([]model.DataWordTagGroup, 0)
	db.Order("id ASC").Find(&groups)

	// gen sheet by tag group
	for _, group := range groups {
		sheetName := group.Name
		filePath := fmt.Sprintf("data/words/v1/%s.xlsx", sheetName)

		fileUtils.MkDirIfNeeded(filepath.Dir(filePath))

		f := excelize.NewFile()
		index := f.NewSheet(sheetName)
		f.SetActiveSheet(index)
		sheet1 := f.GetSheetName(0)
		f.DeleteSheet(sheet1)

		allTags := loadTagsByGroup(group.Id, db)
		words := loadWordsByGroup(group.Id, db)

		// gen headers
		excelColIndex := 'A'
		excelColId := fmt.Sprintf("%c%d", excelColIndex, 1)
		f.SetCellValue(sheetName, excelColId, "词")
		excelColIndex++

		for _, tag := range allTags {
			excelColId := fmt.Sprintf("%c%d", excelColIndex, 1)
			f.SetCellValue(sheetName, excelColId, tag.Name)

			excelColIndex++
		}

		// gen word row
		wordIndex := 2
		for _, word := range words {
			tags := loadTagsByWord(word.Id, group.Id, db)

			// gen word's tag data
			wordTagMap := map[string]bool{}
			for _, tag := range tags {
				wordTagMap[tag.Name] = true
			}

			// gen rows
			excelColIndex := 'A'
			excelColId := fmt.Sprintf("%c%d", excelColIndex, wordIndex)
			f.SetCellValue(sheetName, excelColId, word.Word)
			excelColIndex++

			for _, tag := range allTags {
				excelColId := fmt.Sprintf("%c%d", excelColIndex, wordIndex)

				val := ""
				if wordTagMap[tag.Name] {
					val = "Y"
				}
				f.SetCellValue(sheetName, excelColId, val)

				excelColIndex++
			}

			wordIndex++
		}

		f.SaveAs(filePath)
	}
}

func loadTagsByWord(wordId uint, groupId uint, db *gorm.DB) (tags []model.DataWordTag) {
	tagIds := make([]int, 0)

	sqlTags := fmt.Sprintf("SELECT r.data_word_tag_id "+
		"FROM biz_data_word_biz_data_word_tag r "+
		"WHERE r.data_word_id = %d AND r.data_word_tag_id IN "+
		"	(SELECT data_word_tag_id FROM zendata.biz_data_word_tag_group_biz_data_word_tag "+
		"		WHERE data_word_tag_group_id = %d) ",
		wordId, groupId)

	db.Raw(sqlTags).Scan(&tagIds)

	db.Where("id IN (?) AND NOT deleted", tagIds).Find(&tags)

	return
}

func loadTagsByGroup(groupId uint, db *gorm.DB) (tags []model.DataWordTag) {
	sqlTags := fmt.Sprintf("SELECT t.* "+
		"FROM biz_data_word_tag t "+
		"WHERE t.id IN "+
		"	(SELECT data_word_tag_id FROM zendata.biz_data_word_tag_group_biz_data_word_tag "+
		"		WHERE data_word_tag_group_id = %d) "+
		"AND NOT t.deleted "+
		"ORDER BY t.id ",
		groupId)

	db.Raw(sqlTags).Scan(&tags)

	return
}

func loadWordsByGroup(groupId uint, db *gorm.DB) (words []model.DataWord) {
	sqlWords := fmt.Sprintf("SELECT w.* FROM biz_data_word w "+
		"WHERE w.tag_group_id = %d AND NOT w.deleted "+
		"ORDER BY w.id",
		groupId)
	db.Raw(sqlWords).Scan(&words)

	return
}
