package service

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/easysoft/zendata/internal/pkg/domain"
	"github.com/easysoft/zendata/internal/pkg/helper"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/mattn/go-runewidth"
	"regexp"
	"strconv"
	"strings"
)

type ExpressionService struct {
}

func (s *ExpressionService) GenExpressionValues(field domain.DefField, valuesMap map[string][]interface{},
	fieldMap map[string]domain.DefField) (ret []interface{}) {
	exp := field.Value

	reg := regexp.MustCompile(`\$([_,a-z,A-Z][_,a-z,A-Z,0-9]+)`)
	arr := reg.FindAllStringSubmatch(exp, -1)

	total := 1
	typeGrade := map[string]int{
		"int":    0,
		"float":  1,
		"string": 2,
	}
	expressionType := "int"
	if strings.Contains(exp, "'") {
		expressionType = "string"
	}
	for _, items := range arr { // computer total
		placeholder := items[0]
		fieldName := items[1]
		exp = strings.Replace(exp, placeholder, fieldName, 1)

		size := len(valuesMap[fieldName])
		if total < size {
			total = size
		}

		// judge type of expression
		referField := fieldMap[fieldName]
		tp := s.getValuesType(valuesMap[fieldName], referField.Prefix, referField.Postfix)
		if typeGrade[tp] > typeGrade[expressionType] {
			expressionType = tp
		}
	}

	for i := 0; i < total; i++ {
		params := make(map[string]interface{})

		for _, items := range arr {
			fieldName := items[1]
			referValues := valuesMap[fieldName]
			referField := fieldMap[fieldName]

			valStr := "N/A"
			var val interface{}
			if len(referValues) > 0 {
				valStr = referValues[i%len(referValues)].(string)
				valStr = strings.TrimLeft(valStr, referField.Prefix)
				valStr = strings.TrimRight(valStr, referField.Postfix)

				val = s.parseValue(valStr, expressionType)
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
			if expressionType == "int" {
				mask = "%.0f"
			} else if expressionType == "float" {
				mask = "%f"
			} else {
				mask = "%s"
			}

			str := fmt.Sprintf(mask, result)
			if field.Length > runewidth.StringWidth(str) {
				str = helper.AddPad(str, field)
			}
			str = field.Prefix + str + field.Postfix
			ret = append(ret, str)
		}
	}

	return
}

func (s *ExpressionService) ReplaceVariableValues(exp string, valuesMap map[string][]interface{}) (ret []string) {
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
			referField := vari.GlobalVars.TopFieldMap[fieldName]

			valStr := "N/A"
			if len(referValues) > 0 {
				valStr = referValues[i%len(referValues)].(string)
				valStr = strings.TrimLeft(valStr, referField.Prefix)
				valStr = strings.TrimRight(valStr, referField.Postfix)
			}

			item = strings.ReplaceAll(item, fieldSlot, valStr)
		}

		ret = append(ret, item)
	}

	return
}

func (s *ExpressionService) getValuesType(values []interface{}, prefix string, postfix string) (tp string) {
	tool := map[string]int{
		"int":    0,
		"float":  1,
		"string": 2,
	}
	tp = "int"
	for _, item := range values {
		valStr := strings.TrimLeft(item.(string), prefix)
		valStr = strings.TrimRight(valStr, postfix)
		_, t := s.getType(valStr)
		if tool[t] > tool[tp] {
			tp = t
		}
		if tp == "string" {
			break
		}
	}
	return
}

func (s *ExpressionService) getType(str string) (val interface{}, tp string) {
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
	tp = "string"
	return
}

func (s *ExpressionService) parseValue(str string, tp string) (val interface{}) {
	var err error
	if tp == "int" {
		val, err = strconv.ParseInt(str, 0, 64)
		if err != nil {
			val = 0
		}
	} else if tp == "float" {
		val, err = strconv.ParseFloat(str, 64)
		if err != nil {
			val = 0.0
		}
	} else {
		val = str
	}
	return
}
