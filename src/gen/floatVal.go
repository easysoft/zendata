package gen

import (
	constant "github.com/easysoft/zendata/src/utils/const"
	"math/rand"
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
	for round := 0; round < repeat; round++ {
		for i := 0; true; {
			gap := float64(i) * step.(float64)
			val := start + gap
			if val > end {
				break
			}

			arr = append(arr, val)
			i++
			total++

			if total > constant.MaxNumb {
				break
			}
		}
		if total > constant.MaxNumb {
			break
		}
	}

	return arr
}

func GenerateFloatItemsRand(start float64, end float64, step float64, repeat int) []interface{} {
	arr := make([]interface{}, 0)

	countInRound := (end - start) / step

	total := 0
	for round := 0; round < repeat; round++ {
		for i := float64(0); i < countInRound; {
			val := start + float64(rand.Int63n(int64(countInRound)))*step

			arr = append(arr, val)
			i++
			total++

			if total > constant.MaxNumb {
				break
			}
		}

		if total > constant.MaxNumb {
			break
		}
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