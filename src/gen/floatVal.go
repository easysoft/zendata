package gen

import (
	"math/rand"
	"strconv"
	"strings"
)

func GenerateFloatItems(start float64, end float64, step interface{}, index int, total int, rand bool) []interface{} {
	if !rand {
		return GenerateFloatItemsByStep(start, end, step.(float64), index, total)
	} else {
		return GenerateFloatItemsRand(start, end, step.(float64), index, total)
	}
}

func GenerateFloatItemsByStep(start float64, end float64, step interface{}, index int, total int) []interface{} {
	arr := make([]interface{}, 0)

	count := index
	for i := 0; i < total - index; {
		if count >= total {
			break
		}

		gap := float64(i) * step.(float64)
		val := start + gap
		if val > end {
			break
		}

		arr = append(arr, val)
		count++
		i++
	}

	return arr
}

func GenerateFloatItemsRand(start float64, end float64, step float64, index int, total int) []interface{} {
	arr := make([]interface{}, 0)

	genCount := (end - start) / step
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