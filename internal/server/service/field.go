package serverService

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	serverRepo "github.com/easysoft/zendata/internal/server/repo"
)

type FieldService struct {
	DefRepo   *serverRepo.DefRepo   `inject:""`
	FieldRepo *serverRepo.FieldRepo `inject:""`
	ReferRepo *serverRepo.ReferRepo `inject:""`

	DefService *DefService `inject:""`
}

func (s *FieldService) GetTree(defId uint) (root *model.ZdField, err error) {
	root, err = s.FieldRepo.GetDefFieldTree(defId)
	return
}

func (s *FieldService) Get(fieldId int) (field model.ZdField, err error) {
	field, err = s.FieldRepo.Get(uint(fieldId))

	return
}

func (s *FieldService) Save(field *model.ZdField) (err error) {
	err = s.FieldRepo.Save(field)
	s.DefService.updateYaml(field.DefID)

	return
}
func (s *FieldService) Create(defId, targetId uint, name string, mode string) (field *model.ZdField, err error) {
	field = &model.ZdField{DefID: defId}
	field.Field = name
	if mode == "root" {
		field.DefID = defId
		field.ParentID = 0
	} else {
		var target model.ZdField

		target, err = s.FieldRepo.Get(targetId)
		field.DefID = target.DefID

		if mode == "child" {
			field.ParentID = target.ID
		} else {
			field.ParentID = target.ParentID
		}
		field.Ord = s.FieldRepo.GetMaxOrder(field.ParentID)
	}

	err = s.FieldRepo.Save(field)
	s.ReferRepo.CreateDefault(field.ID, consts.ResTypeDef)

	s.DefService.updateYaml(field.DefID)

	return
}

func (s *FieldService) Remove(id int) (defId int, err error) {
	field, _ := s.FieldRepo.Get(uint(id))
	defId = int(field.DefID)

	err = s.deleteFieldAndChildren(field.DefID, field.ID)

	s.DefService.updateYaml(field.DefID)

	return
}

func (s *FieldService) Move(srcId, targetId uint, mode string) (defId uint, srcField model.ZdField, err error) {
	srcField, err = s.FieldRepo.Get(srcId)
	targetField, err := s.FieldRepo.Get(targetId)
	defId = srcField.DefID

	if "0" == mode {
		srcField.ParentID = targetId
		srcField.Ord = s.FieldRepo.GetMaxOrder(srcField.ParentID)
	} else if "-1" == mode {
		err = s.FieldRepo.AddOrderForTargetAndNextCases(srcField.ID, targetField.Ord, targetField.ParentID)
		if err != nil {
			return
		}

		srcField.ParentID = targetField.ParentID
		srcField.Ord = targetField.Ord
	} else if "1" == mode {
		err = s.FieldRepo.AddOrderForNextCases(srcField.ID, targetField.Ord, targetField.ParentID)
		if err != nil {
			return
		}

		srcField.ParentID = targetField.ParentID
		srcField.Ord = targetField.Ord + 1
	}

	err = s.FieldRepo.UpdateOrdAndParent(srcField)

	s.DefService.updateYaml(defId)

	return
}

func (s *FieldService) deleteFieldAndChildren(defId, fieldId uint) (err error) {
	err = s.FieldRepo.Remove(fieldId)
	if err == nil {
		children, _ := s.FieldRepo.GetChildren(defId, fieldId)
		for _, child := range children {
			s.deleteFieldAndChildren(child.DefID, child.ID)
		}
	}

	return
}

func NewFieldService(defRepo *serverRepo.DefRepo, fieldRepo *serverRepo.FieldRepo,
	referRepo *serverRepo.ReferRepo,
	defService *DefService) *FieldService {
	return &FieldService{DefRepo: defRepo, FieldRepo: fieldRepo, ReferRepo: referRepo,
		DefService: defService}
}
