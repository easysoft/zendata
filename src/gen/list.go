package gen

import (
	"github.com/easysoft/zendata/src/model"
	commonUtils "github.com/easysoft/zendata/src/utils/common"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	"github.com/easysoft/zendata/src/utils/vari"
	"strconv"
	"strings"
)

func CreateListField(field *model.DefField, fieldWithValue *model.FieldWithValues) {
	fieldWithValue.Field = field.Field
	fieldWithValue.Precision = field.Precision

	if len(field.Fields) > 0 {
		for _, child := range field.Fields {
			childFieldWithValue := model.FieldWithValues{}
			CreateListField(&child, &childFieldWithValue)
		}
	} else {
		CreateListFieldValues(field, fieldWithValue)
	}
}

func CreateListFieldValues(field *model.DefField, fieldValue *model.FieldWithValues) {
	if strings.Index(field.Range, ".txt") > -1 {
		CreateFieldValuesFromText(field, fieldValue)
	} else {
		CreateFieldValuesFromList(field, fieldValue)
	}
}

func CreateFieldValuesFromList(field *model.DefField, fieldValue *model.FieldWithValues) {
	rang := field.Range
	if rang == "" {
		fieldValue.Values = append(fieldValue.Values, "")
		return
	}

	rangeSections := ParseRangeProperty(rang) // parse 1

	index := 0
	for _, rangeSection := range rangeSections {
		if index >= constant.MaxNumb {
			break
		}
		if rangeSection == "" {
			continue
		}

		descStr, stepStr, repeat := ParseRangeSection(rangeSection) // parse 2
		if strings.ToLower(stepStr) == "r" {
			(*field).IsRand = true
		}

		typ, desc := ParseRangeSectionDesc(descStr) // parse 3

		items := make([]interface{}, 0)
		if typ == "literal" {
			items = CreateValuesFromLiteral(field, desc, stepStr, repeat)
		} else if typ == "interval" {
			items = CreateValuesFromInterval(field, desc, stepStr, repeat)
		} else if typ == "yaml" { // load from a yaml
			items = CreateValuesFromYaml(field, desc, stepStr, repeat)
			field.ReferToAnotherYaml = true
		}

		fieldValue.Values = append(fieldValue.Values, items...)
		index = index + len(items)
	}

	if len(fieldValue.Values) == 0 {
		fieldValue.Values = append(fieldValue.Values, "N/A")
	}
}

func CheckRangeType(startStr string, endStr string, stepStr string) (string, interface{}, int, bool) {
	rand := false

	stepStr = strings.ToLower(strings.TrimSpace(stepStr))

	int1, errInt1 := strconv.ParseInt(startStr, 0, 64)
	int2, errInt2 := strconv.ParseInt(endStr, 0, 64)
	var errInt3 error
	if stepStr != "" && stepStr != "r" {
		_, errInt3 = strconv.ParseInt(stepStr, 0, 64)
	}
	if errInt1 == nil && errInt2 == nil && errInt3 == nil { // is int
		var step interface{} = 1
		if stepStr != "" && stepStr != "r" {
			stepInt, errInt3 := strconv.Atoi(stepStr)
			if errInt3 == nil {
				step = stepInt
			}
		} else if stepStr == "r" {
			rand = true
		}

		if (int1 > int2 && step.(int) > 0) || (int1 < int2 && step.(int) < 0) {
			step = -1 * step.(int)
		}
		return "int", step, 0, rand

	} else {
		float1, errFloat1 := strconv.ParseFloat(startStr, 64)
		float2, errFloat2 := strconv.ParseFloat(endStr, 64)
		var errFloat3 error
		if stepStr != "" && stepStr != "r" {
			_, errFloat3 = strconv.ParseFloat(stepStr, 64)
		}
		if errFloat1 == nil && errFloat2 == nil && errFloat3 == nil { // is float
			var step interface{} = nil

			if stepStr != "" && stepStr != "r" {
				stepFloat, errFloat3 := strconv.ParseFloat(stepStr, 64)
				if errFloat3 == nil {
					step = stepFloat
				}
			} else if stepStr == "r" {
				rand = true
			}

			precision1, step1 := GetPrecision(float1, step)
			precision2, step2 := GetPrecision(float2, step)
			precision := 0
			if precision1 < precision2 {
				precision = precision2
				step = step2
			} else {
				precision = precision1
				step = step1
			}

			if (float1 > float2 && step.(float64) > 0) || (float1 < float2 && step.(float64) < 0) {
				step = -1 * step.(float64)
			}
			return "float", step, precision, rand

		} else if len(startStr) == 1 && len(endStr) == 1 { // is char
			var step interface{} = 1
			if stepStr != "r" {
				stepChar, errChar3 := strconv.Atoi(stepStr)
				if errChar3 == nil {
					step = stepChar
				}
			} else if stepStr == "r" {
				rand = true
			}

			if (strings.Compare(startStr, endStr) > 0 && step.(int) > 0) ||
				(strings.Compare(startStr, endStr) < 0 && step.(int) < 0) {
				step = -1 * step.(int)
			}
			return "char", step, 0, rand
		}
	}

	return "string", 1, 0, false // is string
}

func CreateValuesFromLiteral(field *model.DefField, desc string, stepStr string, repeat int) (items []interface{}) {
	elemArr := ParseDesc(desc)
	step, _ := strconv.Atoi(stepStr)
	if step == 0 {
		step = 1
	}
	total := 0

	if field.Path != "" && stepStr == "r" {
		items = append(items, Placeholder(field.Path))
		mp := placeholderMapForRandValues("list", elemArr, "", "", "", "", field.Format)

		vari.RandFieldNameToValuesMap[field.Path] = mp
		return
	}

	for i := 0; i < len(elemArr); {
		idx := i
		if field.Path == "" && stepStr == "r" {
			idx = commonUtils.RandNum(len(elemArr)) // should set random here too
		}

		val := elemArr[idx]

		for round := 0; round < repeat; round++ {
			items = append(items, val)

			total++
			if total > constant.MaxNumb {
				break
			}
		}

		if total >= constant.MaxNumb {
			break
		}
		i += step
	}

	if field.Path == "" && stepStr == "r" { // for ranges and instances, random
		items = randomInterfaces(items)
	}

	return
}

func CreateValuesFromInterval(field *model.DefField, desc, stepStr string, repeat int) (items []interface{}) {
	elemArr := strings.Split(desc, "-")
	startStr := elemArr[0]
	endStr := startStr
	if len(elemArr) > 1 {
		endStr = elemArr[1]
	}

	dataType, step, precision, rand := CheckRangeType(startStr, endStr, stepStr)

	if field.Path != "" && dataType != "string" && rand { // random. for res, field.Path == ""
		items = append(items, Placeholder(field.Path))

		mp := placeholderMapForRandValues(dataType, []string{}, startStr, endStr, stepStr, strconv.Itoa(precision), field.Format)
		vari.RandFieldNameToValuesMap[field.Path] = mp

		return
	}

	if dataType == "int" {
		startInt, _ := strconv.ParseInt(startStr, 0, 64)
		endInt, _ := strconv.ParseInt(endStr, 0, 64)

		items = GenerateIntItems(startInt, endInt, step.(int), rand, repeat)

	} else if dataType == "float" {
		startFloat, _ := strconv.ParseFloat(startStr, 64)
		endFloat, _ := strconv.ParseFloat(endStr, 64)
		field.Precision = precision

		items = GenerateFloatItems(startFloat, endFloat, step, rand, repeat)

	} else if dataType == "char" {
		items = GenerateByteItems(startStr[0], endStr[0], step.(int), rand, repeat)

	} else if dataType == "string" {
		if repeat == 0 {
			repeat = 1
		}
		for i := 0; i < repeat; i++ {
			items = append(items, desc)
		}
	}

	if field.Path == "" && stepStr == "r" { // for ranges and instances, random again
		items = randomInterfaces(items)
	}

	return
}

func CreateValuesFromYaml(field *model.DefField, yamlFile, stepStr string, repeat int) (items []interface{}) {
	// keep root def, since vari.ZdDef will be overwrite by refer yaml file
	rootDef := vari.Def
	defaultDir := vari.DefaultFileDir
	configDir := vari.ConfigFileDir
	res := vari.Res

	configFile := fileUtils.ComputerReferFilePath(yamlFile)
	fieldsToExport := make([]string, 0) // set to empty to use all fields
	rows, colIsNumArr, _ := GenerateOnTopLevel("", configFile, &fieldsToExport)
	if field.Rand {
		rows = randomValuesArr(rows)
	}

	items = Print(rows, constant.FormatData, "", colIsNumArr, fieldsToExport)

	if repeat > 0 {
		if repeat > len(items) {
			repeat = len(items)
		}
		items = items[:repeat]
	}

	// rollback root def when finish to deal with refer yaml file
	vari.Def = rootDef
	vari.DefaultFileDir = defaultDir
	vari.ConfigFileDir = configDir
	vari.Res = res

	return
}

func Placeholder(str string) string {
	return "${" + str + "}"
}

func placeholderMapForRandValues(tp string, list []string, start, end, step, precision, format string) map[string]interface{} {
	ret := map[string]interface{}{}

	ret["type"] = tp

	ret["list"] = list // for literal values

	ret["start"] = start // for interval values
	ret["end"] = end
	ret["step"] = step
	ret["precision"] = precision
	ret["format"] = format

	return ret
}
