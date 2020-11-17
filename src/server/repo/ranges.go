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

func (r *RangesRepo) GetItems(rangesId int) (items []*model.ZdRangesItem, err error) {
	err = r.db.Where("rangesId=?", rangesId).Find(&items).Error
	return
}
func (r *RangesRepo) GetItem(itemId uint) (item model.ZdRangesItem, err error) {
	err = r.db.Where("id=?", itemId).First(&item).Error
	return
}
func (r *RangesRepo) SaveItem(item *model.ZdRangesItem) (err error) {
	err = r.db.Save(item).Error
	return
}

func (r *RangesRepo) GetMaxOrder(rangesId int) (ord int) {
	var preChild model.ZdField
	err := r.db.
		Where("rangesID=?", rangesId).
		Order("ord DESC").Limit(1).
		First(&preChild).Error

	if err != nil {
		ord = 1
	}
	ord = preChild.Ord + 1

	return
}

func (r *RangesRepo) RemoveItem(id uint) (err error) {
	item := model.ZdRangesItem{}
	item.ID = id
	err = r.db.Delete(item).Error
	return
}

func NewRangesRepo(db *gorm.DB) *RangesRepo {
	return &RangesRepo{db: db}
}
