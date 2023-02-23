package serverRepo

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	"gorm.io/gorm"
)

type TextRepo struct {
	DB *gorm.DB `inject:""`
}

func (r *TextRepo) ListAll() (models []*model.ZdText) {
	r.DB.Select("id,title,referName,fileName,folder,path,updatedAt").Find(&models)
	return
}

func (r *TextRepo) List(keywords string, page int) (models []*model.ZdText, total int, err error) {
	query := r.DB.Select("id,title,referName,fileName,folder,path").Order("id ASC")
	if keywords != "" {
		query = query.Where("title LIKE ?", "%"+keywords+"%")
	}
	if page > 0 {
		query = query.Offset((page - 1) * consts.PageSize).Limit(consts.PageSize)
	}

	err = query.Find(&models).Error

	var total64 int64
	err = r.DB.Model(&model.ZdText{}).Count(&total64).Error
	total = int(total64)

	return
}

func (r *TextRepo) Get(id uint) (model model.ZdText, err error) {
	err = r.DB.Where("id=?", id).First(&model).Error
	return
}

func (r *TextRepo) Create(model *model.ZdText) (err error) {
	err = r.DB.Create(model).Error
	return
}
func (r *TextRepo) Update(model *model.ZdText) (err error) {
	err = r.DB.Save(model).Error
	return
}

func (r *TextRepo) Remove(id uint) (err error) {
	model := model.ZdText{}
	model.ID = id
	err = r.DB.Delete(model).Error

	return
}

func NewTextRepo(db *gorm.DB) *TextRepo {
	return &TextRepo{DB: db}
}
