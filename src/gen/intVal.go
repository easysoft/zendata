package gen

import (
	"math/rand"
)

func GenerateIntItems(start int64, end int64, step interface{}, index int, total int, rand bool) []interface{} {
	if !rand {
		return GenerateIntItemsByStep(start, end, step.(int), index, total)
	} else {
		return GenerateIntItemsRand(start, end, step.(int), index, total)
	}
}

func GenerateIntItemsByStep(start int64, end int64, step int, index int, total int) []interface{} {
	arr := make([]interface{}, 0)

	count := index
	for i := 0; i < total - index; {
		if count >= total {
			break
		}

		val := start + int64(i * step)
		if val > end {
			break
		}

		arr = append(arr, val)
		count++
		i++
	}

	return arr
}

func GenerateIntItemsRand(start int64, end int64, step int, index int, total int) []interface{} {
	arr := make([]interface{}, 0)

	genCount := (end - start) / int64(step) + 1
	if genCount > int64(total) - int64(index) {
		genCount = int64(total) - int64(index)
	}

	for i := int64(0); i < genCount; {
		val := start + int64(rand.Int63n(genCount))

		arr = append(arr, val)
		i++
	}

	return arr
}