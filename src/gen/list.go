package gen

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	"github.com/easysoft/zendata/src/utils/vari"
	"strconv"
	"strings"
)

func CreateList(field *model.DefField) model.FieldWithValues {
	fieldWithValue := model.FieldWithValues{}
	CreateListField(field, &fieldWithValue)

	return fieldWithValue
}

func CreateListField(field *model.DefField, fieldWithValue *model.FieldWithValues) {
	fieldWithValue.Field = field.Field
	fieldWithValue.Precision = field.Precision

	if len(field.Fields) > 0 {
		for _, child := range field.Fields {
			childFieldWithValue := model.FieldWithValues{}
			CreateListField(&child, &childFieldWithValue)
		}
	} else {
		CreateFieldValues(field, fieldWithValue)
	}
}

func CreateFieldValues(field *model.DefField, fieldValue *model.FieldWithValues) {
	if strings.Index(field.Range, ".txt") > -1 {
		CreateFieldValuesFromText(field, fieldValue)
	} else {
		CreateFieldValuesFromList(field, fieldValue)
	}
}

func CreateFieldValuesFromList(field *model.DefField, fieldValue *model.FieldWithValues) {
	rang := field.Range
	rangeSections := ParseRangeProperty(rang) // 1

	index := 0
	for _, rangeSection := range rangeSections {
		if index >= constant.MaxNumb { break }
		if rangeSection == "" { continue }

		descStr, stepStr, repeat := ParseRangeSection(rangeSection) // 2
		if strings.ToLower(stepStr) == "r" {
			(*field).IsRand = true
		}

		typ, desc := ParseRangeSectionDesc(descStr) // 3

		items := make([]interface{}, 0)
		if typ == "literal" {
			items = CreateValuesFromLiteral(field, desc, stepStr, repeat)
		} else if typ == "interval" {
			items = CreateValuesFromInterval(field, desc, stepStr, repeat)
		} else if typ == "yaml" {
			items = CreateValuesFromYaml(field, desc, stepStr, repeat)
			field.IsReferYaml = true
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

	int1, errInt1 := strconv.ParseInt(startStr, 0, 64)
	int2, errInt2 := strconv.ParseInt(endStr, 0, 64)
	var errInt3 error
	if strings.ToLower(stepStr) != "r" {
		_, errInt3 = strconv.ParseInt(stepStr, 0, 64)
	}
	if errInt1 == nil && errInt2 == nil && errInt3 == nil { // is int
		var step interface{} = 1
		if strings.ToLower(strings.TrimSpace(stepStr)) != "r" {
			stepInt, errInt3 := strconv.Atoi(stepStr)
			if errInt3 == nil {
				step = stepInt
			}
		} else {
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
		if strings.ToLower(stepStr) != "" && strings.ToLower(stepStr) != "r" {
			_, errFloat3 = strconv.ParseFloat(stepStr, 64)
		}
		if errFloat1 == nil && errFloat2 == nil && errFloat3 == nil { // is float
			var step interface{} = 0.1
			if strings.ToLower(strings.TrimSpace(stepStr)) != "r" {
				stepFloat, errFloat3 := strconv.ParseFloat(stepStr, 64)
				if errFloat3 == nil {
					step = stepFloat
				}
			} else {
				rand = true
			}

			precision := getPrecision(float1, step)

			if (float1 > float2 && step.(float64) > 0) || (float1 < float2 && step.(float64) < 0) {
				step = -1 * step.(float64)
			}
			return "float", step, precision, rand

		} else if len(startStr) == 1 && len(endStr) == 1 { // is char
			var step interface{} = 1
			if strings.ToLower(strings.TrimSpace(stepStr)) != "r" {
				stepChar, errChar3 := strconv.Atoi(stepStr)
				if errChar3 == nil {
					step = stepChar
				}
			} else {
				rand = true
			}

			if (strings.Compare(startStr,endStr) > 0 && step.(int) > 0) ||
					(strings.Compare(startStr,endStr) < 0 && step.(int) < 0) {
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
	total := 0

	if stepStr == "r" {
		items = append(items, Placeholder(field.Path))
		mp := placeholderMapForRandValues("list", elemArr, "", "", "", "")

		vari.RandFieldNameToValuesMap[field.Path] = mp
		return
	}

	for i := 0; i < len(elemArr); {
		val := elemArr[i]

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

	return
}

func CreateValuesFromInterval(field *model.DefField, desc, stepStr string, repeat int) (items []interface{}) {
	elemArr := strings.Split(desc, "-")
	startStr := elemArr[0]
	endStr := startStr
	if len(elemArr) > 1 { endStr = elemArr[1] }

	dataType, step, precision, rand := CheckRangeType(startStr, endStr, stepStr)

	if dataType != "string" && rand {
		items = append(items, Placeholder(field.Path))

		mp := placeholderMapForRandValues(dataType, []string{}, startStr, endStr, stepStr, strconv.Itoa(precision))
		vari.RandFieldNameToValuesMap[field.Path] = mp

		return
	}

	if dataType == "int" {
		startInt, _ := strconv.ParseInt(startStr, 0, 64)
		endInt, _ := strconv.ParseInt(endStr, 0, 64)

		items = GenerateIntItemsByStep(startInt, endInt, step.(int), repeat)

	} else if dataType == "float" {
		startFloat, _ := strconv.ParseFloat(startStr, 64)
		endFloat, _ := strconv.ParseFloat(endStr, 64)
		field.Precision = precision

		items = GenerateFloatItemsByStep(startFloat, endFloat, step.(float64), repeat)

	} else if dataType == "char" {
		items = GenerateByteItemsByStep(startStr[0], endStr[0], step.(int), repeat)

	} else if dataType == "string" {
		if repeat == 0 { repeat = 1 }
		for i := 0; i < repeat; i++ {
			items = append(items, desc)
		}
	}

	return
}

func CreateValuesFromYaml(field *model.DefField, yamlFile, stepStr string, repeat int) (items []interface{}) {
	// keep root def, since vari.Def will be overwrite by refer yaml file
	rootDef := vari.Def

	configFile := vari.ConfigDir + yamlFile
	fieldsToExport := make([]string, 0) // set to empty to use all fields
	rows, colIsNumArr, _ := GenerateForDefinition("", configFile, &fieldsToExport)
	items = Print(rows, constant.FormatData, "", colIsNumArr, fieldsToExport)

	if repeat > 0 {
		if repeat > len(items) - 1 {
			repeat = len(items) - 1
		}
		items = items[:repeat]
	}

	// rollback root def when finish to deal with refer yaml file
	vari.Def = rootDef

	return
}

func Placeholder(str string) string {
	return "${" + str + "}"
}

func placeholderMapForRandValues(tp string, list []string, start, end, step, precision string) map[string]interface{} {
	ret := map[string]interface{}{}

	ret["type"] = tp

	ret["list"] = list // for literal values

	ret["start"] = start // for interval values
	ret["end"] = end
	ret["step"] = step
	ret["precision"] = precision

	return ret
}
