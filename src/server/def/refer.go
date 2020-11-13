package defServer

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/utils/vari"
)

func CreateDefFieldRefer(fieldId uint) (err error) {
	refer := &model.Refer{FieldID: fieldId}
	err = vari.GormDB.Create(&refer).Error

	return
}

func GetDefFieldRefer(fieldId uint) (refer model.Refer, err error) {
	err = vari.GormDB.Where("fieldID=?", fieldId).First(&refer).Error
	return
}

func UpdateDefFieldRefer(ref *model.Refer) (err error) {
	err = vari.GormDB.Save(ref).Error

	setDefFieldIsRange(ref.FieldID, false)

	return
}
