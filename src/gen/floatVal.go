package gen

import (
	constant "github.com/easysoft/zendata/src/utils/const"
	"math/rand"
	"strconv"
	"strings"
)

func GenerateFloatItems(start float64, end float64, step interface{}, rand bool, limit int) []interface{} {
	if !rand{
		return GenerateFloatItemsByStep(start, end, step.(float64), limit)
	} else {
		return GenerateFloatItemsRand(start, end, step.(float64), limit)
	}
}

func GenerateFloatItemsByStep(start float64, end float64, step interface{}, limit int) []interface{} {
	arr := make([]interface{}, 0)

	for i := 0; i < constant.MaxNumb; {
		gap := float64(i) * step.(float64)
		val := start + gap
		if val > end || i > limit {
			break
		}

		arr = append(arr, val)
		i++
	}

	return arr
}

func GenerateFloatItemsRand(start float64, end float64, step float64, limit int) []interface{} {
	arr := make([]interface{}, 0)

	count := (end - start) / step
	if count > float64(limit) {
		count = float64(limit)
	}
	if count > float64(constant.MaxNumb) {
		count = float64(constant.MaxNumb)
	}

	for i := float64(0); i < count; {
		val := start + float64(rand.Int63n(int64(count))) * step

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