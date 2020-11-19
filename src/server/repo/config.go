package serverRepo

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/jinzhu/gorm"
)

type ConfigRepo struct {
	db *gorm.DB
}

func (r *ConfigRepo) List() (models []*model.ZdConfig, err error) {
	err = r.db.Select("id,title,name,folder,path").Where("true").Order("id ASC").Find(&models).Error
	return
}

func (r *ConfigRepo) Get(id uint) (model model.ZdConfig, err error) {
	err = r.db.Where("id=?", id).First(&model).Error
	return
}

func (r *ConfigRepo) Create(model *model.ZdConfig) (err error) {
	err = r.db.Create(model).Error
	return
}
func (r *ConfigRepo) Update(model *model.ZdConfig) (err error) {
	err = r.db.Update(model).Error
	return
}

func (r *ConfigRepo) Remove(id uint) (err error) {
	model := model.ZdConfig{}
	model.ID = id
	err = r.db.Delete(model).Error

	return
}

func NewConfigRepo(db *gorm.DB) *ConfigRepo {
	return &ConfigRepo{db: db}
}
