package gen

import (
	"fmt"
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	"math/rand"
	"regexp"
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
	//rang := strings.TrimSpace(field.Range)
	rang := field.Range

	rangeItems := ParseRange(rang)

	index := 0
	for _, rangeItem := range rangeItems {
		if index >= constant.MaxNumb { break }
		if rangeItem == "" { continue }

		entry, stepStr, repeat := ParseRangeItem(rangeItem)
		typ, desc := ParseEntry(entry)

		items := make([]interface{}, 0)
		if typ == "literal" {
			items = GenerateValuesFromLiteral(field, desc, stepStr, repeat)
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

	_, errInt1 := strconv.ParseInt(startStr, 0, 64)
	_, errInt2 := strconv.ParseInt(endStr, 0, 64)
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

		return "int", step, 0, rand

	} else {
		startFloat, errFloat1 := strconv.ParseFloat(startStr, 64)
		_, errFloat2 := strconv.ParseFloat(endStr, 64)
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

			precision := getPrecision(startFloat, step)
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

			return "char", step, 0, rand
		}
	}

	return "string", 1, 0, false // is string
}

func ParseRange(rang string) []string {
	items := make([]string, 0)

	tagOpen := false
	temp := ""

	runeArr := make([]rune, 0)
	for _, c := range rang {
		runeArr = append(runeArr, c)
	}

	for i := 0; i < len(runeArr); i++ {
		c := runeArr[i]

		if c == constant.RightChar {
			tagOpen = false
		} else if c == constant.LeftChar  {
			tagOpen = true
		}

		if i == len(runeArr) - 1 {
			temp += fmt.Sprintf("%c", c)
			items = append(items, temp)
		} else if !tagOpen && c == ',' {
			items = append(items, temp)
			temp = ""
			tagOpen = false
		} else {
			temp += fmt.Sprintf("%c", c)
		}
	}

	return items
}

func ParseRangeItem(item string) (string, string, int) {
	entry := ""
	step := "1"
	repeat := 1

	item = strings.TrimSpace(item)
	if string(item[0]) == string(constant.LeftChar) &&  // It's a whole when meet (xxx)
			string(item[len(item) - 1]) == string(constant.RightChar) {
		return item, step, repeat
	}

	regx := regexp.MustCompile(`\{(.*)\}`)
	arr := regx.FindStringSubmatch(item)
	if len(arr) == 2 {
		repeat, _ = strconv.Atoi(arr[1])
	}
	item = regx.ReplaceAllString(item, "")

	sectionArr := strings.Split(item, ":")
	entry = sectionArr[0]
	if len(sectionArr) == 2 {
		step = sectionArr[1]
	}

	return entry, step, repeat
}

func GenerateValuesFromLiteral(field *model.DefField, desc string, stepStr string, repeat int) []interface{} {
	items := make([]interface{}, 0)

	elemArr := strings.Split(desc, ",")
	stepStr = strings.ToLower(strings.TrimSpace(stepStr))
	step, _ := strconv.Atoi(stepStr)
	total := 0
	for round := 0; round < repeat; round++ {
		for i := 0; i < len(elemArr); {
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

	// deal with exp like user-01


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
		//items = append(items, startStr)
		//if startStr != endStr {
		//	items = append(items, endStr)
		//}

		items = append(items, desc)
	}

	return items
}

func ParseEntry(str string) (string, string) {
	typ := ""
	desc := ""

	str = strings.TrimSpace(str)
	if int32(str[0]) == constant.LeftChar {
		typ = "literal"
		desc = strings.ReplaceAll(str, string(constant.LeftChar), "")
		desc = strings.ReplaceAll(desc,string(constant.RightChar), "")
	} else {
		typ = "interval"
		desc = str
	}

	return typ, desc
}
