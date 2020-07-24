package gen

import (
	"fmt"
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/mattn/go-runewidth"
	"strconv"
	"strings"
)

func GenerateForDefinition(defaultFile, configFile string, fieldsToExport *[]string, total int) ([][]string, []bool) {
	vari.DefaultDir = fileUtils.GetAbsDir(defaultFile)
	vari.ConfigDir = fileUtils.GetAbsDir(configFile)

	vari.Def = LoadConfigDef(defaultFile, configFile, fieldsToExport)
	vari.Res = LoadResDef(*fieldsToExport)

	fieldNameToValues := map[string][]string{}

	colTypes := make([]bool, 0)

	// 为每个field生成值列表
	for index, field := range vari.Def.Fields {
		if !stringUtils.FindInArr(field.Field, *fieldsToExport) {
			continue
		}

		values := GenerateForField(&field, total, true)
		vari.Def.Fields[index].Precision = field.Precision

		fieldNameToValues[field.Field] = values
		colTypes = append(colTypes, field.IsNumb)
	}

	// 生成指定数量行的数据
	rows := make([][]string, 0)
	for i := 0; i < total; i++ {
		for _, field := range vari.Def.Fields {
			if !stringUtils.FindInArr(field.Field, *fieldsToExport) {
				continue
			}

			values := fieldNameToValues[field.Field]
			fieldVal := values[i % len(values)]
			if len(rows) == i { rows = append(rows, make([]string, 0)) }
			rows[i] = append(rows[i], fieldVal)
		}
	}

	return rows, colTypes
}

func GenerateForField(field *model.DefField, total int, withFix bool) []string {
	values := make([]string, 0)

	if len(field.Fields) > 0 { // sub fields
		arr := make([][]string, 0) // 2 dimension arr for child, [ [a,b,c], [1,2,3] ]
		for _, child := range field.Fields {
			childValues := GenerateForField(&child, total, withFix)
			arr = append(arr, childValues)
		}

		for i := 0; i < total; i++ {
			concat := ""
			for _, row := range arr {
				concat = concat + row[i] // get one item from each child, grouped as a1 or b2
			}

			// should be done by calling LoopSubFields func as below, so disable this line
			//concat = field.Prefix + concat + field.Postfix
			values = append(values, concat)
		}
		values = LoopSubFields(field, values, total, true)

	} else if field.From != "" { // refer to res

		if field.Use != "" { // refer to instance
			if field.Use == "privateIP" {
				fmt.Println("")
			}

			groupValues := vari.Res[field.From]
			groups := strings.Split(field.Use, ",")
			for _, group := range groups {
				if group == "all" {
					for _, arr := range groupValues { // add all
						values = append(values, arr...)
					}
				} else {
					values = append(values, groupValues[group]...)
				}
			}
		} else if field.Select != "" { // refer to excel
			groupValues := vari.Res[field.From]
			slct := field.Select
			values = append(values, groupValues[slct]...)
		}

		values = LoopSubFields(field, values, total, true)

	} else if field.Config != "" { // refer to another define
		groupValues := vari.Res[field.Config]
		values = append(values, groupValues["all"]...)
	} else { // basic field
		values = GenerateFieldItemsFromDefinition(field)
	}

	return values
}

func GenerateFieldItemsFromDefinition(field *model.DefField) []string {
	if field.Loop == 0 {field.Loop = 1}

	values := make([]string, 0)

	fieldValue := GenerateList(field)

	index := 0
	count := 0
	for {
		// 处理格式、前后缀、loop等
		str := GenerateFieldValWithFix(*field, fieldValue, &index, true)
		values = append(values, str)

		count++
		if count >= constant.Total {
			break
		}
	}

	return values
}

func GenerateFieldVal(field model.DefField, fieldValue model.FieldValue, index *int) string {
	str := ""

	// 叶节点
	idx := *index % len(fieldValue.Values)
	val := fieldValue.Values[idx]
	str = GetFieldValStr(field, val)

	return str
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

func LoopSubFields(field *model.DefField, oldValues []string, total int, withFix bool) []string {
	if field.Loop == 0 {field.Loop = 1}

	values := make([]string, 0)
	fieldValue := model.FieldValue{}

	for _, val := range oldValues {
		fieldValue.Values = append(fieldValue.Values, val)
	}

	index := 0
	count := 0
	for {
		// 处理格式、前后缀、loop等
		str := GenerateFieldValWithFix(*field, fieldValue, &index, withFix)
		values = append(values, str)

		count++
		if count >= total {
			break
		}
	}

	return values
}

func GenerateFieldValWithFix(field model.DefField, fieldValue model.FieldValue, indexOfRow *int, withFix bool) string {
	prefix := field.Prefix
	postfix := field.Postfix

	loopStr := ""
	for j := 0; j < field.Loop; j++ {
		if loopStr != "" {
			loopStr = loopStr + field.Loopfix
		}

		str := GenerateFieldVal(field, fieldValue, indexOfRow)
		loopStr = loopStr + str

		*indexOfRow++
	}

	if withFix && !vari.Trim {
		loopStr = prefix + loopStr + postfix
	}

	if field.Width > runewidth.StringWidth(loopStr) {
		loopStr = stringUtils.AddPad(loopStr, field)
	}

	return loopStr
}