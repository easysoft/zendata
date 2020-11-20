package serverRepo

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	"github.com/jinzhu/gorm"
)

type ConfigRepo struct {
	db *gorm.DB
}

func (r *ConfigRepo) ListAll() (models []*model.ZdConfig) {
	r.db.Select("id,title,name,folder,path,updatedAt").Find(&models)
	return
}

func (r *ConfigRepo) List(keywords string, page int) (models []*model.ZdConfig, total int, err error) {
	query := r.db.Select("id,title,name,folder,path").Order("id ASC")
	if keywords != "" {
		query = query.Where("title LIKE ?", "%"+keywords+"%")
	}
	if page > 0 {
		query = query.Offset((page-1) * constant.PageSize).Limit(constant.PageSize)
	}

	err = query.Find(&models).Error

	err = r.db.Model(&model.ZdConfig{}).Count(&total).Error

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
	err = r.db.Save(model).Error
	return
}

func (r *ConfigRepo) Remove(id uint) (err error) {
	model := model.ZdConfig{}
	model.ID = id
	err = r.db.Delete(model).Error

	return
}

func (r *ConfigRepo) UpdateYaml(po model.ZdConfig) (err error) {
	err = r.db.Model(&model.ZdConfig{}).Where("id=?", po.ID).Update("yaml", po.Yaml).Error
	return
}

func NewConfigRepo(db *gorm.DB) *ConfigRepo {
	return &ConfigRepo{db: db}
}
