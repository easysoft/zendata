package defServer

import (
	"github.com/easysoft/zendata/src/model"
	commonUtils "github.com/easysoft/zendata/src/utils/common"
	constant "github.com/easysoft/zendata/src/utils/const"
	"github.com/easysoft/zendata/src/utils/vari"
)

func List() (defs []model.Def, err error) {
	err = vari.GormDB.Find(&defs).Error

	return
}
func Get(id int) (def model.Def, err error) {
	vari.GormDB.Where("id=?", id).First(&def)

	return
}

func Create(def *model.Def) (err error) {
	def.Folder = commonUtils.AddPathSep(def.Folder)
	def.Path = def.Folder + constant.PthSep + def.Name
	err = vari.GormDB.Save(def).Error

	return
}

func Update(def *model.Def) (err error) {
	def.Folder = commonUtils.AddPathSep(def.Folder)
	def.Path = def.Folder + constant.PthSep + def.Name
	err = vari.GormDB.Save(def).Error

	return
}

func Remove(id int) (err error) {
	var def model.Def
	def.Id = uint(id)
	err = vari.GormDB.Delete(&def).Error
	return
}
