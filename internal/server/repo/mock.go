package serverRepo

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	"gorm.io/gorm"
)

type MockRepo struct {
	DB *gorm.DB `inject:""`
}

func (r *MockRepo) List(keywords string, page int) (pos []*model.ZdMock, total int, err error) {
	query := r.DB
	if keywords != "" {
		query = query.Where("title LIKE ?", "%"+keywords+"%")
	}
	if page > 0 {
		query = query.Offset((page - 1) * consts.PageSize).Limit(consts.PageSize)
	}

	err = query.Find(&pos).Error

	var total64 int64
	err = r.DB.Model(&model.ZdMock{}).Count(&total64).Error
	total = int(total64)

	return
}

func (r *MockRepo) Get(id uint) (po model.ZdMock, err error) {
	err = r.DB.Where("id=?", id).First(&po).Error

	return
}

func (r *MockRepo) GetByPath(path string) (po model.ZdMock, err error) {
	err = r.DB.Where("path=?", path).First(&po).Error

	return
}

func (r *MockRepo) GetByDefID(defId uint) (po model.ZdMock) {
	r.DB.Where("def_id=?", defId).First(&po)

	return
}

func (r *MockRepo) Save(po *model.ZdMock) (err error) {
	err = r.DB.Save(po).Error

	return
}

func (r *MockRepo) Remove(id uint) (err error) {
	var po model.ZdMock
	po.ID = id

	err = r.DB.Delete(&po).Error
	err = r.DB.Where("id = ?", id).Delete(&model.ZdField{}).Error

	return
}

func (r *MockRepo) ListSampleSrc(mockId int) (pos []model.ZdMockSampleSrc, err error) {
	err = r.DB.Find(&pos).Error
	return
}

func (r *MockRepo) GetSampleSrc(mockId uint, key string) (po model.ZdMockSampleSrc, err error) {
	err = r.DB.Where("mock_id=? AND key=?", mockId, key).
		Find(&po).Error
	return
}

func (r *MockRepo) ChangeSampleSrc(mockId int, req model.ZdMockSampleSrc) (err error) {
	po := model.ZdMockSampleSrc{}

	err = r.DB.Where("mock_id=? AND key=?", mockId, req.Key).First(&po).Error
	if err != nil { // not found
		po.MockId = uint(mockId)
		po.Key = req.Key
	}

	po.Value = req.Value
	err = r.DB.Save(&po).Error

	return
}
