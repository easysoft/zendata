package defServer

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/utils/vari"
)

func ListDefFieldSection(fieldId uint) (sections []*model.Section, err error) {
	err = vari.GormDB.Where("fieldID=?", fieldId).Find(&sections).Error
	return
}

func CreateDefFieldSection(fieldId, sectionsId uint) (err error) {
	var preSection model.Section
	err = vari.GormDB.Where("id=?", sectionsId).Find(&preSection).Error

	section := &model.Section{Value: "0-9", FieldID: fieldId, Ord: preSection.Ord + 1,
		Start: "0", End: "9"}
	err = vari.GormDB.Create(&section).Error
	return
}

func UpdateDefFieldSection(field *model.Section) (err error) {
	err = vari.GormDB.Save(field).Error
	return
}

func RemoveDefFieldSection(sectionId int) (fieldId uint, err error) {
	var section model.Section
	err = vari.GormDB.Where("id=?", sectionId).First(&section).Error
	fieldId = section.FieldID

	err = vari.GormDB.Where("id=?", sectionId).Delete(&model.Section{}).Error
	return
}
