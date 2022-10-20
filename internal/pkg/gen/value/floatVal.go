package valueGen

import (
	"fmt"
	constant "github.com/easysoft/zendata/internal/pkg/const"
	commonUtils "github.com/easysoft/zendata/pkg/utils/common"
	"strconv"
	"strings"
)

func GenerateFloatItems(start float64, end float64, step interface{}, rand bool, precision, repeat int, repeatTag string) []interface{} {
	if !rand {
		return generateFloatItemsByStep(start, end, step.(float64), precision, repeat, repeatTag)
	} else {
		return generateFloatItemsRand(start, end, step.(float64), precision, repeat, repeatTag)
	}
}

func generateFloatItemsByStep(start float64, end float64, step float64, precision, repeat int, repeatTag string) []interface{} {
	arr := make([]interface{}, 0)

	total := 0

	if repeatTag == "" {
		for i := 0; true; {
			val := start + float64(i)*step
			val = ChangePrecision(val, precision)

			if (val > end && step > 0) || (val < end && step < 0) {
				break
			}

			for round := 0; round < repeat; round++ {
				arr = append(arr, val)

				total++
				if total > constant.MaxNumb {
					break
				}
			}
			if total > constant.MaxNumb {
				break
			}
			i++
		}
	} else if repeatTag == "!" {
		for round := 0; round < repeat; round++ {
			for i := 0; true; {
				val := start + float64(i)*step
				val = ChangePrecision(val, precision)

				if (val > end && step > 0) || (val < end && step < 0) {
					break
				}

				arr = append(arr, val)

				if total > constant.MaxNumb {
					break
				}
				i++
			}

			total++
			if total > constant.MaxNumb {
				break
			}
		}
	}

	return arr
}

func generateFloatItemsRand(start float64, end float64, step float64, precision, repeat int, repeatTag string) []interface{} {
	arr := make([]interface{}, 0)

	countInRound := (end - start) / step
	total := 0

	if repeatTag == "" {
		for i := float64(0); i < countInRound; {
			rand := commonUtils.RandNum64(int64(countInRound))
			if step < 0 {
				rand = rand * -1
			}

			val := start + float64(rand)*step

			for round := 0; round < repeat; round++ {
				arr = append(arr, val)

				total++
				if total > constant.MaxNumb {
					break
				}
			}

			if total > constant.MaxNumb {
				break
			}
			i++
		}
	} else if repeatTag == "!" {
		for round := 0; round < repeat; round++ {
			for i := float64(0); i < countInRound; {
				rand := commonUtils.RandNum64(int64(countInRound))
				if step < 0 {
					rand = rand * -1
				}

				val := start + float64(rand)*step

				arr = append(arr, val)

				if total > constant.MaxNumb {
					break
				}
				i++
			}

			total++
			if total > constant.MaxNumb {
				break
			}
		}
	}

	return arr
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

func ChangePrecision(flt float64, precision int) float64 {
	format := fmt.Sprintf("%%.%df", precision)
	ret, _ := strconv.ParseFloat(fmt.Sprintf(format, flt), 64)
	return ret
}
