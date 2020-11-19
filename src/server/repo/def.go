package serverRepo

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/jinzhu/gorm"
)

type DefRepo struct {
	db *gorm.DB
}

func (r *DefRepo) List() (defs []*model.ZdDef, err error) {
	err = r.db.Find(&defs).Error
	return
}

func (r *DefRepo) Get(id uint) (def model.ZdDef, err error) {
	err = r.db.Where("id=?", id).First(&def).Error

	return
}

func (r *DefRepo) Create(def *model.ZdDef) (err error) {
	err = r.db.Create(def).Error

	return
}

func (r *DefRepo) Update(def *model.ZdDef) (err error) {
	err = r.db.Save(def).Error

	return
}

func (r *DefRepo) Remove(id uint) (err error) {
	var def model.ZdDef
	def.ID = id

	err = r.db.Delete(&def).Error
	err = r.db.Where("defID = ?", id).Delete(&model.ZdField{}).Error

	return
}


func (r *DefRepo) UpdateYaml(def model.ZdDef) (err error) {
	err = r.db.Model(&model.ZdDef{}).Where("id=?", def.ID).Update("yaml", def.Yaml).Error
	return
}

func (r *DefRepo) GenDef(def model.ZdDef, data *model.DefData) () {
	data.Title = def.Title
	data.Desc = def.Desc
	data.Type = def.Type
}

func NewDefRepo(db *gorm.DB) *DefRepo {
	return &DefRepo{db: db}
}
