package gen

import (
	commonUtils "github.com/easysoft/zendata/src/utils/common"
	constant "github.com/easysoft/zendata/src/utils/const"
	"strconv"
	"strings"
)

func GenerateFloatItems(start float64, end float64, step interface{}, rand bool, repeat int) []interface{} {
	if !rand{
		return GenerateFloatItemsByStep(start, end, step.(float64), repeat)
	} else {
		return GenerateFloatItemsRand(start, end, step.(float64), repeat)
	}
}

func GenerateFloatItemsByStep(start float64, end float64, step interface{}, repeat int) []interface{} {
	arr := make([]interface{}, 0)

	total := 0

	for i := 0; true; {
		val := start + float64(i) * step.(float64)
		if (val > end && step.(float64) > 0) || (val < end && step.(float64) < 0)  {
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

	return arr
}

func GenerateFloatItemsRand(start float64, end float64, step float64, repeat int) []interface{} {
	arr := make([]interface{}, 0)

	countInRound := (end - start) / step
	total := 0
	for i := float64(0); i < countInRound; {
		rand := commonUtils.RandNum64(int64(countInRound))
		if step < 0 {
			rand = rand * -1
		}

		val := start + float64(rand) * step

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

	return arr
}

func getPrecision(base float64, step interface{}) int {
	val := base + step.(float64)

	str1 := strconv.FormatFloat(base, 'f', -1, 64)
	str2 := strconv.FormatFloat(val, 'f', -1, 64)

	index1 := strings.LastIndex(str1, ".")
	index2 := strings.LastIndex(str2, ".")

	if index1 < index2 {
		return len(str1) - index1 - 1
	} else {
		return len(str2) - index2 - 1
	}
}