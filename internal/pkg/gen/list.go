package gen

import (
	"fmt"
	"github.com/easysoft/zendata/internal/pkg/gen/helper"
	valueGen "github.com/easysoft/zendata/internal/pkg/gen/value"
	"math"
	"strconv"
	"strings"

	constant "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	commonUtils "github.com/easysoft/zendata/pkg/utils/common"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"github.com/easysoft/zendata/pkg/utils/vari"
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
		for i := 0; i < vari.Total; i++ {
			fieldValue.Values = append(fieldValue.Values, "")

			if strings.Index(field.Format, "uuid") == -1 {
				break
			}
		}

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

		descStr, stepStr, count, countTag := ParseRangeSection(rangeSection) // parse 2
		if strings.ToLower(stepStr) == "r" {
			(*field).IsRand = true
		}

		typ, desc := ParseRangeSectionDesc(descStr) // parse 3

		items := make([]interface{}, 0)
		if typ == "literal" {
			items = CreateValuesFromLiteral(field, desc, stepStr, count, countTag)
		} else if typ == "interval" {
			items = CreateValuesFromInterval(field, desc, stepStr, count, countTag)
		} else if typ == "yaml" { // load from a yaml
			items = CreateValuesFromYaml(field, desc, stepStr, count, countTag)
			field.ReferToAnotherYaml = true
		}

		fieldValue.Values = append(fieldValue.Values, items...)
		index = index + len(items)
	}

	if len(fieldValue.Values) == 0 {
		fieldValue.Values = append(fieldValue.Values, "N/A")
	}
}

func CreateFieldFixValuesFromList(strRang string, field *model.DefField) (rang *model.Range) {
	rang = &model.Range{}

	if strRang == "" {
		return
	}

	rangeSections := ParseRangeProperty(strRang) // parse 1

	index := 0
	for _, rangeSection := range rangeSections {
		if index >= constant.MaxNumb {
			break
		}
		if rangeSection == "" {
			continue
		}

		descStr, stepStr, count, countTag := ParseRangeSection(rangeSection) // parse 2
		if strings.ToLower(stepStr) == "r" {
			rang.IsRand = true
		}

		typ, desc := ParseRangeSectionDesc(descStr) // parse 3

		items := make([]interface{}, 0)
		if typ == "literal" {
			items = CreateValuesFromLiteral(field, desc, stepStr, count, countTag)
		} else if typ == "interval" {
			items = CreateValuesFromInterval(field, desc, stepStr, count, countTag)
		} else if typ == "yaml" { // load from a yaml
			items = CreateValuesFromYaml(field, desc, stepStr, count, countTag)
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

func CheckRangeType(startStr string, endStr string, stepStr string) (dataType string, step interface{}, precision int,
	rand bool, count int) {
	step = 1

	stepStr = strings.ToLower(strings.TrimSpace(stepStr))

	if start, end, stepi, ok := checkRangeTypeIsInt(startStr, endStr, stepStr); ok { // is int
		if stepStr == "r" {
			rand = true
		}

		count = (int)(start-end) / int(stepi)
		if count < 0 {
			count = count * -1
		}

		dataType = "int"
		step = int(stepi)
		return
	} else if start, end, stepf, ok := checkRangeTypeIsFloat(startStr, endStr, stepStr); ok { // is float
		if stepStr == "r" {
			rand = true
		}

		step = stepf

		precision1, step1 := valueGen.GetPrecision(start, step)
		precision2, step2 := valueGen.GetPrecision(end, step)
		if precision1 < precision2 {
			precision = precision2
			step = step2
		} else {
			precision = precision1
			step = step1
		}

		if (start > end && stepf > 0) || (start < end && stepf < 0) {
			step = -1 * stepf
		}

		dataType = "float"
		count = int(math.Floor(math.Abs(start-end)/stepf + 0.5))
		return

	} else if len(startStr) == 1 && len(endStr) == 1 { // is char
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

		dataType = "char"

		startChar := startStr[0]
		endChar := endStr[0]

		gap := 0
		if endChar > startChar { // 负数转uint单独处理
			gap = int(endChar - startChar)
		} else {
			gap = int(startChar - endChar)
		}
		count = gap / step.(int)
		if count < 0 {
			count = count * -1
		}
		return
	}

	// else is string
	dataType = "string"
	step = 1
	return
}

func CreateValuesFromLiteral(field *model.DefField, desc string, stepStr string, repeat int, repeatTag string) (items []interface{}) {
	elemArr := ParseDesc(desc)
	step, _ := strconv.Atoi(stepStr)
	if step == 0 {
		step = 1
	}
	total := 0

	if field.Path != "" && stepStr == "r" {
		pth := field.Path
		key := helper.GetRandFieldSection(pth)

		items = append(items, Placeholder(key))
		mp := placeholderMapForRandValues("list", elemArr, "", "", "", "",
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
			total = appendValues(&items, val, repeat, total)

			if total >= constant.MaxNumb {
				break
			}
			i += step
		}
	} else if repeatTag == "!" {
		isRand := field.Path == "" && stepStr == "r"
		for i := 0; i < repeat; {
			total = appendArrItems(&items, elemArr, total, isRand)

			if total >= constant.MaxNumb {
				break
			}
			i += step
		}
	}

	if field.Path == "" && stepStr == "r" { // for ranges and instances, random
		items = randomInterfaces(items)
	}

	return
}

func CreateValuesFromInterval(field *model.DefField, desc, stepStr string, repeat int, repeatTag string) (items []interface{}) {
	elemArr := strings.Split(desc, "-")
	startStr := elemArr[0]
	endStr := startStr
	if len(elemArr) > 1 {
		endStr = elemArr[1]
	}

	dataType, step, precision, rand, _ := CheckRangeType(startStr, endStr, stepStr)

	// 1. random replacement
	if field.Path != "" && dataType != "string" && rand { // random. for res, field.Path == ""
		pth := field.Path + "->" + desc
		key := helper.GetRandFieldSection(pth)

		val := Placeholder(key)
		strItems := make([]string, 0)

		//for i := 0; i < repeat*count; i++ { // chang to add only one placeholder item
		items = append(items, val)
		strItems = append(strItems, val)
		//}

		mp := placeholderMapForRandValues(dataType, strItems, startStr, endStr, stepStr,
			strconv.Itoa(precision), field.Format, repeat, repeatTag)

		vari.RandFieldSectionShortKeysToPathMap[key] = pth
		vari.RandFieldSectionPathToValuesMap[key] = mp

		return
	}

	if dataType == "int" {
		startInt, _ := strconv.ParseInt(startStr, 0, 64)
		endInt, _ := strconv.ParseInt(endStr, 0, 64)

		items = valueGen.GenerateIntItems(startInt, endInt, step.(int), rand, repeat, repeatTag)

	} else if dataType == "float" {
		startFloat, _ := strconv.ParseFloat(startStr, 64)
		endFloat, _ := strconv.ParseFloat(endStr, 64)
		field.Precision = precision

		items = valueGen.GenerateFloatItems(startFloat, endFloat, step, rand, precision, repeat, repeatTag)

	} else if dataType == "char" {
		items = valueGen.GenerateByteItems(startStr[0], endStr[0], step.(int), rand, repeat, repeatTag)

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

func CreateValuesFromYaml(field *model.DefField, yamlFile, stepStr string, repeat int, repeatTag string) (items []interface{}) {
	// keep root def, since vari.ZdDef will be overwrite by refer yaml file
	rootDef := vari.Def
	configDir := vari.ConfigFileDir
	res := vari.Res

	configFile := fileUtils.ComputerReferFilePath(yamlFile, field)
	fieldsToExport := make([]string, 0) // set to empty to use all fields
	rows, colIsNumArr, _ := GenerateFromYaml([]string{configFile}, &fieldsToExport)
	if field.Rand {
		rows = randomValuesArr(rows)
	}

	items = PrintLines(rows, constant.FormatData, "", colIsNumArr, fieldsToExport)

	if repeat > 0 {
		if repeat > len(items) {
			repeat = len(items)
		}
		items = items[:repeat]
	}

	// rollback root def when finish to deal with refer yaml file
	vari.Def = rootDef
	vari.ConfigFileDir = configDir
	vari.Res = res

	return
}

func Placeholder(key int) string {
	return fmt.Sprintf("${%d}", key)
}

func placeholderMapForRandValues(tp string, list []string, start, end, step, precision, format string,
	repeat int, repeatTag string) map[string]interface{} {
	ret := map[string]interface{}{}

	ret["type"] = tp

	// for literal values
	ret["list"] = list

	// for interval values
	ret["start"] = start
	ret["end"] = end
	ret["step"] = step
	ret["precision"] = precision
	ret["format"] = format

	ret["repeat"] = repeat
	ret["repeatTag"] = repeatTag

	return ret
}

func appendValues(items *[]interface{}, val string, repeat int, total int) int {
	for round := 0; round < repeat; round++ {
		*items = append(*items, val)

		total++
		if total > constant.MaxNumb {
			break
		}
	}

	return total
}

func appendArrItems(items *[]interface{}, arr []string, total int, isRand bool) int {
	for i := 0; i < len(arr); i++ {
		idx := i
		if isRand {
			idx = commonUtils.RandNum(len(arr)) // should set random here too
		}

		*items = append(*items, arr[idx])

		total++
		if total > constant.MaxNumb {
			break
		}
	}

	return total
}

func checkRangeTypeIsInt(startStr string, endStr string, stepStr string) (
	start int64, end int64, step int64, ok bool) {
	step = 1

	stepStr = strings.ToLower(strings.TrimSpace(stepStr))

	start, errInt1 := strconv.ParseInt(startStr, 0, 64)
	end, errInt2 := strconv.ParseInt(endStr, 0, 64)
	var errInt3 error

	if stepStr != "" && stepStr != "r" {
		step, errInt3 = strconv.ParseInt(stepStr, 0, 64)
	}

	if errInt1 == nil && errInt2 == nil && errInt3 == nil { // is int
		if (start > end && step > 0) || (start < end && step < 0) {
			step = -1 * step
		}

		ok = true
		return

	}

	return
}

func checkRangeTypeIsFloat(startStr string, endStr string, stepStr string) (
	start float64, end float64, step float64, ok bool) {
	step = 1.0
	stepStr = strings.ToLower(strings.TrimSpace(stepStr))

	start, errFloat1 := strconv.ParseFloat(startStr, 64)
	end, errFloat2 := strconv.ParseFloat(endStr, 64)
	var errFloat3 error

	if stepStr != "" && stepStr != "r" {
		step, errFloat3 = strconv.ParseFloat(stepStr, 64)
	}

	if errFloat1 == nil && errFloat2 == nil && errFloat3 == nil { // is float
		if (start > end && step > 0) || (start < end && step < 0) {
			step = -1 * step
		}

		ok = true
		return
	}

	return
}
