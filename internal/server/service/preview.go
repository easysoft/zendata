package serverService

import (
	"fmt"
	"github.com/easysoft/zendata/internal/pkg/action"
	constant "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	"github.com/easysoft/zendata/internal/server/repo"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/jinzhu/copier"
	"gopkg.in/yaml.v3"
)

type PreviewService struct {
	DefRepo       *serverRepo.DefRepo       `inject:""`
	FieldRepo     *serverRepo.FieldRepo     `inject:""`
	ReferRepo     *serverRepo.ReferRepo     `inject:""`
	InstancesRepo *serverRepo.InstancesRepo `inject:""`
}

func (s *PreviewService) PreviewDefData(defId uint) (data string) {
	def, _ := s.DefRepo.Get(defId)

	vari.GenVars.Total = 10
	lines := action.Generate([]string{def.Path}, "", constant.FormatData, "")
	data = s.linesToStr(lines)

	return
}
func (s *PreviewService) PreviewFieldData(fieldId uint, fieldType string) (data string) {
	var field model.ZdField

	if fieldType == constant.ResTypeDef {
		field, _ = s.FieldRepo.Get(fieldId)
	} else if fieldType == constant.ResTypeInstances {
		instItem, _ := s.InstancesRepo.GetItem(fieldId)
		field.From = instItem.From
		copier.Copy(&field, instItem)
	}

	ref := model.ZdRefer{}
	if !field.IsRange {
		ref, _ = s.ReferRepo.GetByOwnerId(field.ID)
	}

	fld := model.DefField{}
	genFieldFromZdField(field, ref, &fld)

	def := model.DefData{}
	def.Fields = append(def.Fields, fld)
	defContent, _ := yaml.Marshal(def)

	configFile := vari.ZdPath + "tmp" + constant.PthSep + ".temp.yaml"
	fileUtils.WriteFile(configFile, string(defContent))

	lines := action.Generate([]string{configFile}, field.Field, constant.FormatData, "")
	data = s.linesToStr(lines)

	return
}
func (s *PreviewService) linesToStr(lines []interface{}) (data string) {
	for index, line := range lines {
		if index > 0 {
			data += "<br/>"
		}
		data += fmt.Sprint(line)
	}

	return
}

func NewPreviewService(defRepo *serverRepo.DefRepo, fieldRepo *serverRepo.FieldRepo, referRepo *serverRepo.ReferRepo, instancesRepo *serverRepo.InstancesRepo) *PreviewService {
	return &PreviewService{DefRepo: defRepo, FieldRepo: fieldRepo, ReferRepo: referRepo, InstancesRepo: instancesRepo}
}
