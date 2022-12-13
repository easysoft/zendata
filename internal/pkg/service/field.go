package service

import (
	"fmt"
	constant "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/gen"
	"github.com/easysoft/zendata/internal/pkg/gen/helper"
	valueGen "github.com/easysoft/zendata/internal/pkg/gen/value"
	"github.com/easysoft/zendata/internal/pkg/model"
	commonUtils "github.com/easysoft/zendata/pkg/utils/common"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"strconv"
	"strings"
)

type FieldService struct {
	TextService *TextService
}

func NewFieldService() *FieldService {
	return &FieldService{
		TextService: NewTextService(),
	}
}

func (s *FieldService) Generate(field *model.DefField) {

}

func (s *FieldService) CreateListField(field *model.DefField) {
	if len(field.Fields) > 0 {
		for _, child := range field.Fields {
			s.CreateListField(&child)
		}
	} else {
		s.CreateListFieldValues(field)
	}
}

func (s *FieldService) CreateListFieldValues(field *model.DefField) {
	if strings.Index(field.Range, ".txt") > -1 {
		s.TextService.CreateFieldValuesFromText(field)
	} else {
		s.CreateFieldValuesFromList(field)
	}
}

func (s *FieldService) CreateFieldValuesFromList(field *model.DefField) {
	rang := field.Range
	if rang == "" {
		for i := 0; i < vari.GenVars.Total; i++ {
			field.Values = append(field.Values, "")

			if strings.Index(field.Format, "uuid") == -1 {
				break
			}
		}

		return
	}

	rangeSections := gen.ParseRangeProperty(rang) // parse 1

	index := 0
	for _, rangeSection := range rangeSections {
		if index >= constant.MaxNumb {
			break
		}
		if rangeSection == "" {
			continue
		}

		descStr, stepStr, count, countTag := gen.ParseRangeSection(rangeSection) // parse 2
		if strings.ToLower(stepStr) == "r" {
			(*field).IsRand = true
		}

		typ, desc := gen.ParseRangeSectionDesc(descStr) // parse 3

		items := make([]interface{}, 0)
		if typ == "literal" {
			items = s.CreateValuesFromLiteral(field, desc, stepStr, count, countTag)
		} else if typ == "interval" {
			items = s.CreateValuesFromInterval(field, desc, stepStr, count, countTag)
		} else if typ == "yaml" { // load from a yaml
			items = s.CreateValuesFromYaml(field, desc, stepStr, count, countTag)
			field.ReferToAnotherYaml = true
		}

		field.Values = append(field.Values, items...)
		index = index + len(items)
	}

	if len(field.Values) == 0 {
		field.Values = append(field.Values, "N/A")
	}
}

func (s *FieldService) CreateFieldFixValuesFromList(strRang string, field *model.DefField) (rang *model.Range) {
	rang = &model.Range{}

	if strRang == "" {
		return
	}

	rangeSections := gen.ParseRangeProperty(strRang) // parse 1

	index := 0
	for _, rangeSection := range rangeSections {
		if index >= constant.MaxNumb {
			break
		}
		if rangeSection == "" {
			continue
		}

		descStr, stepStr, count, countTag := gen.ParseRangeSection(rangeSection) // parse 2
		if strings.ToLower(stepStr) == "r" {
			rang.IsRand = true
		}

		typ, desc := gen.ParseRangeSectionDesc(descStr) // parse 3

		items := make([]interface{}, 0)
		if typ == "literal" {
			items = s.CreateValuesFromLiteral(field, desc, stepStr, count, countTag)
		} else if typ == "interval" {
			items = s.CreateValuesFromInterval(field, desc, stepStr, count, countTag)
		} else if typ == "yaml" { // load from a yaml
			items = s.CreateValuesFromYaml(field, desc, stepStr, count, countTag)
			field.ReferToAnotherYaml = true
		}

		rang.Values = append(rang.Values, items...)
		index = index + len(items)
	}

	if len(rang.Values) == 0 {
		rang.Values = append(rang.Values, "N/A")
	}

	return
}

func (s *FieldService) CreateValuesFromLiteral(field *model.DefField, desc string, stepStr string, repeat int, repeatTag string) (items []interface{}) {
	elemArr := gen.ParseDesc(desc)
	step, _ := strconv.Atoi(stepStr)
	if step == 0 {
		step = 1
	}
	total := 0

	if field.Path != "" && stepStr == "r" {
		pth := field.Path
		key := helper.GetRandFieldSection(pth)

		items = append(items, gen.Placeholder(key))
		mp := gen.PlaceholderMapForRandValues("list", elemArr, "", "", "", "",
			field.Format, repeat, repeatTag)

		vari.RandFieldSectionShortKeysToPathMap[key] = pth
		vari.RandFieldSectionPathToValuesMap[key] = mp
		return
	}

	if repeatTag == "" {
		for i := 0; i < len(elemArr); {
			idx := i
			if field.Path == "" && stepStr == "r" {
				idx = commonUtils.RandNum(len(elemArr)) // should set random here too
			}

			val := elemArr[idx]
			total = gen.AppendValues(&items, val, repeat, total)

			if total >= constant.MaxNumb {
				break
			}
			i += step
		}
	} else if repeatTag == "!" {
		isRand := field.Path == "" && stepStr == "r"
		for i := 0; i < repeat; {
			total = gen.AppendArrItems(&items, elemArr, total, isRand)

			if total >= constant.MaxNumb {
				break
			}
			i += step
		}
	}

	if field.Path == "" && stepStr == "r" { // for ranges and instances, random
		items = gen.RandomInterfaces(items)
	}

	return
}

func (s *FieldService) CreateValuesFromInterval(field *model.DefField, desc, stepStr string, repeat int, repeatTag string) (items []interface{}) {
	elemArr := strings.Split(desc, "-")
	startStr := elemArr[0]
	endStr := startStr
	if len(elemArr) > 1 {
		endStr = elemArr[1]
	}

	dataType, step, precision, rand, _ := gen.CheckRangeType(startStr, endStr, stepStr)
	field.Precision = precision

	// 1. random replacement
	if field.Path != "" && dataType != "string" && rand { // random. for res, field.Path == ""
		pth := field.Path + "->" + desc
		key := helper.GetRandFieldSection(pth)

		val := gen.Placeholder(key)
		strItems := make([]string, 0)

		//for i := 0; i < repeat*count; i++ { // chang to add only one placeholder item
		items = append(items, val)
		strItems = append(strItems, val)
		//}

		mp := gen.PlaceholderMapForRandValues(dataType, strItems, startStr, endStr, fmt.Sprintf("%v", step),
			strconv.Itoa(precision), field.Format, repeat, repeatTag)

		vari.RandFieldSectionShortKeysToPathMap[key] = pth
		vari.RandFieldSectionPathToValuesMap[key] = mp

		return
	}

	if dataType == "int" {
		startInt, _ := strconv.ParseInt(startStr, 0, 64)
		endInt, _ := strconv.ParseInt(endStr, 0, 64)

		items = valueGen.GenerateItems(startInt, endInt, int64(step.(int)), 0, rand, repeat, repeatTag, 0)

	} else if dataType == "char" {
		items = valueGen.GenerateItems(startStr[0], endStr[0], int64(step.(int)), 0, rand, repeat, repeatTag, 0)

	} else if dataType == "float" {
		startFloat, _ := strconv.ParseFloat(startStr, 64)
		endFloat, _ := strconv.ParseFloat(endStr, 64)
		field.Precision = precision

		items = valueGen.GenerateItems(startFloat, endFloat, step.(float64), field.Precision, rand, repeat, repeatTag, 0)

	} else if dataType == "string" {
		if repeat == 0 {
			repeat = 1
		}
		for i := 0; i < repeat; i++ {
			items = append(items, desc)
		}
	}

	if field.Path == "" && stepStr == "r" { // for ranges and instances, random again
		items = gen.RandomInterfaces(items)
	}

	return
}

func (s *FieldService) CreateValuesFromYaml(field *model.DefField, yamlFile, stepStr string, repeat int, repeatTag string) (items []interface{}) {
	// keep root def, since vari.ZdDef will be overwrite by refer yaml file
	rootDef := vari.GenVars.DefData
	configDir := vari.GenVars.ConfigFileDir
	res := vari.Res

	configFile := fileUtils.ComputerReferFilePath(yamlFile, field)
	fieldsToExport := make([]string, 0) // set to empty to use all fields
	rows, colIsNumArr, _ := gen.GenerateFromYaml([]string{configFile}, &fieldsToExport)
	if field.Rand {
		rows = gen.RandomValuesArr(rows)
	}

	items = gen.PrintLines(rows, constant.FormatData, "", colIsNumArr, fieldsToExport)

	if repeat > 0 {
		if repeat > len(items) {
			repeat = len(items)
		}
		items = items[:repeat]
	}

	// rollback root def when finish to deal with refer yaml file
	vari.GenVars.DefData = rootDef
	vari.GenVars.ConfigFileDir = configDir
	vari.Res = res

	return
}
