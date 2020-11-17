package serverRepo

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/jinzhu/gorm"
)

type RangesRepo struct {
	db *gorm.DB
}

func (r *RangesRepo) List() (list []*model.ZdRanges, err error) {
	err = r.db.Where("true").Order("id ASC").Find(&list).Error
	return
}

func (r *RangesRepo) Get(id uint) (ranges model.ZdRanges, err error) {
	err = r.db.Where("id=?", id).First(&ranges).Error
	return
}

func (r *RangesRepo) Save(ranges *model.ZdRanges) (err error) {
	err = r.db.Save(ranges).Error
	return
}

func (r *RangesRepo) Remove(id uint) (err error) {
	model := model.ZdRanges{}
	model.ID = id
	err = r.db.Delete(model).Error

	return
}

func NewRangesRepo(db *gorm.DB) *RangesRepo {
	return &RangesRepo{db: db}
}
