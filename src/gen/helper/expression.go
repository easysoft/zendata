package helper

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/easysoft/zendata/src/model"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/mattn/go-runewidth"
	"regexp"
	"strconv"
	"strings"
)

func GenExpressionValues(field model.DefField, valuesMap map[string][]string) (ret []string) {
	exp := field.Value

	reg := regexp.MustCompile(`\$([_,a-z,A-Z,0-9]+)`)
	arr := reg.FindAllStringSubmatch(exp, -1)

	total := 1
	for _, items := range arr { // computer total
		placeholder := items[0]
		fieldName := items[1]
		exp = strings.Replace(exp, placeholder, fieldName, 1)

		size := len(valuesMap[fieldName])
		if total < size {
			total = size
		}
	}

	dataTypeFromValues := "int"
	for i := 0; i < total; i++ {
		params := make(map[string]interface{})

		//expr1, err1 := govaluate.NewEvaluableExpression("1+1")
		//result, err := expr1.Evaluate(params)
		//fmt.Sprintf("%v, %v, %v", expr1, err1, result)

		for _, items := range arr {
			fieldName := items[1]
			referValues := valuesMap[fieldName]
			referField := vari.TopFieldMap[fieldName]

			valStr := "N/A"
			tp := ""
			var val interface{}
			if len(referValues) > 0 {
				valStr = referValues[i%len(referValues)]
				valStr = strings.TrimLeft(valStr, referField.Prefix)
				valStr = strings.TrimRight(valStr, referField.Postfix)

				val, tp = getNumType(valStr)
				if tp != "int" {
					dataTypeFromValues = tp
				}
			}
			params[fieldName] = val
		}

		expr, err := govaluate.NewEvaluableExpression(exp)
		if err != nil {
			logUtils.PrintErrMsg(err.Error())
			ret = append(ret, "ERR")
		} else {
			result, err := expr.Evaluate(params)
			if err != nil {
				logUtils.PrintErrMsg(err.Error())
			}

			mask := ""
			if dataTypeFromValues == "int" {
				mask = "%.0f"
			} else {
				mask = "%d"
			}

			str := fmt.Sprintf(mask, result)
			if field.Length > runewidth.StringWidth(str) {
				str = stringUtils.AddPad(str, field)
			}
			str = field.Prefix + str + field.Postfix
			ret = append(ret, str)
		}
	}

	return
}

func ReplaceVariableValues(exp string, valuesMap map[string][]string) (ret []string) {
	reg := regexp.MustCompile(`\$\{([_,a-z,A-Z,0-9]+)\}`)
	arr := reg.FindAllStringSubmatch(exp, -1)

	total := 1
	for _, items := range arr { // computer total
		fieldName := items[1]

		size := len(valuesMap[fieldName])
		if total < size {
			total = size
		}
	}

	for i := 0; i < total; i++ {
		item := exp
		for _, items := range arr {
			fieldSlot := items[0]
			fieldName := items[1]
			referValues := valuesMap[fieldName]
			referField := vari.TopFieldMap[fieldName]

			valStr := "N/A"
			if len(referValues) > 0 {
				valStr = referValues[i%len(referValues)]
				valStr = strings.TrimLeft(valStr, referField.Prefix)
				valStr = strings.TrimRight(valStr, referField.Postfix)
			}

			item = strings.ReplaceAll(item, fieldSlot, valStr)
		}

		ret = append(ret, item)
	}

	return
}

func getNumType(str string) (val interface{}, tp string) {
	val, errInt := strconv.ParseInt(str, 0, 64)
	if errInt == nil {
		tp = "int"
		return
	}

	val, errFloat := strconv.ParseFloat(str, 64)
	if errFloat == nil {
		tp = "float"
		return
	}

	val = str
	tp = "float"

	return
}
