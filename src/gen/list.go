package gen

import (
	"github.com/easysoft/zendata/src/model"
	"strconv"
	"strings"
)

func GenerateList(field *model.Field, total int, fieldMap map[string][]interface{}) {
	name := strings.TrimSpace(field.Name)
	rang := strings.TrimSpace(field.Range)

	rangeItems := strings.Split(rang, ",")

	index := 0
	for itemIndex, item := range rangeItems {
		if index >= total { break }
		if strings.TrimSpace(item) == "" { continue }

		elemArr := strings.Split(item, "-")
		startStr := elemArr[0]
		endStr := startStr
		if len(elemArr) > 1 { endStr = elemArr[1] }

		items := make([]interface{}, 0)
		isLast := itemIndex == len(rangeItems) - 1

		dataType, step, precision := CheckRangeType(startStr, endStr, stepStr)

		if dataType == "int" { // int
			startInt, _ := strconv.ParseInt(startStr, 0, 64)
			endInt, _ := strconv.ParseInt(endStr, 0, 64)

			items = GenerateIntItems(startInt, endInt, int64(step.(int)), index, total, isLast)
		} else if dataType == "float" {
			startFloat, _ := strconv.ParseFloat(startStr, 64)
			endFloat, _ := strconv.ParseFloat(endStr, 64)

			field.Precision = precision
			items = GenerateFloatItems(startFloat, endFloat, step.(float64), index, total, isLast)
		} else if dataType == "char" {
			items = GenerateByteItems(byte(startStr[0]), byte(endStr[0]), step.(int), index, total, isLast)
		}

		fieldMap[name] = append(fieldMap[name], items...)
		index = index + len(items)
	}
}

func CheckRangeType(startStr string, endStr string, stepStr string) (string, interface{}, int) {
	_, errInt1 := strconv.ParseInt(startStr, 0, 64)
	_, errInt2 := strconv.ParseInt(endStr, 0, 64)
	if errInt1 == nil && errInt2 == nil {
		step, errInt3 := strconv.Atoi(stepStr)
		if errInt3 != nil {
			step = 1
		}

		return "int", step, 0
	} else {
		startFloat, errFloat1 := strconv.ParseFloat(startStr, 64)
		_, errFloat2 := strconv.ParseFloat(endStr, 64)
		if errFloat1 == nil && errFloat2 == nil {
			step, errFloat3 := strconv.ParseFloat(stepStr, 64)
			if errFloat3 != nil { step = 1.0 }
			precision := getPrecision(startFloat, step)

			return "float", step, precision
		} else if len(startStr) == 1 && len(endStr) == 1 {
			step, errChar3 := strconv.Atoi(stepStr)
			if errChar3 != nil {
				step = 1
			}

			return "char", step, 0
		}
	}
}