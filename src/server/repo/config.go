package serverRepo

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	"gorm.io/gorm"
)

type ConfigRepo struct {
	DB *gorm.DB `inject:""`
}

func (r *ConfigRepo) ListAll() (models []*model.ZdConfig) {
	r.DB.Select("id,title,referName,fileName,folder,path,updatedAt").Find(&models)
	return
}

func (r *ConfigRepo) List(keywords string, page int) (models []*model.ZdConfig, total int, err error) {
	query := r.DB.Select("id,title,referName,fileName,folder,path").Order("id ASC")
	if keywords != "" {
		query = query.Where("title LIKE ?", "%"+keywords+"%")
	}
	if page > 0 {
		query = query.Offset((page - 1) * constant.PageSize).Limit(constant.PageSize)
	}

	err = query.Find(&models).Error

	var total64 int64
	err = r.DB.Model(&model.ZdConfig{}).Count(&total64).Error
	total = int(total64)

	return
}

func (r *ConfigRepo) Get(id uint) (model model.ZdConfig, err error) {
	err = r.DB.Where("id=?", id).First(&model).Error
	return
}

func (r *ConfigRepo) Create(model *model.ZdConfig) (err error) {
	err = r.DB.Create(model).Error
	return
}
func (r *ConfigRepo) Update(model *model.ZdConfig) (err error) {
	err = r.DB.Save(model).Error
	return
}

func (r *ConfigRepo) Remove(id uint) (err error) {
	model := model.ZdConfig{}
	model.ID = id
	err = r.DB.Delete(model).Error

	return
}

func (r *ConfigRepo) UpdateYaml(po model.ZdConfig) (err error) {
	err = r.DB.Model(&model.ZdConfig{}).Where("id=?", po.ID).Update("yaml", po.Yaml).Error
	return
}

func (r *ConfigRepo) GenConfigRes(config model.ZdConfig, res *model.ResConfig) {
	res.Title = config.Title
	res.Desc = config.Desc
	res.Prefix = config.Prefix
	res.Postfix = config.Postfix
	res.Loop = config.Loop
	res.Loopfix = config.Loopfix
	res.Format = config.Format

	res.Range = config.Range
}

func (r *ConfigRepo) UpdateConfigRange(rang string, id uint) (err error) {
	err = r.DB.Model(&model.ZdConfig{}).Where("id=?", id).Update("range", rang).Error

	return
}

func NewConfigRepo(db *gorm.DB) *ConfigRepo {
	return &ConfigRepo{DB: db}
}
