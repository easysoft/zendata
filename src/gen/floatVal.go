package gen

import (
	"strconv"
	"strings"
)

func GenerateFloatItems(start float64, end float64, step interface{}, index int, total int) []interface{} {
	if step != nil {
		return GenerateFloatItemsByStep(start, end, int64(step.(int)), index, total)
	} else {
		return GenerateFloatItemsRand(start, end, index, total)
	}
}

func GenerateFloatItemsByStep(start float64, end float64, step interface{}, index int, total int) []interface{} {
	arr := make([]interface{}, 0)

	count := index
	for i := 0; i < total - index; {
		gap := float64(i) * float64(step.(int))
		val := start + gap

		if val > end {
			if count < total { // loop if it's last item and not enough
				i = 0
				continue
			} else {
				break
			}
		}

		arr = append(arr, val)
		count++
		i++
	}

	return arr
}

func GenerateFloatItemsRand(start float64, end float64, index int, total int) []interface{} {

	return nil
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