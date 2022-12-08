package valueGen

import (
	"fmt"
	constant "github.com/easysoft/zendata/internal/pkg/const"
	commonUtils "github.com/easysoft/zendata/pkg/utils/common"
	"math"
	"reflect"
	"strconv"
	"strings"
)

func GenerateItems[TV ValType, TS StepType](start, end TV, step TS, precision int, isRand bool, repeat int, repeatTag string) (
	arr []interface{}) {

	typ := GetType(start)

	limit := getLimit(start, end, step, typ, isRand)

	if repeatTag == "" { // repeat one by one
		for i := int64(0); i < limit; i++ {
			val := GetValue(start, step, precision, i, limit, isRand)

			if IsFinish(val, end, step) {
				break
			}

			RepeatSameVal(val, repeat, &arr)
		}
	} else if repeatTag == "!" { // repeat the list
		for round := 0; round < repeat; round++ {
			for i := int64(0); i < limit; i++ {
				val := GetValue(start, step, precision, i, limit, isRand)

				if IsFinish(val, end, step) {
					break
				}

				arr = append(arr, val)
				if len(arr) > constant.MaxNumb {
					break
				}
			}
		}
	}

	return
}

type ValType interface {
	int64 | byte | float64
}
type StepType interface {
	int64 | float64
}

func GetValue[TV ValType, TS StepType](start TV, step TS, precision int, it, limit int64, isRand bool) (ret TV) {
	typ := GetType(start)

	var val interface{}

	if typ == "int" {
		if !isRand {
			val = int64(start) + it*int64(step)
		} else {
			rand := commonUtils.RandNum64(limit)
			if step < 0 {
				rand = rand * -1
			}

			val = int64(start) + rand
		}

	} else if typ == "char" {
		if !isRand {
			val = byte(start) + byte(int(it)*int(step))
		} else {
			rand := commonUtils.RandNum(int(limit))
			if step < 0 {
				rand = rand * -1
			}

			val = byte(start) + byte(rand)
		}
	} else if typ == "float" {
		if !isRand {
			valFloat := float64(start) + float64(it)*float64(step)
			val = ChangePrecision(valFloat, precision)
		} else {
			rand := commonUtils.RandNum64(limit)
			if step < 0 {
				rand = rand * -1
			}

			val = float64(start) + float64(rand)*float64(step)
		}
	}

	ret = val.(TV)

	return
}

func IsFinish[TV ValType, TS StepType](a interface{}, b TV, step TS) bool {
	typ := GetType(a)

	if typ == "int" {
		if (a.(int64) > int64(b) && step > 0) || (a.(int64) < int64(b) && step < 0) {
			return true
		}
	} else if typ == "char" {
		if (a.(byte) > byte(b) && step > 0) || (a.(byte) < byte(b) && step < 0) {
			return true
		}

	} else if typ == "float" {
		if (a.(float64) > float64(b) && step > 0) || (a.(float64) < float64(b) && step < 0) {
			return true
		}
	}

	return false
}

func getLimit[TV ValType, TS StepType](start TV, end TV, step TS, typ string, isRand bool) (limit int64) {
	limit = int64(constant.MaxNumb)

	if isRand {
		if typ == "int" || typ == "char" {
			limit = (int64(end) - int64(start)) / int64(step)
		} else if typ == "float" {
			limitFloat := (float64(end) - float64(start)) / float64(step)
			limit = int64(math.Floor(limitFloat))
		}

		if limit > int64(constant.MaxNumb) {
			limit = int64(constant.MaxNumb)
		}
	}

	return
}

func RepeatSameVal[TV ValType](val TV, repeat int, arr *[]interface{}) {
	for round := 0; round < repeat; round++ {
		*arr = append(*arr, val)

		if len(*arr) > constant.MaxNumb {
			break
		}
	}
}

func GetType(value interface{}) string {
	v := reflect.ValueOf(value)

	switch v.Kind() {
	case reflect.Int64:
		return "int"
	case reflect.Uint8:
		return "char"
	case reflect.Float64:
		return "float"
	default:
		return ""
	}

	return ""
}

func GetPrecision(base float64, step interface{}) (precision int, newStep float64) {
	baseStr := strconv.FormatFloat(base, 'f', -1, 64)

	var stepFloat float64 = 1

	switch step.(type) {
	case float64:
		stepFloat = step.(float64)
	case int:
		stepFloat = float64(step.(int))
	}
	stepStr := strconv.FormatFloat(stepFloat, 'f', -1, 64)

	baseIndex := strings.LastIndex(baseStr, ".")
	stepIndex := strings.LastIndex(stepStr, ".")

	if baseIndex < 0 {
		baseIndex = 0
	}
	if stepIndex < 0 {
		stepIndex = 0
	}

	baseWidth := len(baseStr) - baseIndex - 1
	stepWidth := len(stepStr) - stepIndex - 1

	if baseWidth > stepWidth {
		precision = baseWidth
	} else {
		precision = stepWidth
	}

	if step == nil || step.(float64) == 0 {
		newStep = float64(1)
		for i := 0; i < precision; i++ {
			newStep = newStep / 10
		}
	} else {
		switch step.(type) {
		case float64:
			newStep = step.(float64)
		case int:
			newStep = float64(step.(int))
		}
	}

	return
}

func ChangePrecision(flt float64, precision int) float64 {
	format := fmt.Sprintf("%%.%df", precision)
	ret, _ := strconv.ParseFloat(fmt.Sprintf(format, flt), 64)
	return ret
}

func InterfaceToStr(val interface{}) string {
	str := "n/a"

	switch val.(type) {
	case int64:
		str = strconv.FormatInt(val.(int64), 10)
	case float64:
		precision, _ := GetPrecision(val.(float64), 0)
		str = strconv.FormatFloat(val.(float64), 'f', precision, 64)
	case byte:
		str = string(val.(byte))
	case string:
		str = val.(string)
	default:
	}
	return str
}
