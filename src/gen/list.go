package gen

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	"github.com/easysoft/zendata/src/utils/vari"
	"strconv"
	"strings"
)

func GenerateList(field *model.DefField) model.FieldWithValues {
	fieldValue := model.FieldWithValues{}
	GenerateListField(field, &fieldValue)

	return fieldValue
}

func GenerateListField(field *model.DefField, fieldValue *model.FieldWithValues) {
	fieldValue.Field = field.Field
	fieldValue.Precision = field.Precision

	if len(field.Fields) > 0 {
		for _, child := range field.Fields {
			childValue := model.FieldWithValues{}
			GenerateListField(&child, &childValue)
		}
	} else {
		GenerateFieldValues(field, fieldValue)
	}
}

func GenerateFieldValues(field *model.DefField, fieldValue *model.FieldWithValues) {
	if strings.Index(field.Range, ".txt") > -1 {
		GenerateFieldValuesFromText(field, fieldValue)
	} else {
		GenerateFieldValuesFromList(field, fieldValue)
	}
}

func GenerateFieldValuesFromList(field *model.DefField, fieldValue *model.FieldWithValues) {
	rang := field.Range
	rangeItems := ParseRange(rang) // 1

	index := 0
	for _, rangeItem := range rangeItems {
		if index >= constant.MaxNumb { break }
		if rangeItem == "" { continue }

		entry, stepStr, repeat := ParseRangeItem(rangeItem)
		typ, desc := ParseEntry(entry) // 2

		items := make([]interface{}, 0)
		itemsWithPlaceholder := make([]string, 0)
		if typ == "literal" {
			items = GenerateValuesFromLiteral(field, desc, stepStr, repeat)
		} else if typ == "interval" {
			items = GenerateValuesFromInterval(field, desc, stepStr, repeat)
		}

		fieldValue.Values = append(fieldValue.Values, items...)
		fieldValue.ValuesWithPlaceholder = append(fieldValue.ValuesWithPlaceholder, itemsWithPlaceholder...)
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

func GenerateValuesFromLiteral(field *model.DefField, desc string, stepStr string, repeat int) (items []interface{}) {
	elemArr := strings.Split(desc, ",")
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

func GenerateValuesFromInterval(field *model.DefField, desc string, stepStr string, repeat int) (items []interface{}) {
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
