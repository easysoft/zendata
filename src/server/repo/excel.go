package serverRepo

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	"github.com/jinzhu/gorm"
)

type ExcelRepo struct {
	db *gorm.DB
}

func (r *ExcelRepo) ListAll() (models []*model.ZdExcel) {
	r.db.Select("id,title,referName,fileName,folder,path,updatedAt").Find(&models)
	return
}

func (r *ExcelRepo) List(keywords string, page int) (models []*model.ZdExcel, total int, err error) {
	query := r.db.Select("id,title,referName,fileName,folder,path").Order("id ASC")
	if keywords != "" {
		query = query.Where("title LIKE ?", "%"+keywords+"%")
	}
	if page > 0 {
		query = query.Offset((page-1) * constant.PageSize).Limit(constant.PageSize)
	}

	err = query.Find(&models).Error

	err = r.db.Model(&model.ZdExcel{}).Count(&total).Error

	return
}

func (r *ExcelRepo) Get(id uint) (model model.ZdExcel, err error) {
	err = r.db.Where("id=?", id).First(&model).Error
	return
}

func (r *ExcelRepo) Create(model *model.ZdExcel) (err error) {
	err = r.db.Create(model).Error
	return
}
func (r *ExcelRepo) Update(model *model.ZdExcel) (err error) {
	err = r.db.Save(model).Error
	return
}

func (r *ExcelRepo) Remove(id uint) (err error) {
	model := model.ZdExcel{}
	model.ID = id
	err = r.db.Delete(model).Error

	return
}

func NewExcelRepo(db *gorm.DB) *ExcelRepo {
	return &ExcelRepo{db: db}
}
