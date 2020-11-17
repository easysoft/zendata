package serverService

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/server/repo"
)

type FieldService struct {
	defRepo *serverRepo.DefRepo
	fieldRepo *serverRepo.FieldRepo
	referRepo *serverRepo.ReferRepo
}

func (s *FieldService) GetTree(defId uint) (root *model.ZdField, err error) {
	root, err = s.fieldRepo.GetDefFieldTree(defId)
	return
}

func (s *FieldService) Get(fieldId int) (field model.ZdField, err error) {
	field, err = s.fieldRepo.Get(uint(fieldId))

	return
}

func (s *FieldService) Save(field *model.ZdField) (err error) {
	err = s.fieldRepo.Save(field)

	return
}
func (s *FieldService) Create(defId, targetId uint, name string, mode string) (field *model.ZdField, err error) {
	field = &model.ZdField{Field: name, DefID: defId}
	if mode == "root" {
		field.DefID = defId
		field.ParentID = 0
	} else {
		var target model.ZdField

		target, err = s.fieldRepo.Get(targetId)
		field.DefID = target.DefID

		if mode == "child" {
			field.ParentID = target.ID
		} else {
			field.ParentID = target.ParentID
		}
		field.Ord = s.fieldRepo.GetMaxOrder(field.ParentID)
	}

	err = s.fieldRepo.Save(field)
	s.referRepo.CreateDefault(field.ID)
	return
}

func (s *FieldService) Remove(id int) (defId int, err error) {
	field, _ := s.fieldRepo.Get(uint(id))
	defId = int(field.DefID)

	err = s.deleteFieldAndChildren(field.DefID, field.ID)
	return
}

func (s *FieldService) Move(srcId, targetId uint, mode string) (defId int, srcField model.ZdField, err error) {
	srcField, err = s.fieldRepo.Get(srcId)
	targetField, err := s.fieldRepo.Get(targetId)

	if "0" == mode {
		srcField.ParentID = targetId
		srcField.Ord = s.fieldRepo.GetMaxOrder(srcField.ParentID)
	} else if "-1" == mode {
		err = s.fieldRepo.AddOrderForTargetAndNextCases(srcField.ID, targetField.Ord, targetField.ParentID)
		if err != nil {
			return
		}

		srcField.ParentID = targetField.ParentID
		srcField.Ord = targetField.Ord
	} else if "1" == mode {
		err = s.fieldRepo.AddOrderForNextCases(srcField.ID, targetField.Ord, targetField.ParentID)
		if err != nil {
			return
		}

		srcField.ParentID = targetField.ParentID
		srcField.Ord = targetField.Ord + 1
	}

	err = s.fieldRepo.UpdateOrdAndParent(srcField)

	return
}

func (s *FieldService) deleteFieldAndChildren(defId, fieldId uint) (err error) {
	err = s.fieldRepo.Remove(fieldId)
	if err == nil {
		children, _ := s.fieldRepo.GetChildren(defId, fieldId)
		for _, child := range children {
			s.deleteFieldAndChildren(child.DefID, child.ID)
		}
	}

	return
}

func NewFieldService(defRepo *serverRepo.DefRepo, fieldRepo *serverRepo.FieldRepo, referRepo *serverRepo.ReferRepo) *FieldService {
	return &FieldService{defRepo: defRepo, fieldRepo: fieldRepo, referRepo: referRepo}
}
