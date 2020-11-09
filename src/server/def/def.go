package defServer

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/jinzhu/gorm"
	"strings"
)

func List() (defs []model.Def, err error) {
	err = vari.GormDB.Find(&defs).Error

	return
}
func Get(id int) (def model.Def, err error) {
	vari.GormDB.Where("id=?", id).First(&def)
	def.Folder = GetFolder(def.Path)

	return
}

func Create(def *model.Def) (err error) {
	def.Folder = DealWithPathSepRight(def.Folder)

	def.Path = def.Folder + def.Name
	def.Path = AddExt(def.Path)
	err = vari.GormDB.Save(def).Error

	return
}

func Update(def *model.Def) (err error) {
	def.Folder = DealWithPathSepRight(def.Folder)

	def.Path = def.Folder + def.Name
	def.Path = AddExt(def.Path)

	var oldDef model.Def
	err = vari.GormDB.Where("id=?", def.Id).First(&oldDef).Error
	if err == gorm.ErrRecordNotFound {
		return
	}

	if def.Path != oldDef.Path {
		fileUtils.RemoveExist(oldDef.Path)
	}

	err = vari.GormDB.Save(def).Error

	return
}

func Remove(id int) (err error) {
	var oldDef model.Def
	err = vari.GormDB.Where("id=?", id).First(&oldDef).Error
	if err == gorm.ErrRecordNotFound {
		return
	}

	fileUtils.RemoveExist(oldDef.Path)

	var def model.Def
	def.Id = uint(id)
	err = vari.GormDB.Delete(&def).Error
	return
}

func AddExt(pth string) string {
	if strings.LastIndex(pth, ".yaml") != len(pth) - 4 {
		pth += ".yaml"
	}

	return pth
}

func DealWithPathSepRight(pth string) string {
	pth = fileUtils.RemovePathSepLeftIfNeeded(pth)
	pth = fileUtils.AddPathSepRightIfNeeded(pth)

	return pth
}
func GetFolder(pth string) string {
	idx := strings.LastIndex(pth, constant.PthSep)
	return pth[:idx+1]
}
