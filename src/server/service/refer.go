package serverService

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/server/repo"
)

type ReferService struct {
	fieldRepo *serverRepo.FieldRepo
	referRepo *serverRepo.ReferRepo
}

func (s *ReferService) Get(ownerId uint, ownerType string) (refer model.ZdRefer, err error) {
	refer, err = s.referRepo.Get(ownerId, ownerType)
	return
}

func (s *ReferService) Update(ref *model.ZdRefer) (err error) {
	err = s.referRepo.Save(ref)

	s.fieldRepo.SetIsRange(ref.OwnerID, false)

	return
}

func NewReferService(fieldRepo *serverRepo.FieldRepo, referRepo *serverRepo.ReferRepo) *ReferService {
	return &ReferService{fieldRepo: fieldRepo, referRepo: referRepo}
}
