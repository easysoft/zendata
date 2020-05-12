package gen

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func GenerateList(field *model.Field, total int) model.FieldValue {
	fieldValue := model.FieldValue{}
	GenerateListField(field, &fieldValue, 0)

	return fieldValue
}

func GenerateListField(field *model.Field, fieldValue *model.FieldValue, level int) {
	fieldValue.Name = field.Name
	fieldValue.Precision = field.Precision
	fieldValue.Level = level

	if len(field.Fields) > 0 {
		GenerateFieldChildren(field, fieldValue, level)
	} else {
		GenerateFieldValues(field, fieldValue, level)
	}
}
func GenerateFieldChildren(field *model.Field, fieldValue *model.FieldValue, level int) {
	for _, child := range field.Fields {
		childValue := model.FieldValue{}
		GenerateListField(&child, &childValue, level + 1)

		fieldValue.Children = append(fieldValue.Children, childValue)
	}
}

func GenerateFieldValues(field *model.Field, fieldValue *model.FieldValue, level int) {
	if strings.Index(field.Range, ".txt") > -1 {
		GenerateFieldValuesFromText(field, fieldValue, level)
	} else {
		GenerateFieldValuesFromList(field, fieldValue, level)
	}
}

func GenerateFieldValuesFromList(field *model.Field, fieldValue *model.FieldValue, level int) {
	rang := strings.TrimSpace(field.Range)

	rangeItems := strings.Split(rang, ",")
	index := 0
	for _, item := range rangeItems {
		if index >= constant.MaxNumb { break }
		if strings.TrimSpace(item) == "" { continue }

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

func GenerateFieldValuesFromText(field *model.Field, fieldValue *model.FieldValue, level int) {
	// get file and step string
	rang := strings.TrimSpace(field.Range)
	sectionArr := strings.Split(rang, ":")
	file := sectionArr[0]
	stepStr := "1"
	if len(sectionArr) == 2 { stepStr = sectionArr[1] }

	// read from file
	list := make([]string, 0)
	relaPath := constant.ResDir + file
	content, err := ioutil.ReadFile(relaPath)
	if err != nil {
		logUtils.Screen("fail to read " + relaPath + ", try to use global config")

		relaPath = "def" + string(os.PathSeparator) + file
		content, err = ioutil.ReadFile(relaPath)
		if err != nil {
			logUtils.Screen("fail to read " + relaPath + ", will return")

			fieldValue.Values = append(fieldValue.Values, "N/A")
			return
		} else {
			logUtils.Screen("success to read " + relaPath)
		}
	} else {
		logUtils.Screen("success to read " + relaPath)
	}
	str := string(content)
	str = strings.Replace(str, "\\r\\n", "\\n", -1)
	list = strings.Split(str, "\n")

	// get step and rand
	rand := false
	step := 1
	if strings.ToLower(strings.TrimSpace(stepStr)) != "r" {
		stepInt, err := strconv.Atoi(stepStr)
		if err == nil {
			step = stepInt
		}
	} else {
		rand = true
	}

	// get index for data retrieve
	numbs := GenerateIntItems(0, (int64)(len(list) - 1), step, rand)
	// get data by index
	index := 0
	for _, numb := range numbs {
		item := list[numb.(int64)]

		if index >= constant.MaxNumb { break }
		if strings.TrimSpace(item) == "" { continue }

		fieldValue.Values = append(fieldValue.Values, item)
		index = index + 1
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