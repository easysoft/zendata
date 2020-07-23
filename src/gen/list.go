package gen

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	"math/rand"
	"strconv"
	"strings"
)

func GenerateList(field *model.DefField) model.FieldValue {
	fieldValue := model.FieldValue{}
	GenerateListField(field, &fieldValue)

	return fieldValue
}

func GenerateListField(field *model.DefField, fieldValue *model.FieldValue) {
	fieldValue.Field = field.Field
	fieldValue.Precision = field.Precision

	if len(field.Fields) > 0 {
		for _, child := range field.Fields {
			childValue := model.FieldValue{}
			GenerateListField(&child, &childValue)
		}
	} else {
		GenerateFieldValues(field, fieldValue)
	}
}

func GenerateFieldValues(field *model.DefField, fieldValue *model.FieldValue) {
	if strings.Index(field.Range, ".txt") > -1 {
		GenerateFieldValuesFromText(field, fieldValue)
	} else {
		GenerateFieldValuesFromList(field, fieldValue)
	}
}

func GenerateFieldValuesFromList(field *model.DefField, fieldValue *model.FieldValue) {
	rang := field.Range
	rangeItems := ParseRange(rang) // 1

	index := 0
	for _, rangeItem := range rangeItems {
		if index >= constant.MaxNumb { break }
		if rangeItem == "" { continue }

		entry, stepStr, repeat := ParseRangeItem(rangeItem)
		typ, desc := ParseEntry(entry) // 2

		items := make([]interface{}, 0)
		if typ == "literal" {
			items = GenerateValuesFromLiteral(desc, stepStr, repeat)
		} else if typ == "interval" {
			items = GenerateValuesFromInterval(field, desc, stepStr, repeat)
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

		if int1 > int2 && step.(int) > 0 {
			step = -1 * step.(int)
		}
		return "int", step, 0, rand

	} else {
		float1, errFloat1 := strconv.ParseFloat(startStr, 64)
		float2, errFloat2 := strconv.ParseFloat(endStr, 64)
		var errFloat3 error
		if strings.ToLower(stepStr) != "r" {
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

			if float1 > float2 && step.(int) > 0 {
				step = -1 * step.(int)
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

			if strings.Compare(startStr,endStr) > 0 && step.(int) > 0 {
				step = -1 * step.(int)
			}
			return "char", step, 0, rand
		}
	}

	return "string", 1, 0, false // is string
}

func GenerateValuesFromLiteral(desc string, stepStr string, repeat int) []interface{} {
	items := make([]interface{}, 0)

	elemArr := strings.Split(desc, ",")
	step, _ := strconv.Atoi(stepStr)
	total := 0

	for i := 0; i < len(elemArr); {
		for round := 0; round < repeat; round++ {
			val := ""
			if stepStr == "r" {
				val = elemArr[rand.Intn(len(elemArr))]
			} else {
				val = elemArr[i]
			}

			items = append(items, val)
			i += step
			total++

			if total > constant.MaxNumb {
				break
			}
		}

		if total >= constant.MaxNumb {
			break
		}
	}

	return items
}

func GenerateValuesFromInterval(field *model.DefField, desc string, stepStr string, repeat int) []interface{} {
	elemArr := strings.Split(desc, "-")
	startStr := elemArr[0]
	endStr := startStr
	if len(elemArr) > 1 { endStr = elemArr[1] }

	items := make([]interface{}, 0)
	dataType, step, precision, rand := CheckRangeType(startStr, endStr, stepStr)

	if dataType == "int" {
		startInt, _ := strconv.ParseInt(startStr, 0, 64)
		endInt, _ := strconv.ParseInt(endStr, 0, 64)

		items = GenerateIntItems(startInt, endInt, step, rand, repeat)
	} else if dataType == "float" {
		startFloat, _ := strconv.ParseFloat(startStr, 64)
		endFloat, _ := strconv.ParseFloat(endStr, 64)
		field.Precision = precision

		items = GenerateFloatItems(startFloat, endFloat, step.(float64), rand, repeat)
	} else if dataType == "char" {
		items = GenerateByteItems(byte(startStr[0]), byte(endStr[0]), step, rand, repeat)
	} else if dataType == "string" {
		if repeat == 0 { repeat = 1 }
		for i := 0; i < repeat; i++ {
			items = append(items, desc)
		}
	}

	return items
}
