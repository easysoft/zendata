package valueGen

import (
	constant "github.com/easysoft/zendata/internal/pkg/const"
	commonUtils "github.com/easysoft/zendata/pkg/utils/common"
	"math"
	"reflect"
)

func GenerateItems[TV ValType, TS StepType](start, end TV, step TS, precision int, isRand bool, repeat int, repeatTag string) (
	arr []interface{}) {

	typ := getType(start)

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
	int64 | byte | float32
}
type StepType interface {
	int64 | float32
}

func GetValue[TV ValType, TS StepType](start TV, step TS, precision int, it, limit int64, isRand bool) (ret TV) {
	typ := getType(start)

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
	typ := getType(a)

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

func getType(value interface{}) string {
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
