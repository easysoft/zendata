package gen

import (
	"fmt"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"math/rand"
)

func GenerateByteItems(start byte, end byte, step interface{}, index int, total int, rand bool) []interface{} {
	if !rand {
		return GenerateByteItemsByStep(start, end, step.(int), index, total)
	} else {
		return GenerateByteItemsRand(start, end, step.(int), index, total)
	}
}

func GenerateByteItemsByStep(start byte, end byte, step int, index int, total int) []interface{} {
	arr := make([]interface{}, 0)

	count := index
	for i := 0; i < total - index; {
		if count >= total {
			break
		}

		val := start + byte(int(i) * step)
		if val > end {
			break
		}

		arr = append(arr, val)
		count++
		i++
	}

	return arr
}

func GenerateByteItemsRand(start byte, end byte, step int, index int, total int) []interface{} {
	arr := make([]interface{}, 0)

	genCount := int(end - start) / step + 1
	if genCount > total - index {
		genCount = total - index
	}

	for i := 0; i < genCount; {
		ran := rand.Intn(genCount)
		logUtils.Screen(fmt.Sprintf("=== %d", ran))
		val := start + byte(ran)

		arr = append(arr, val)
		i++
	}

	return arr
}