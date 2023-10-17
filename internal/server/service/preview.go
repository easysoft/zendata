package serverService

import (
	"fmt"

	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/domain"
	"github.com/easysoft/zendata/internal/pkg/model"
	"github.com/easysoft/zendata/internal/pkg/service"
	serverRepo "github.com/easysoft/zendata/internal/server/repo"
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

	MainService   *service.MainService   `inject:""`
	OutputService *service.OutputService `inject:""`
	RangeService  *service.RangeService  `inject:""`
}

func (s *PreviewService) PreviewDefData(defId uint) (data string) {
	def, _ := s.DefRepo.Get(defId)

	vari.GlobalVars.Total = 10

	//lines := action.Generate([]string{def.Path}, "", consts.FormatData, "")
	s.MainService.GenerateDataByFile([]string{def.Path})
	lines := s.OutputService.GenText(false)

	data = s.linesToStr(lines)

	return
}
func (s *PreviewService) PreviewFieldData(fieldId uint, fieldType string) (data string) {
	var field model.ZdField

	if fieldType == consts.ResTypeDef {
		field, _ = s.FieldRepo.Get(fieldId)
	} else if fieldType == consts.ResTypeInstances {
		instItem, _ := s.InstancesRepo.GetItem(fieldId)
		field.From = instItem.From
		copier.Copy(&field, instItem)
	}

	ref := model.ZdRefer{}
	if !field.IsRange {
		ref, _ = s.ReferRepo.GetByOwnerIdAndType(field.ID, fieldType)
	}

	fld := domain.DefField{}
	genFieldFromZdField(field, ref, &fld)

	def := domain.DefData{}
	def.Fields = append(def.Fields, fld)
	defContent, _ := yaml.Marshal(def)

	configFile := vari.WorkDir + "tmp" + consts.PthSep + ".temp.yaml"
	fileUtils.WriteFile(configFile, string(defContent))

	//lines := action.Generate([]string{configFile}, field.Field, consts.FormatData, "")
	s.MainService.GenerateDataByFile([]string{configFile})
	lines := s.OutputService.GenText(false)

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
