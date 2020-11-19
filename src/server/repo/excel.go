package serverRepo

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/jinzhu/gorm"
)

type ExcelRepo struct {
	db *gorm.DB
}

func (r *ExcelRepo) List() (models []*model.ZdExcel, err error) {
	err = r.db.Select("id,title,name,folder,path").Where("true").Order("id ASC").Find(&models).Error
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
