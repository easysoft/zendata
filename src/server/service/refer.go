package serverService

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/server/repo"
)

type ReferService struct {
	fieldRepo *serverRepo.FieldRepo
	referRepo *serverRepo.ReferRepo

	defService *DefService
}

func (s *ReferService) Get(ownerId uint, ownerType string) (refer model.ZdRefer, err error) {
	refer, err = s.referRepo.GetByOwnerIdAndType(ownerId, ownerType)
	return
}

func (s *ReferService) Update(ref *model.ZdRefer) (err error) {
	err = s.referRepo.Save(ref)

	s.fieldRepo.SetIsRange(ref.OwnerID, false)
	s.defService.updateYamlByField(ref.OwnerID)

	return
}

func NewReferService(fieldRepo *serverRepo.FieldRepo, referRepo *serverRepo.ReferRepo, defService *DefService) *ReferService {
	return &ReferService{fieldRepo: fieldRepo, referRepo: referRepo, defService: defService}
}
