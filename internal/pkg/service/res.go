package service

import (
	"fmt"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/helper"
	"github.com/easysoft/zendata/internal/pkg/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	stringUtils "github.com/easysoft/zendata/pkg/utils/string"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"gopkg.in/yaml.v3"
	"os"
)

type ResService struct {
	FieldService        *FieldService        `inject:""`
	ResYamlService      *ResYamlService      `inject:""`
	ResExcelService     *ResExcelService     `inject:""`
	ResRangesService    *ResRangesService    `inject:""`
	ResInstancesService *ResInstancesService `inject:""`
}

func (s *ResService) LoadResDef(fieldsToExport []string) (res map[string]map[string][]interface{}) {
	vari.GlobalVars.ResData = map[string]map[string][]interface{}{}

	for index, field := range vari.GlobalVars.DefData.Fields {
		if !stringUtils.StrInArr(field.Field, fieldsToExport) {
			continue
		}

		if (field.Use != "" || field.Select != "") && field.From == "" {
			field.From = vari.GlobalVars.DefData.From
			vari.GlobalVars.DefData.Fields[index].From = vari.GlobalVars.DefData.From
		}
		s.loadResForFieldRecursive(&field)
	}
	return
}

func (s *ResService) loadResForFieldRecursive(field *model.DefField) {
	if len(field.Fields) > 0 { // child fields
		for _, child := range field.Fields {
			if child.Use != "" && child.From == "" {
				child.From = field.From
			}

			child.FileDir = field.FileDir
			s.loadResForFieldRecursive(&child)
		}

	} else if len(field.Froms) > 0 { // multiple from
		for _, child := range field.Froms {
			if child.Use != "" && child.From == "" {
				child.From = field.From
			}

			child.FileDir = field.FileDir
			s.loadResForFieldRecursive(&child)
		}

	} else if field.From != "" && field.Type != consts.FieldTypeArticle { // from a res
		var valueMap map[string][]interface{}
		resFile, resType, sheet := fileUtils.GetResProp(field.From, field.FileDir) // relate to current file
		valueMap, _ = s.GetResValueFromExcelOrYaml(resFile, resType, sheet, field)

		if vari.GlobalVars.ResData[s.GetFromKey(field)] == nil {
			vari.GlobalVars.ResData[s.GetFromKey(field)] = map[string][]interface{}{}
		}

		for key, val := range valueMap {
			resKey := key
			// avoid article key to be duplicate
			if vari.GlobalVars.DefData.Type == consts.DefTypeArticle {
				resKey = resKey + "_" + field.Field
			}
			vari.GlobalVars.ResData[s.GetFromKey(field)][resKey] = val
		}

	} else if field.Config != "" { // from a config
		resFile, resType, _ := fileUtils.GetResProp(field.Config, field.FileDir)
		values, _ := s.GetResValueFromExcelOrYaml(resFile, resType, "", field)
		vari.GlobalVars.ResData[field.Config] = values
	}
}

func (s *ResService) GetResValueFromExcelOrYaml(resFile, resType, sheet string, field *model.DefField) (map[string][]interface{}, string) {
	resName := ""
	groupedValues := map[string][]interface{}{}

	if resType == "yaml" {
		groupedValues = s.ResYamlService.GetResFromYaml(resFile)
	} else if resType == "excel" {
		groupedValues = s.ResExcelService.GetResFromExcel(resFile, sheet, field)
	}

	return groupedValues, resName
}

func (s *ResService) GetReferencedRangeOrInstant(inst model.DefField) (referencedRanges model.ResRanges, referencedInsts model.ResInstances) {
	resFile, _, _ := fileUtils.GetResProp(inst.From, inst.FileDir)

	yamlContent, err := os.ReadFile(resFile)
	yamlContent = helper.ReplaceSpecialChars(yamlContent)
	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_read_file", resFile))
		return
	}

	err1 := yaml.Unmarshal(yamlContent, &referencedRanges)
	if err1 != nil || referencedRanges.Ranges == nil || len(referencedRanges.Ranges) == 0 { // parse ranges failed
		err2 := yaml.Unmarshal(yamlContent, &referencedInsts)
		if err2 != nil || referencedInsts.Instances == nil || len(referencedInsts.Instances) == 0 { // parse instances failed
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_parse_file", resFile))
			return
		} else { // is instances
			referencedInsts.FileDir = fileUtils.GetAbsDir(resFile)
		}
	} else { // is ranges
		referencedRanges.FileDir = fileUtils.GetAbsDir(resFile)
	}

	return
}

func (s *ResService) GetFromKey(field *model.DefField) string {
	return fmt.Sprintf("%s-%s-%s", field.From, field.Use, field.Select)
}
