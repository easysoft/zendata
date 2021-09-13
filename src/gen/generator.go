package gen

import (
	"errors"
	"fmt"
	"github.com/easysoft/zendata/src/gen/helper"
	"github.com/easysoft/zendata/src/model"
	commonUtils "github.com/easysoft/zendata/src/utils/common"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GenerateOnTopLevel(defaultFile, configFile string, fieldsToExport *[]string) (
	rows [][]string, colIsNumArr []bool, err error) {

	vari.DefaultFileDir = fileUtils.GetAbsDir(defaultFile)
	vari.ConfigFileDir = fileUtils.GetAbsDir(configFile)

	vari.Def = LoadDataDef(defaultFile, configFile, fieldsToExport)

	if len(vari.Def.Fields) == 0 {
		err = errors.New("")
		return
	} else if vari.Def.Type == constant.ConfigTypeArticle && vari.Out == "" {
		errMsg := i118Utils.I118Prt.Sprintf("gen_article_must_has_out_param")
		logUtils.PrintErrMsg(errMsg)
		err = errors.New(errMsg)
		return
	}

	if vari.Total < 0 {
		if vari.Def.Type == constant.ConfigTypeArticle {
			vari.Total = 1
		} else {
			vari.Total = constant.DefaultNumber
		}
	}

	vari.ResLoading = true // not to use placeholder when loading res
	vari.Res = LoadResDef(*fieldsToExport)
	vari.ResLoading = false

	topLevelFieldNameToValuesMap := map[string][]string{}

	// 为每个field生成值列表
	for index, field := range vari.Def.Fields {
		if !stringUtils.StrInArr(field.Field, *fieldsToExport) {
			continue
		}

		if field.Use != "" && field.From == "" {
			field.From = vari.Def.From
		}
		values := GenerateForField(&field, true)

		vari.Def.Fields[index].Precision = field.Precision

		topLevelFieldNameToValuesMap[field.Field] = values
		colIsNumArr = append(colIsNumArr, field.IsNumb)
	}

	// 处理数据
	arrOfArr := make([][]string, 0) // 2 dimension arr for child, [ [a,b,c], [1,2,3] ]
	for _, child := range vari.Def.Fields {
		if !stringUtils.StrInArr(child.Field, *fieldsToExport) {
			continue
		}

		childValues := topLevelFieldNameToValuesMap[child.Field]

		// is value expression
		if child.Value != "" {
			childValues = helper.GenExpressionValues(child, topLevelFieldNameToValuesMap, vari.TopFieldMap)
		}

		// select from excel with expr
		if helper.SelectExcelWithExpr(child) {
			selects := helper.ReplaceVariableValues(child.Select, topLevelFieldNameToValuesMap)
			wheres := helper.ReplaceVariableValues(child.Where, topLevelFieldNameToValuesMap)

			childValues = make([]string, 0)
			for index, slct := range selects {
				temp := child
				temp.Select = slct
				temp.Where = wheres[index%len(wheres)]

				resFile, _, sheet := fileUtils.GetResProp(temp.From, temp.FileDir)

				selectCount := vari.Total / len(selects)
				mp := generateFieldValuesFromExcel(resFile, sheet, &temp, selectCount) // re-generate values
				for _, items := range mp {
					childValues = append(childValues, items...)
				}
			}
		}

		arrOfArr = append(arrOfArr, childValues)
	}
	rows = putChildrenToArr(arrOfArr, vari.Recursive)

	return
}

func GenerateForField(field *model.DefField, withFix bool) (values []string) {
	if len(field.Fields) > 0 { // sub fields
		arrOfArr := make([][]string, 0) // 2 dimension arr for child, [ [a,b,c], [1,2,3] ]

		for _, child := range field.Fields {
			if child.From == "" {
				child.From = field.From
			}

			child.FileDir = field.FileDir
			childValues := GenerateForField(&child, withFix)
			arrOfArr = append(arrOfArr, childValues)
		}

		count := vari.Total
		count = getRecordCount(arrOfArr)
		if count > vari.Total {
			count = vari.Total
		}

		recursive := vari.Recursive
		if stringUtils.InArray(field.Mode, constant.Modes) { // set on field level
			recursive = field.Mode == constant.ModeRecursive || field.Mode == constant.ModeRecursiveShort
		}

		values = combineChildrenValues(arrOfArr, recursive)
		values = loopFieldValues(field, values, count, true)

	} else if len(field.Froms) > 0 { // from muti items
		unionValues := make([]string, 0) // 2 dimension arr for child, [ [a,b,c], [1,2,3] ]
		for _, child := range field.Froms {
			if child.From == "" {
				child.From = field.From
			}

			child.FileDir = field.FileDir
			childValues := GenerateForField(&child, withFix)
			unionValues = append(unionValues, childValues...)
		}

		count := len(unionValues)
		if count > vari.Total {
			count = vari.Total
		}
		values = loopFieldValues(field, unionValues, count, true)

	} else if field.From != "" && field.Type != constant.FieldTypeArticle { // refer to res
		if field.Use != "" { // refer to ranges or instance
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
						if i <= 0 {
							continue
						}
					}
				} else {
					valuesFromGroup := make([]string, 0)
					if num == 0 {
						valuesFromGroup = groupValues[group]
					} else {
						mode := num % len(groupValues[group])
						if mode == 0 {
							mode = len(groupValues[group])
						}
						valuesFromGroup = groupValues[group][:mode]
					}

					values = append(values, valuesFromGroup...)

					i = i - len(valuesFromGroup)
					if i <= 0 {
						continue
					}
				}
			}
		} else if field.Select != "" { // refer to excel
			groupValues := vari.Res[field.From]
			resKey := field.Select

			// deal with the key
			if vari.Def.Type == constant.ConfigTypeArticle {
				resKey = resKey + "_" + field.Field
			}

			values = append(values, groupValues[resKey]...)
		}

		values = loopFieldValues(field, values, vari.Total, true)

	} else if field.Config != "" { // refer to config
		groupValues := vari.Res[field.Config]
		values = append(values, groupValues["all"]...)

		values = loopFieldValues(field, values, vari.Total, true)

	} else { // leaf field
		values = GenerateFieldValuesForDef(field)
	}

	if field.Rand && field.Type != constant.FieldTypeArticle {
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
		// 2. random replacement
		isRandomAndLoopEnd := !vari.ResLoading && //  ignore rand in resource
			!(*field).ReferToAnotherYaml && (*field).IsRand && (*field).LoopIndex == (*field).LoopEnd
		// isNotRandomAndValOver := !(*field).IsRand && indexOfRow >= len(fieldWithValues.Values)
		if count >= vari.Total || count >= len(fieldWithValues.Values) || isRandomAndLoopEnd {
			for _, v := range fieldWithValues.Values {
				values = append(values, fmt.Sprintf("%v", v))
			}
			break
		}

		// 处理格式、前后缀、loop等
		val := loopFieldValWithFix(field, fieldWithValues, &indexOfRow, true)
		values = append(values, val)

		count++

		if count >= vari.Total || count >= len(fieldWithValues.Values) {
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

	format := strings.TrimSpace(field.Format)

	if field.Type == constant.FieldTypeTimestamp && field.Format != "" {
		str = time.Unix(val.(int64), 0).Format(field.Format)
		return str
	}

	switch val.(type) {
	case int64:
		if format != "" {
			str, success = stringUtils.FormatStr(format, val.(int64), 0)
		}
		if !success {
			str = strconv.FormatInt(val.(int64), 10)
		}
	case float64:
		precision := 0
		if field.Precision > 0 {
			precision = field.Precision
		}
		if format != "" {
			str, success = stringUtils.FormatStr(format, val.(float64), precision)
		}
		if !success {
			str = strconv.FormatFloat(val.(float64), 'f', precision, 64)
		}
	case byte:
		str = string(val.(byte))
		if format != "" {
			str, success = stringUtils.FormatStr(format, str, 0)
		}
		if !success {
			str = string(val.(byte))
		}
	case string:
		str = val.(string)

		match, _ := regexp.MatchString("%[0-9]*d", format)
		if match {
			valInt, err := strconv.Atoi(str)
			if err == nil {
				str, success = stringUtils.FormatStr(format, valInt, 0)
			}
		} else {
			str, success = stringUtils.FormatStr(format, str, 0)
		}
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
	divider := field.Divider

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

	if field.Length > runewidth.StringWidth(loopStr) {
		loopStr = stringUtils.AddPad(loopStr, *field)
	}
	if withFix && !vari.Trim {
		loopStr = prefix + loopStr + postfix
	}
	if vari.Format == constant.FormatText && !vari.Trim {
		loopStr += divider
	}

	return
}

func GenerateFieldVal(field model.DefField, fieldValue model.FieldWithValues, index *int) (val string, err error) {
	// 叶节点
	if len(fieldValue.Values) == 0 {
		if helper.SelectExcelWithExpr(field) {
			logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("fail_to_generate_field", field.Field), color.FgCyan)
			err = errors.New("")
		}
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

func putChildrenToArr(arrOfArr [][]string, recursive bool) (values [][]string) {
	indexArr := make([]int, 0)
	if recursive {
		indexArr = getModArr(arrOfArr)
	}

	for i := 0; i < vari.Total; i++ {
		strArr := make([]string, 0)
		for j := 0; j < len(arrOfArr); j++ {
			child := arrOfArr[j]

			var index int
			if recursive {
				mod := indexArr[j]
				index = i / mod % len(child)
			} else {
				index = i % len(child)
			}

			val := child[index]
			strArr = append(strArr, val)
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

func combineChildrenValues(arrOfArr [][]string, recursive bool) (ret []string) {
	valueArr := putChildrenToArr(arrOfArr, recursive)

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
