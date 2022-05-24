package serverService

import (
	"github.com/easysoft/zendata/internal/pkg/model"
	"github.com/easysoft/zendata/internal/server/repo"
)

type ReferService struct {
	FieldRepo *serverRepo.FieldRepo `inject:""`
	ReferRepo *serverRepo.ReferRepo `inject:""`

	DefService *DefService `inject:""`
}

func (s *ReferService) Get(ownerId uint, ownerType string) (refer model.ZdRefer, err error) {
	refer, err = s.ReferRepo.GetByOwnerIdAndType(ownerId, ownerType)
	return
}

func (s *ReferService) Update(ref *model.ZdRefer) (err error) {
	err = s.ReferRepo.Save(ref)

	s.FieldRepo.SetIsRange(ref.OwnerID, false)
	s.DefService.updateYamlByField(ref.OwnerID)

	return
}

func NewReferService(fieldRepo *serverRepo.FieldRepo, referRepo *serverRepo.ReferRepo, defService *DefService) *ReferService {
	return &ReferService{FieldRepo: fieldRepo, ReferRepo: referRepo, DefService: defService}
}
