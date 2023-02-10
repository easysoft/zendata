package gen

import (
	"errors"
	"fmt"
	"github.com/easysoft/zendata/internal/pkg/domain"
	genHelper "github.com/easysoft/zendata/internal/pkg/gen/helper"
	"github.com/easysoft/zendata/internal/pkg/helper"
	"regexp"
	"strconv"
	"strings"
	"time"

	consts "github.com/easysoft/zendata/internal/pkg/const"
	commonUtils "github.com/easysoft/zendata/pkg/utils/common"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	stringUtils "github.com/easysoft/zendata/pkg/utils/string"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
)

func DealwithFixRange(field *domain.DefField) {
	if isRangeFix(field.Prefix) {
		field.PrefixRange = CreateFieldFixValuesFromList(field.Prefix, field)
	} else {
		var tmp interface{}
		tmp = field.Prefix
		field.PrefixRange = &domain.Range{Values: []interface{}{tmp}}
	}

	if isRangeFix(field.Postfix) {
		field.PostfixRange = CreateFieldFixValuesFromList(field.Postfix, field)
	} else {
		var tmp interface{}
		tmp = field.Postfix
		field.PostfixRange = &domain.Range{Values: []interface{}{tmp}}
	}
}

func genValuesForChildFields(field *domain.DefField, withFix bool, total int) (values []string) {
	fieldNameToValuesMap := map[string][]string{}
	fieldNameToFieldMap := map[string]domain.DefField{}

	// 1. generate values for sub fields
	for _, child := range field.Fields {
		if child.From == "" {
			child.From = field.From
		}

		child.FileDir = field.FileDir
		childValues := GenerateForFieldRecursive(&child, withFix, total)
		fieldNameToValuesMap[child.Field] = childValues
		fieldNameToFieldMap[child.Field] = child
	}

	// 2. deal with expression
	arrOfArr := make([][]string, 0) // 2 dimension arr for child, [ [a,b,c], [1,2,3] ]
	for _, child := range field.Fields {
		childValues := fieldNameToValuesMap[child.Field]

		if child.Value != "" {
			childValues = genHelper.GenExpressionValues(child, fieldNameToValuesMap, fieldNameToFieldMap)
		}
		arrOfArr = append(arrOfArr, childValues)
	}

	// 3. get combined values for parent field
	isRecursive := vari.GlobalVars.Recursive
	if stringUtils.InArray(field.Mode, consts.Modes) { // set on field level
		isRecursive = field.Mode == consts.ModeRecursive || field.Mode == consts.ModeRecursiveShort
	}

	values = combineChildrenValues(arrOfArr, isRecursive, total)
	values = loopFieldValues(field, values, len(values), true)

	return
}

func GenValuesForMultiRes(field *domain.DefField, withFix bool, total int) (values []string) {
	unionValues := make([]string, 0) // 2 dimension arr for child, [ [a,b,c], [1,2,3] ]

	for _, child := range field.Froms {
		if child.From == "" {
			child.From = field.From
		}

		child.FileDir = field.FileDir
		childValues := GenerateForFieldRecursive(&child, withFix, total)
		unionValues = append(unionValues, childValues...)
	}

	count := len(unionValues)
	if count > vari.GlobalVars.Total {
		count = vari.GlobalVars.Total
	}
	values = loopFieldValues(field, unionValues, count, true)

	return
}

func GenValuesForSingleRes(field *domain.DefField, total int) (values []string) {
	if field.Use != "" { // refer to ranges or instance
		groupValues := vari.Res[getFromKey(field)]

		uses := strings.TrimSpace(field.Use) // like group{limit:repeat}
		use, numLimit, repeat := getNum(uses)
		if strings.Index(use, "all") == 0 {
			valuesForAdd := getRepeatValuesFromAll(groupValues, numLimit, repeat)
			values = append(values, valuesForAdd...)
		} else {
			infos := parseUse(uses)
			valuesForAdd := getRepeatValuesFromGroups(groupValues, infos)
			values = append(values, valuesForAdd...)
		}
	} else if field.Select != "" { // refer to excel
		groupValues := vari.Res[getFromKey(field)]
		resKey := field.Select

		// deal with the key
		if vari.GlobalVars.DefData.Type == consts.DefTypeArticle {
			resKey = resKey + "_" + field.Field
		}

		values = append(values, groupValues[resKey]...)
	}

	values = loopFieldValues(field, values, total, true)

	return
}

func GenValuesForConfig(field *domain.DefField, total int) (values []string) {
	groupValues := vari.Res[field.Config]
	values = append(values, groupValues["all"]...)

	values = loopFieldValues(field, values, total, true)

	return
}

func isRangeFix(fix string) bool {
	index := strings.Index(fix, "-")

	return index > 0 && index < len(fix)-1
}

func GetFieldValStr(field domain.DefField, val interface{}) string {
	str := "n/a"
	success := false

	format := strings.TrimSpace(field.Format)

	if field.Type == consts.FieldTypeTimestamp && field.Format != "" {
		str = time.Unix(val.(int64), 0).Format(field.Format)
		return str
	}

	switch val.(type) {
	case int64:
		if format != "" {
			str, success = helper.FormatStr(format, val.(int64), 0)
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
			str, success = helper.FormatStr(format, val.(float64), precision)
		}
		if !success {
			str = strconv.FormatFloat(val.(float64), 'f', precision, 64)
		}
	case byte:
		str = string(val.(byte))
		if format != "" {
			str, success = helper.FormatStr(format, str, 0)
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
				str, success = helper.FormatStr(format, valInt, 0)
			}
		} else {
			str, success = helper.FormatStr(format, str, 0)
		}
	default:
	}

	return str
}

func loopFieldValues(field *domain.DefField, oldValues []string, total int, withFix bool) (values []string) {
	fieldValue := domain.FieldWithValues{}

	for _, val := range oldValues {
		fieldValue.Values = append(fieldValue.Values, val)
	}

	computerLoop(field)
	indexOfRow := 0
	count := 0
	for {
		// 处理格式、前后缀、loop等
		str := loopFieldValWithFix(field, fieldValue, &indexOfRow, count, withFix)
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

func loopFieldValWithFix(field *domain.DefField, fieldValue domain.FieldWithValues,
	indexOfRow *int, count int, withFix bool) (loopStr string) {

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

	loopStr = addFix(loopStr, field, count, withFix)

	return
}

func addFix(str string, field *domain.DefField, count int, withFix bool) (ret string) {
	prefix := GetStrValueFromRange(field.PrefixRange, count)
	postfix := GetStrValueFromRange(field.PostfixRange, count)
	divider := field.Divider

	if field.Length > runewidth.StringWidth(str) {
		str = helper.AddPad(str, *field)
	}
	if withFix && !vari.GlobalVars.Trim {
		str = prefix + str + postfix
	}
	if vari.GlobalVars.OutputFormat == consts.FormatText && !vari.GlobalVars.Trim {
		str += divider
	}

	ret = str
	return
}

func GetStrValueFromRange(rang *domain.Range, index int) string {
	if len(rang.Values) == 0 {
		return ""
	}

	idx := index % len(rang.Values)
	x := rang.Values[idx]
	return convPrefixVal2Str(x, "")
}

func convPrefixVal2Str(val interface{}, format string) string {
	str := "n/a"
	success := false

	switch val.(type) {
	case int64:
		if format != "" {
			str, success = helper.FormatStr(format, val.(int64), 0)
		}
		if !success {
			str = strconv.FormatInt(val.(int64), 10)
		}
	case float64:
		precision := 0
		if format != "" {
			str, success = helper.FormatStr(format, val.(float64), precision)
		}
		if !success {
			str = strconv.FormatFloat(val.(float64), 'f', precision, 64)
		}
	case byte:
		str = string(val.(byte))
		if format != "" {
			str, success = helper.FormatStr(format, str, 0)
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
				str, success = helper.FormatStr(format, valInt, 0)
			}
		} else {
			str, success = helper.FormatStr(format, str, 0)
		}
	default:
	}

	return str
}

func GenerateFieldVal(field domain.DefField, fieldValue domain.FieldWithValues, index *int) (val string, err error) {
	// 叶节点
	if len(fieldValue.Values) == 0 {
		if genHelper.IsSelectExcelWithExpr(field) {
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

func computerLoop(field *domain.DefField) {
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

func populateRowsFromTwoDimArr(arrOfArr [][]string, isRecursive, isOnTopLevel bool, total int) (values [][]string) {
	count := total
	if !isOnTopLevel {
		if isRecursive {
			count = getRecordCountForRecursive(arrOfArr)
		} else {
			count = getRecordCountForParallel(arrOfArr)
		}
	}

	indexArr := make([]int, 0)
	if isRecursive {
		indexArr = getModArr(arrOfArr)
	}

	for i := 0; i < count; i++ {
		strArr := make([]string, 0)
		for j := 0; j < len(arrOfArr); j++ {
			child := arrOfArr[j]

			var index int
			if isRecursive {
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

func RandomValuesArr(values [][]string) (ret [][]string) {
	length := len(values)
	for i := 0; i < length; i++ {
		val := commonUtils.RandNum(length)
		ret = append(ret, values[val])
	}

	return
}
func RandomInterfaces(values []interface{}) (ret []interface{}) {
	length := len(values)
	for i := 0; i < length; i++ {
		val := commonUtils.RandNum(length)
		ret = append(ret, values[val])
	}

	return
}
func RandomValues(values []interface{}) (ret []interface{}) {
	length := len(values)

	for i := 0; i < length; i++ {
		num := commonUtils.RandNum(length * 10000)
		ret = append(ret, values[num%len(values)])
	}

	return
}
func RandomStrValues(values []string) (ret []string) {
	length := len(values)

	for i := 0; i < length; i++ {
		num := commonUtils.RandNum(length * 10000)
		ret = append(ret, values[num%len(values)])
	}

	return
}

func combineChildrenValues(arrOfArr [][]string, isRecursive bool, total int) (ret []string) {
	valueArr := populateRowsFromTwoDimArr(arrOfArr, isRecursive, false, total)

	for _, arr := range valueArr {
		ret = append(ret, strings.Join(arr, ""))
	}
	return
}

func getRecordCountForParallel(arrOfArr [][]string) int {
	// get max count of 2nd dim arr
	count := 1
	for _, arr := range arrOfArr {
		if count < len(arr) {
			count = len(arr)
		}
	}

	if count > vari.GlobalVars.Total {
		count = vari.GlobalVars.Total
	}

	return count
}

func getRecordCountForRecursive(arrOfArr [][]string) int {
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

func getNum(group string) (ret string, numLimit, repeat int) {
	regx := regexp.MustCompile(`\{([^:]*):?([^:]*)\}`)
	arr := regx.FindStringSubmatch(group)
	if len(arr) >= 2 {
		numLimit, _ = strconv.Atoi(arr[1])
	}
	if len(arr) >= 3 {
		repeat, _ = strconv.Atoi(arr[2])
	}

	ret = regx.ReplaceAllString(group, "")

	return
}

// pars Uses
type retsInfo struct {
	ret      string
	numLimit int
	repeat   int
}

func parseUse(groups string) (results []retsInfo) {
	rets := strings.Split(groups, ",")
	results = make([]retsInfo, len(rets))
	regx := regexp.MustCompile(`\{([^:]*):?([^:]*)\}`)
	for k, v := range rets {
		results[k].ret = regx.ReplaceAllString(v, "")
		arr := regx.FindStringSubmatch(v)
		if len(arr) >= 2 {
			results[k].numLimit, _ = strconv.Atoi(arr[1])
		}
		if len(arr) >= 3 {
			results[k].repeat, _ = strconv.Atoi(arr[2])
			if results[k].repeat == 0 {
				results[k].repeat = 1
			}
		}
	}
	return
}

func getRepeatValuesFromAll(groupValues map[string][]string, numLimit, repeat int) (ret []string) {
	if repeat == 0 {
		repeat = 1
	}

	count := 0
exit:
	for _, arr := range groupValues {
		for _, item := range arr {
			for i := 0; i < repeat; i++ {
				ret = append(ret, item)
				count++

				if numLimit > 0 && count >= numLimit {
					break exit
				}
			}
		}
	}

	return
}

func getRepeatValuesFromGroups(groupValues map[string][]string, info []retsInfo) (ret []string) {
	count := 0

exit:
	for _, v := range info {
		if v.repeat == 0 {
			v.repeat = 1
		}

		arr := groupValues[v.ret]
		if len(arr) == 0 {
			break exit
		}
		if v.numLimit != 0 { // privateB{n}
			for i := 0; (v.numLimit > 0 && i < v.numLimit) && i < len(arr) && i < vari.GlobalVars.Total; i++ {
				index := i / v.repeat
				ret = append(ret, arr[index])
				count++
			}
		} else { // privateA
			for i := 0; i < len(arr) && i < vari.GlobalVars.Total; i++ {
				index := i / v.repeat % len(arr)
				ret = append(ret, arr[index])
				count++
			}
		}

		if count >= vari.GlobalVars.Total {
			break exit
		}

	}
	return
}

func getFromKey(field *domain.DefField) string {
	return fmt.Sprintf("%s-%s-%s", field.From, field.Use, field.Select)
}
