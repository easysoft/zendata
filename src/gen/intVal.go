package gen

import (
	constant "github.com/easysoft/zendata/src/utils/const"
	"math/rand"
)

func GenerateIntItems(start int64, end int64, step interface{}, rand bool) []interface{} {
	if !rand {
		return GenerateIntItemsByStep(start, end, step.(int))
	} else{
		return GenerateIntItemsRand(start, end, step.(int))
	}
}

func GenerateIntItemsByStep(start int64, end int64, step int) []interface{} {
	arr := make([]interface{}, 0)

	count := 0
	for i := 0; i < constant.MaxNumb; {
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

func GenerateIntItemsRand(start int64, end int64, step int) []interface{} {
	arr := make([]interface{}, 0)

	genCount := (end - start) / int64(step) + 1
	if genCount > int64(constant.MaxNumb) {
		genCount = int64(constant.MaxNumb)
	}

	for i := int64(0); i < genCount; {
		val := start + int64(rand.Int63n(genCount))

		arr = append(arr, val)
		i++
	}

	return arr
}