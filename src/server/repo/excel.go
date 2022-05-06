package serverRepo

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	"gorm.io/gorm"
)

type ExcelRepo struct {
	DB *gorm.DB `inject:""`
}

func (r *ExcelRepo) ListAll() (models []*model.ZdExcel) {
	r.DB.Select("id,title,referName,fileName,folder,path,updatedAt").Find(&models)
	return
}

func (r *ExcelRepo) List(keywords string, page int) (models []*model.ZdExcel, total int, err error) {
	query := r.DB.Select("id,title,referName,fileName,folder,path").Order("id ASC")
	if keywords != "" {
		query = query.Where("title LIKE ?", "%"+keywords+"%")
	}
	if page > 0 {
		query = query.Offset((page - 1) * constant.PageSize).Limit(constant.PageSize)
	}

	err = query.Find(&models).Error

	var total64 int64
	err = r.DB.Model(&model.ZdExcel{}).Count(&total64).Error

	total = int(total64)
	return
}

func (r *ExcelRepo) ListFiles() (models []*model.ZdExcel) {
	tbl := (&model.ZdExcel{}).TableName()
	sql := "select id,title,referName,fileName,folder,path from " +
		tbl + " where id in (select max(id) from " + tbl + " group by referName)"
	r.DB.Raw(sql).Find(&models)
	return
}
func (r *ExcelRepo) ListSheets(referName string) (models []*model.ZdExcel) {
	r.DB.Select("id,sheet").Where("referName=?", referName).Find(&models)
	return
}

func (r *ExcelRepo) Get(id uint) (model model.ZdExcel, err error) {
	err = r.DB.Where("id=?", id).First(&model).Error
	return
}
func (r *ExcelRepo) GetBySheet(referName, sheet string) (model model.ZdExcel, err error) {
	err = r.DB.Where("referName=? AND sheet=?", referName, sheet).First(&model).Error
	return
}

func (r *ExcelRepo) Create(model *model.ZdExcel) (err error) {
	err = r.DB.Create(model).Error
	return
}
func (r *ExcelRepo) Update(model *model.ZdExcel) (err error) {
	err = r.DB.Save(model).Error
	return
}

func (r *ExcelRepo) Remove(id uint) (err error) {
	model := model.ZdExcel{}
	model.ID = id
	err = r.DB.Delete(model).Error

	return
}

func NewExcelRepo(db *gorm.DB) *ExcelRepo {
	return &ExcelRepo{DB: db}
}
