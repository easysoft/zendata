package serverService

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/server/repo"
)

type SectionService struct {
	fieldRepo   *serverRepo.FieldRepo
	sectionRepo *serverRepo.SectionRepo
}

func (s *SectionService) List(ownerId uint, ownerType string) (sections []*model.ZdSection, err error) {
	sections, err = s.sectionRepo.List(ownerId, ownerType)
	return
}

func (s *SectionService) Create(ownerId, sectionsId uint, ownerType string) (err error) {
	preSection, err := s.sectionRepo.Get(sectionsId, ownerType)

	section := model.ZdSection{Value: "0-9", OwnerID: ownerId, OwnerType: ownerType, Ord: preSection.Ord + 1,
		Start: "0", End: "9"}
	err = s.sectionRepo.Create(&section)

	s.fieldRepo.SetIsRange(section.OwnerID, true)

	return
}

func (s *SectionService) Update(section *model.ZdSection) (err error) {
	err = s.sectionRepo.Update(section)
	s.fieldRepo.SetIsRange(section.OwnerID, true)
	return
}

func (s *SectionService) Remove(sectionId int, ownerType string) (ownerId uint, err error) {
	section, err := s.sectionRepo.Get(uint(sectionId), ownerType)
	ownerId = section.OwnerID

	err = s.sectionRepo.Remove(uint(sectionId), ownerType)

	s.fieldRepo.SetIsRange(ownerId, true)
	return
}

func NewSectionService(fieldRepo *serverRepo.FieldRepo, sectionRepo *serverRepo.SectionRepo) *SectionService {
	return &SectionService{fieldRepo: fieldRepo, sectionRepo: sectionRepo}
}
