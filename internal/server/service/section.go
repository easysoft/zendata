package serverService

import (
	"github.com/easysoft/zendata/internal/pkg/model"
	"github.com/easysoft/zendata/internal/pkg/service"
	"github.com/easysoft/zendata/internal/server/repo"
)

type SectionService struct {
	FieldRepo     *serverRepo.FieldRepo     `inject:""`
	ConfigRepo    *serverRepo.ConfigRepo    `inject:""`
	RangesRepo    *serverRepo.RangesRepo    `inject:""`
	InstancesRepo *serverRepo.InstancesRepo `inject:""`

	SectionRepo      *serverRepo.SectionRepo `inject:""`
	DefService       *DefService             `inject:""`
	ConfigService    *ConfigService          `inject:""`
	InstancesService *InstancesService       `inject:""`
	RangesService    *RangesService          `inject:""`

	RangesService2 *service.RangeService `inject:""`
}

func (s *SectionService) List(ownerId uint, ownerType string) (sections []*model.ZdSection, err error) {
	sections, err = s.SectionRepo.List(ownerId, ownerType)
	return
}

func (s *SectionService) Create(ownerId, sectionsId uint, ownerType string) (err error) {
	preSection, err := s.SectionRepo.Get(sectionsId, ownerType)

	section := model.ZdSection{Value: "0-9", OwnerID: ownerId, OwnerType: ownerType, Ord: preSection.Ord + 1,
		Start: "0", End: "9"}
	err = s.SectionRepo.Create(&section)

	if ownerType == "def" {
		s.FieldRepo.SetIsRange(section.OwnerID, true)
	} else if ownerType == "instances" {
		s.InstancesRepo.SetIsRange(section.OwnerID, true)
	}

	return
}

func (s *SectionService) Update(section *model.ZdSection) (err error) {
	err = s.SectionRepo.Update(section)

	ownerId := section.OwnerID
	ownerType := section.OwnerType

	s.updateFieldRangeProp(ownerId, ownerType)
	if ownerType == "def" {
		s.FieldRepo.SetIsRange(section.OwnerID, true)
		s.DefService.updateYamlByField(section.OwnerID)

	} else if ownerType == "config" {
		s.ConfigService.updateYaml(section.OwnerID)
	} else if ownerType == "ranges" {
		s.RangesService.updateYamlByItem(section.OwnerID)

	} else if ownerType == "instances" {
		s.InstancesRepo.SetIsRange(section.OwnerID, true)
		s.InstancesService.updateYamlByItem(section.OwnerID)
	}

	return
}

func (s *SectionService) Remove(sectionId int, ownerType string) (ownerId uint, err error) {
	section, err := s.SectionRepo.Get(uint(sectionId), ownerType)
	ownerId = section.OwnerID

	err = s.SectionRepo.Remove(uint(sectionId), ownerType)

	s.updateFieldRangeProp(ownerId, ownerType)
	if ownerType == "def" {
		s.FieldRepo.SetIsRange(section.OwnerID, true)
		s.DefService.updateYamlByField(section.OwnerID)
	} else if ownerType == "config" {
		s.ConfigService.updateYaml(section.OwnerID)
	} else if ownerType == "ranges" {
		s.RangesService.updateYamlByItem(section.OwnerID)

	} else if ownerType == "instances" {
		s.InstancesRepo.SetIsRange(section.OwnerID, true)
		s.InstancesService.updateYamlByItem(section.OwnerID)
	}

	return
}

func (s *SectionService) updateFieldRangeProp(ownerId uint, ownerType string) (err error) {
	rangeStr := ""

	sections, _ := s.SectionRepo.List(ownerId, ownerType)
	for index, sect := range sections {
		if index > 0 {
			rangeStr += ","
		}
		rangeStr += sect.Value
	}

	if ownerType == "def" {
		s.FieldRepo.UpdateRange(rangeStr, ownerId)
	} else if ownerType == "config" {
		s.ConfigRepo.UpdateConfigRange(rangeStr, ownerId)
	} else if ownerType == "ranges" {
		s.RangesRepo.UpdateItemRange(rangeStr, ownerId)
	} else if ownerType == "instances" {
		s.InstancesRepo.UpdateItemRange(rangeStr, ownerId)
	}

	return
}

func (s *SectionService) SaveFieldSectionToDB(rangeSection string, ord int, fieldID uint, ownerType string) {
	descStr, stepStr, count, countTag := s.RangesService2.ParseRangeSection(rangeSection)
	typ, desc := s.RangesService2.ParseRangeSectionDesc(descStr)

	s.SectionRepo.SaveFieldSectionToDB(typ, desc, stepStr, count, countTag,
		ord, fieldID, ownerType)
}
