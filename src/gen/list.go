package gen

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	"strconv"
	"strings"
)

func GenerateList(field *model.DefField, total int) model.FieldValue {
	fieldValue := model.FieldValue{}
	GenerateListField(field, &fieldValue, 0)

	return fieldValue
}

func GenerateListField(field *model.DefField, fieldValue *model.FieldValue, level int) {
	fieldValue.Field = field.Field
	fieldValue.Precision = field.Precision

	if len(field.Fields) > 0 {
		for _, child := range field.Fields {
			childValue := model.FieldValue{}
			GenerateListField(&child, &childValue, level + 1)
		}
	} else {
		GenerateFieldValues(field, fieldValue, level)
	}
}

func GenerateFieldValues(field *model.DefField, fieldValue *model.FieldValue, level int) {
	if strings.Index(field.Range, ".txt") > -1 {
		GenerateFieldValuesFromText(field, fieldValue, level)
	} else if strings.Index(field.Range, ".xlsx") > -1 {
		GenerateFieldValuesFromExcel(field, fieldValue, level)
	} else {
		GenerateFieldValuesFromList(field, fieldValue, level)
	}
}

func GenerateFieldValuesFromList(field *model.DefField, fieldValue *model.FieldValue, level int) {
	//rang := strings.TrimSpace(field.Range)
	rang := field.Range

	rangeItems := strings.Split(rang, ",")
	index := 0
	for _, item := range rangeItems {
		if index >= constant.MaxNumb { break }
		if item == "" { continue }

		sectionArr := strings.Split(item, ":")
		if len(sectionArr) == 0 { continue }

		stepStr := "1"
		if len(sectionArr) == 2 { stepStr = sectionArr[1] }

		elemArr := strings.Split(sectionArr[0], "-")
		startStr := elemArr[0]
		endStr := startStr
		if len(elemArr) > 1 { endStr = elemArr[1] }

		items := make([]interface{}, 0)

		dataType, step, precision, rand := CheckRangeType(startStr, endStr, stepStr)

		if dataType == "int" {
			startInt, _ := strconv.ParseInt(startStr, 0, 64)
			endInt, _ := strconv.ParseInt(endStr, 0, 64)

			items = GenerateIntItems(startInt, endInt, step, rand)
		} else if dataType == "float" {
			startFloat, _ := strconv.ParseFloat(startStr, 64)
			endFloat, _ := strconv.ParseFloat(endStr, 64)
			field.Precision = precision

			items = GenerateFloatItems(startFloat, endFloat, step.(float64), rand)
		} else if dataType == "char" {
			items = GenerateByteItems(byte(startStr[0]), byte(endStr[0]), step, rand)
		} else if dataType == "string" {
			items = append(items, startStr)
			if startStr != endStr {
				items = append(items, endStr)
			}
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

	_, errInt1 := strconv.ParseInt(startStr, 0, 64)
	_, errInt2 := strconv.ParseInt(endStr, 0, 64)
	if errInt1 == nil && errInt2 == nil {
		var step interface{} = 1
		if strings.ToLower(strings.TrimSpace(stepStr)) != "r" {
			stepInt, errInt3 := strconv.Atoi(stepStr)
			if errInt3 != nil {
				step = stepInt
			}
		} else {
			rand = true
		}

		return "int", step, 0, rand

	} else {
		startFloat, errFloat1 := strconv.ParseFloat(startStr, 64)
		_, errFloat2 := strconv.ParseFloat(endStr, 64)
		if errFloat1 == nil && errFloat2 == nil {
			var step interface{} = 0.1
			if strings.ToLower(strings.TrimSpace(stepStr)) != "r" {
				stepFloat, errFloat3 := strconv.ParseFloat(stepStr, 64)
				if errFloat3 == nil {
					step = stepFloat
				}
			} else {
				rand = true
			}

			precision := getPrecision(startFloat, step)

			return "float", step, precision, rand

		} else if len(startStr) == 1 && len(endStr) == 1 {
			var step interface{} = 1
			if strings.ToLower(strings.TrimSpace(stepStr)) != "r" {
				stepChar, errChar3 := strconv.Atoi(stepStr)
				if errChar3 != nil {
					step = stepChar
				}
			} else {
				rand = true
			}

			return "char", step, 0, rand
		}
	}

	return "string", 1, 0, false
}