package serverService

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/server/repo"
)

type SectionService struct {
	fieldRepo   *serverRepo.FieldRepo
	instancesRepo   *serverRepo.InstancesRepo

	sectionRepo *serverRepo.SectionRepo
	defService *DefService
	instancesService *InstancesService
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

	if ownerType == "field" {
		s.fieldRepo.SetIsRange(section.OwnerID, true)
	} else if ownerType == "instances" {
		s.instancesRepo.SetIsRange(section.OwnerID, true)
	}

	return
}

func (s *SectionService) Update(section *model.ZdSection) (err error) {
	err = s.sectionRepo.Update(section)

	ownerId := section.OwnerID
	ownerType := section.OwnerType

	s.updateFieldRangeProp(ownerId, ownerType)
	if ownerType == "field" {
		s.fieldRepo.SetIsRange(section.OwnerID, true)
		s.defService.updateYamlByField(section.OwnerID)
	} else if ownerType == "instances" {
		s.instancesRepo.SetIsRange(section.OwnerID, true)
		s.instancesService.updateYamlByItem(section.OwnerID)
	}

	return
}

func (s *SectionService) Remove(sectionId int, ownerType string) (ownerId uint, err error) {
	section, err := s.sectionRepo.Get(uint(sectionId), ownerType)
	ownerId = section.OwnerID

	err = s.sectionRepo.Remove(uint(sectionId), ownerType)


	s.updateFieldRangeProp(ownerId, ownerType)
	if ownerType == "field" {
		s.fieldRepo.SetIsRange(section.OwnerID, true)
		s.defService.updateYamlByField(section.OwnerID)
	} else if ownerType == "instances" {
		s.instancesRepo.SetIsRange(section.OwnerID, true)
		s.instancesService.updateYamlByItem(section.OwnerID)
	}

	return
}

func (s *SectionService) updateFieldRangeProp(ownerId uint, ownerType string) (err error) {
	rangeStr := ""

	sections, _ := s.sectionRepo.List(ownerId, ownerType)
	for index, sect := range sections {
		if index > 0 {
			rangeStr += ","
		}
		rangeStr += sect.Value
	}

	if ownerType == "field" {
		s.fieldRepo.UpdateRange(rangeStr, ownerId)
	} else if ownerType == "instances" {
		s.instancesRepo.UpdateItemRange(rangeStr, ownerId)
	}

	return
}

func NewSectionService(fieldRepo *serverRepo.FieldRepo, instancesRepo *serverRepo.InstancesRepo,
	sectionRepo *serverRepo.SectionRepo, defService *DefService, instancesService *InstancesService) *SectionService {
	return &SectionService{fieldRepo: fieldRepo, sectionRepo: sectionRepo,
		defService: defService, instancesService: instancesService, instancesRepo: instancesRepo}
}
