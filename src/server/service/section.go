package serverService

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/server/repo"
)

type SectionService struct {
	fieldRepo   *serverRepo.FieldRepo
	sectionRepo *serverRepo.SectionRepo
}

func (s *SectionService) List(fieldId uint) (sections []*model.Section, err error) {
	sections, err = s.sectionRepo.List(fieldId)
	return
}

func (s *SectionService) Create(fieldId, sectionsId uint) (err error) {
	preSection, err := s.sectionRepo.Get(sectionsId)

	section := model.Section{Value: "0-9", FieldID: fieldId, Ord: preSection.Ord + 1,
		Start: "0", End: "9"}
	err = s.sectionRepo.Create(&section)

	s.fieldRepo.SetIsRange(section.FieldID, true)

	return
}

func (s *SectionService) Update(section *model.Section) (err error) {
	err = s.sectionRepo.Update(section)
	s.fieldRepo.SetIsRange(section.FieldID, true)
	return
}

func (s *SectionService) Remove(sectionId int) (fieldId uint, err error) {
	section, err := s.sectionRepo.Get(uint(sectionId))
	fieldId = section.FieldID

	err = s.sectionRepo.Remove(uint(sectionId))

	s.fieldRepo.SetIsRange(fieldId, true)
	return
}

func NewSectionService(fieldRepo *serverRepo.FieldRepo, sectionRepo *serverRepo.SectionRepo) *SectionService {
	return &SectionService{fieldRepo: fieldRepo, sectionRepo: sectionRepo}
}
