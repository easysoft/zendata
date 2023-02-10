package serverRepo

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/domain"
	"github.com/easysoft/zendata/internal/pkg/model"
	"gorm.io/gorm"
)

type DefRepo struct {
	DB *gorm.DB `inject:""`
}

func (r *DefRepo) ListAll() (models []*model.ZdDef) {
	r.DB.Select("id,title,folder,path,referName,updatedAt").Find(&models)
	return
}

func (r *DefRepo) List(keywords string, page int) (models []*model.ZdDef, total int, err error) {
	query := r.DB.Select("id,title,folder,path,referName").Order("id ASC")
	if keywords != "" {
		query = query.Where("title LIKE ?", "%"+keywords+"%")
	}
	if page > 0 {
		query = query.Offset((page - 1) * consts.PageSize).Limit(consts.PageSize)
	}

	err = query.Find(&models).Error

	var total64 int64
	err = r.DB.Model(&model.ZdDef{}).Count(&total64).Error
	total = int(total64)

	return
}

func (r *DefRepo) Get(id uint) (def model.ZdDef, err error) {
	err = r.DB.Where("id=?", id).First(&def).Error

	return
}

func (r *DefRepo) Create(def *model.ZdDef) (err error) {
	err = r.DB.Create(def).Error

	return
}

func (r *DefRepo) Update(def *model.ZdDef) (err error) {
	err = r.DB.Save(def).Error

	return
}

func (r *DefRepo) Remove(id uint) (err error) {
	var def model.ZdDef
	def.ID = id

	err = r.DB.Delete(&def).Error
	err = r.DB.Where("defID = ?", id).Delete(&model.ZdField{}).Error

	return
}

func (r *DefRepo) UpdateYaml(po model.ZdDef) (err error) {
	err = r.DB.Model(&model.ZdDef{}).Where("id=?", po.ID).Update("yaml", po.Yaml).Error
	return
}

func (r *DefRepo) GenDef(def model.ZdDef, data *domain.DefData) {
	data.Title = def.Title
	data.Desc = def.Desc
	data.Type = def.Type
	if data.Type == consts.ResTypeText {
		data.Type = ""
	}
}

func NewDefRepo(db *gorm.DB) *DefRepo {
	return &DefRepo{DB: db}
}
