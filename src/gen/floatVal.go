package gen

import (
	constant "github.com/easysoft/zendata/src/utils/const"
	"math/rand"
	"strconv"
	"strings"
)

func GenerateFloatItems(start float64, end float64, step interface{}, rand bool) []interface{} {
	if !rand{
		return GenerateFloatItemsByStep(start, end, step.(float64))
	} else {
		return GenerateFloatItemsRand(start, end, step.(float64))
	}
}

func GenerateFloatItemsByStep(start float64, end float64, step interface{}) []interface{} {
	arr := make([]interface{}, 0)

	for i := 0; i < constant.MaxNumb; {
		gap := float64(i) * step.(float64)
		val := start + gap
		if val > end {
			break
		}

		arr = append(arr, val)
		i++
	}

	return arr
}

func GenerateFloatItemsRand(start float64, end float64, step float64) []interface{} {
	arr := make([]interface{}, 0)

	genCount := (end - start) / step
	if genCount > float64(constant.MaxNumb) {
		genCount = float64(constant.MaxNumb)
	}
	for i := float64(0); i < genCount; {
		val := start + float64(rand.Int63n(int64(genCount))) * step

		arr = append(arr, val)
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