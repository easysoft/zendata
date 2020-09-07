package gen

import (
	"errors"
	"fmt"
	"github.com/easysoft/zendata/src/model"
	commonUtils "github.com/easysoft/zendata/src/utils/common"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/fatih/color"
	"regexp"
	"strconv"
	"strings"
)

func GenerateForDefinition(defaultFile, configFile string, fieldsToExport *[]string,
		) (rows [][]string, colIsNumArr []bool, err error) {

	vari.DefaultDir = fileUtils.GetAbsDir(defaultFile)
	vari.ConfigDir = fileUtils.GetAbsDir(configFile)

	vari.Def = LoadConfigDef(defaultFile, configFile, fieldsToExport)
	if len(vari.Def.Fields) == 0 {
		err = errors.New("")
		return
	}
	vari.Res = LoadResDef(*fieldsToExport)

	topFieldNameToValuesMap := map[string][]string{}

	// 为每个field生成值列表
	for index, field := range vari.Def.Fields {
		if !stringUtils.FindInArr(field.Field, *fieldsToExport) {
			continue
		}

		if field.Use != "" && field.From == "" {
			field.From = vari.Def.From
		}
		values := GenerateForField(&field, true)

		vari.Def.Fields[index].Precision = field.Precision

		topFieldNameToValuesMap[field.Field] = values
		colIsNumArr = append(colIsNumArr, field.IsNumb)
	}

	// 处理数据
	arrOfArr := make([][]string, 0) // 2 dimension arr for child, [ [a,b,c], [1,2,3] ]
	for _, child := range vari.Def.Fields {
		if !stringUtils.FindInArr(child.Field, *fieldsToExport) {
			continue
		}

		childValues := topFieldNameToValuesMap[child.Field]
		arrOfArr = append(arrOfArr, childValues)
	}
	rows = putChildrenToArr(arrOfArr)

	return
}

func GenerateForField(field *model.DefField, withFix bool) (values []string) {
	if len(field.Fields) > 0 { // sub fields
		arrOfArr := make([][]string, 0) // 2 dimension arr for child, [ [a,b,c], [1,2,3] ]
		for _, child := range field.Fields {
			if child.From == "" {
				child.From = field.From
			}

			childValues := GenerateForField(&child, withFix)
			arrOfArr = append(arrOfArr, childValues)
		}

		count := vari.Total
		count = getRecordCount(arrOfArr)
		if count > vari.Total {
			count = vari.Total
		}
		values = combineChildrenValues(arrOfArr, count)
		values = loopFieldValues(field, values, count, true)

	} else if len(field.Froms) > 0 { // from muti items
		unionValues := make([]string, 0) // 2 dimension arr for child, [ [a,b,c], [1,2,3] ]
		for _, child := range field.Froms {
			if child.From == "" {
				child.From = field.From
			}

			childValues := GenerateForField(&child, withFix)
			unionValues = append(unionValues, childValues...)
		}

		count := len(unionValues)
		if count > vari.Total {
			count = vari.Total
		}
		values = loopFieldValues(field, unionValues, count, true)

	} else if field.From != "" { // refer to res

		if field.Use != "" { // refer to instance
			groupValues := vari.Res[field.From]
			groups := strings.Split(field.Use, ",")
			for _, group := range groups {
				regx := regexp.MustCompile(`\{(.*)\}`)
				arr := regx.FindStringSubmatch(group)
				group = regx.ReplaceAllString(group, "")
				num := 0
				if len(arr) == 2 {
					num, _ = strconv.Atoi(arr[1])
				}

				i := num
				if group == "all" {
					for _, arr := range groupValues { // add all
						valuesFromGroup := make([]string, 0)
						if num == 0 {
							valuesFromGroup = arr
						} else {
							valuesFromGroup = arr[:num]
						}

						values = append(values, valuesFromGroup...)

						i = i - len(valuesFromGroup)
						if i <= 0 { break }
					}
				} else {
					valuesFromGroup := make([]string, 0)
					if num == 0 {
						valuesFromGroup = groupValues[group]
					} else {
						valuesFromGroup = groupValues[group][:num]
					}

					values = append(values, valuesFromGroup...)

					i = i - len(valuesFromGroup)
					if i <= 0 { break }
				}
			}
		} else if field.Select != "" { // refer to excel
			groupValues := vari.Res[field.From]
			slct := field.Select
			values = append(values, groupValues[slct]...)
		}

		values = loopFieldValues(field, values, vari.Total, true)

	} else if field.Config != "" { // refer to config
		groupValues := vari.Res[field.Config]
		values = append(values, groupValues["all"]...)

		values = loopFieldValues(field, values, vari.Total, true)

	} else { // leaf field
		values = GenerateFieldValuesForDef(field)
	}

	if field.Rand {
		values = randomValues(values)
	}

	return values
}

func GenerateFieldValuesForDef(field *model.DefField) []string {
	values := make([]string, 0)

	fieldWithValues := CreateField(field)

	computerLoop(field)
	indexOfRow := 0
	count := 0
	for {
		// 处理格式、前后缀、loop等
		val := loopFieldValWithFix(field, fieldWithValues, &indexOfRow, true)
		values = append(values, val)

		count++
		isRandomAndLoopEnd := !(*field).IsReferYaml && (*field).IsRand && (*field).LoopIndex == (*field).LoopEnd
		// isNotRandomAndValOver := !(*field).IsRand && indexOfRow >= len(fieldWithValues.Values)
		if count >= vari.Total || isRandomAndLoopEnd {
			break
		}

		(*field).LoopIndex = (*field).LoopIndex + 1
		if (*field).LoopIndex > (*field).LoopEnd {
			(*field).LoopIndex = (*field).LoopStart
		}
	}

	return values
}

func GetFieldValStr(field model.DefField, val interface{}) string {
	str := "n/a"
	success := false

	switch val.(type) {
		case int64:
			if field.Format != "" {
				str, success = stringUtils.FormatStr(field.Format, val.(int64))
			}
			if !success {
				str = strconv.FormatInt(val.(int64), 10)
			}
		case float64:
			precision := 0
			if field.Precision > 0 {
				precision = field.Precision
			}
			if field.Format != "" {
				str, success = stringUtils.FormatStr(field.Format, val.(float64))
			}
			if !success {
				str = strconv.FormatFloat(val.(float64), 'f', precision, 64)
			}
		case byte:
			str = string(val.(byte))
			if field.Format != "" {
				str, success = stringUtils.FormatStr(field.Format, str)
			}
			if !success {
				str = string(val.(byte))
			}
		case string:
			str = val.(string)
			fmt.Sprintf(str)
		default:
	}

	return str
}

func loopFieldValues(field *model.DefField, oldValues []string, total int, withFix bool) (values []string) {
	fieldValue := model.FieldWithValues{}

	for _, val := range oldValues {
		fieldValue.Values = append(fieldValue.Values, val)
	}

	computerLoop(field)
	indexOfRow := 0
	count := 0
	for {
		// 处理格式、前后缀、loop等
		str := loopFieldValWithFix(field, fieldValue, &indexOfRow, withFix)
		values = append(values, str)

		count++
		isRandomAndLoopEnd := (*field).IsRand && (*field).LoopIndex == (*field).LoopEnd
		isNotRandomAndValOver := !(*field).IsRand && indexOfRow >= len(fieldValue.Values)
		if count >= total || isRandomAndLoopEnd || isNotRandomAndValOver {
			break
		}

		(*field).LoopIndex = (*field).LoopIndex + 1
		if (*field).LoopIndex > (*field).LoopEnd {
			(*field).LoopIndex = (*field).LoopStart
		}
	}

	return
}

func loopFieldValWithFix(field *model.DefField, fieldValue model.FieldWithValues,
		indexOfRow *int, withFix bool) (loopStr string) {
	prefix := field.Prefix
	postfix := field.Postfix

	for j := 0; j < (*field).LoopIndex; j++ {
		if loopStr != "" {
			loopStr = loopStr + field.Loopfix
		}

		str, err := GenerateFieldVal(*field, fieldValue, indexOfRow)
		if err != nil {
			str = "N/A"
		}
		loopStr = loopStr + str

		*indexOfRow++
	}

	if withFix && !vari.Trim {
		loopStr = prefix + loopStr + postfix
	}

	return
}

func GenerateFieldVal(field model.DefField, fieldValue model.FieldWithValues, index *int) (val string, err error) {
	// 叶节点
	if len(fieldValue.Values) == 0 {
		logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("fail_to_generate_field", field.Field), color.FgCyan)
		err = errors.New("")
		return
	}

	idx := *index % len(fieldValue.Values)
	str := fieldValue.Values[idx]
	val = GetFieldValStr(field, str)

	return
}

func computerLoop(field *model.DefField) {
	if (*field).LoopIndex != 0 {
		return
	}

	arr := strings.Split(field.Loop, "-")
	(*field).LoopStart, _ = strconv.Atoi(arr[0])
	if len(arr) > 1 {
		field.LoopEnd, _ = strconv.Atoi(arr[1])
	}

	if (*field).LoopStart == 0 {
		(*field).LoopStart = 1
	}
	if (*field).LoopEnd == 0 {
		(*field).LoopEnd = 1
	}

	(*field).LoopIndex = (*field).LoopStart
}

func putChildrenToArr(arrOfArr [][]string) (values [][]string) {
	indexArr := make([]int, 0)
	if vari.Recursive {
		indexArr = getModArr(arrOfArr)
	}

	for i := 0; i < vari.Total; i++ {
		strArr := make([]string, 0)
		for j := 0; j < len(arrOfArr); j++ {
			child := arrOfArr[j]

			var index int
			if vari.Recursive {
				mod := indexArr[j]
				index = i / mod % len(child)
			} else {
				index = i % len(child)
			}

			strArr = append(strArr, child[index])
		}

		values = append(values, strArr)
	}

	return
}

func randomValuesArr(values [][]string) (ret [][]string) {
	length := len(values)
	for i := 0; i < length; i++ {
		val := commonUtils.RandNum(length)
		ret = append(ret, values[val])
	}

	return
}
func randomInterfaces(values []interface{}) (ret []interface{}) {
	length := len(values)
	for i := 0; i < length; i++ {
		val := commonUtils.RandNum(length)
		ret = append(ret, values[val])
	}

	return
}
func randomValues(values []string) (ret []string) {
	length := len(values)
	for i := 0; i < length; i++ {
		val := commonUtils.RandNum(length)
		ret = append(ret, values[val])
	}

	return
}

func combineChildrenValues(arrOfArr [][]string, total int) (ret []string)  {
	valueArr := putChildrenToArr(arrOfArr)

	for _, arr := range valueArr {
		ret = append(ret, strings.Join(arr, ""))
	}
	return
}

func getRecordCount(arrOfArr [][]string) int {
	count := 1
	for i := 0; i < len(arrOfArr); i++ {
		arr := arrOfArr[i]
		count = len(arr) * count
	}
	return count
}

func getModArr(arrOfArr [][]string) []int {
	indexArr := make([]int, 0)
	for _, _ = range arrOfArr {
		indexArr = append(indexArr, 0)
	}

	for i := 0; i < len(arrOfArr); i++ {
		loop := 1
		for j := i + 1; j < len(arrOfArr); j++ {
			loop = loop * len(arrOfArr[j])
		}

		indexArr[i] = loop
	}

	return indexArr
}
