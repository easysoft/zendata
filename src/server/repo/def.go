package serverRepo

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	"github.com/jinzhu/gorm"
)

type DefRepo struct {
	db *gorm.DB
}

func (r *DefRepo) ListAll() (models []*model.ZdDef) {
	r.db.Select("id,title,folder,path,referName,updatedAt").Find(&models)
	return
}

func (r *DefRepo) List(keywords string, page int) (models []*model.ZdDef, total int, err error) {
	query := r.db.Select("id,title,folder,path,referName").Order("id ASC")
	if keywords != "" {
		query = query.Where("title LIKE ?", "%"+keywords+"%")
	}
	if page > 0 {
		query = query.Offset((page-1) * constant.PageSize).Limit(constant.PageSize)
	}

	err = query.Find(&models).Error

	err = r.db.Model(&model.ZdDef{}).Count(&total).Error

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

func (r *DefRepo) UpdateYaml(po model.ZdDef) (err error) {
	err = r.db.Model(&model.ZdDef{}).Where("id=?", po.ID).Update("yaml", po.Yaml).Error
	return
}

func (r *DefRepo) GenDef(def model.ZdDef, data *model.DefData) () {
	data.Title = def.Title
	data.Desc = def.Desc
	data.Type = def.Type
	if data.Type == constant.ResTypeText {
		data.Type = ""
	}
}

func NewDefRepo(db *gorm.DB) *DefRepo {
	return &DefRepo{db: db}
}
