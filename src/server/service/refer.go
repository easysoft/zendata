package serverService

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/server/repo"
)

type ReferService struct {
	fieldRepo *serverRepo.FieldRepo
	referRepo *serverRepo.ReferRepo
}

func (s *ReferService) CreateDefault(fieldId uint) (err error) {
	refer := &model.ZdRefer{FieldID: fieldId}
	err = s.referRepo.Create(refer)

	return
}

func (s *ReferService) Get(fieldId uint) (refer model.ZdRefer, err error) {
	refer, err = s.referRepo.Get(fieldId)
	return
}

func (s *ReferService) Update(ref *model.ZdRefer) (err error) {
	err = s.referRepo.Save(ref)

	s.fieldRepo.SetIsRange(ref.FieldID, false)

	return
}

func NewReferService(fieldRepo *serverRepo.FieldRepo, referRepo *serverRepo.ReferRepo) *ReferService {
	return &ReferService{fieldRepo: fieldRepo, referRepo: referRepo}
}
