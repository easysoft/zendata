package serverRepo

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/jinzhu/gorm"
)

type DefRepo struct {
	db *gorm.DB
}

func (r *DefRepo) List() (defs []*model.Def, err error) {
	err = r.db.Find(&defs).Error
	return
}

func (r *DefRepo) Get(id uint) (def model.Def, err error) {
	err = r.db.Where("id=?", id).First(&def).Error

	return
}

func (r *DefRepo) Create(def *model.Def) (err error) {
	err = r.db.Create(def).Error

	return
}

func (r *DefRepo) Update(def *model.Def) (err error) {
	err = r.db.Save(def).Error

	return
}

func (r *DefRepo) Remove(id uint) (err error) {
	var def model.Def
	def.ID = uint(id)

	err = r.db.Delete(&def).Error

	return
}


func (r *DefRepo) UpdateYaml(def model.Def) (err error) {
	err = r.db.Model(&model.Def{}).Where("id=?", def.ID).Update("yaml", def.Yaml).Error
	return
}

func (r *DefRepo) GenDef(def model.Def, data *model.DefData) () {
	data.Title = def.Title
	data.Desc = def.Desc
	data.Type = def.Type
}

func NewDefRepo(db *gorm.DB) *DefRepo {
	return &DefRepo{db: db}
}
