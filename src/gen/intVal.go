package gen

import (
	constant "github.com/easysoft/zendata/src/utils/const"
	"math/rand"
)

func GenerateIntItems(start int64, end int64, step interface{}, rand bool, limit int) []interface{} {
	if !rand {
		return GenerateIntItemsByStep(start, end, step.(int), limit)
	} else{
		return GenerateIntItemsRand(start, end, step.(int), limit)
	}
}

func GenerateIntItemsByStep(start int64, end int64, step int, limit int) []interface{} {
	arr := make([]interface{}, 0)

	for i := 0; i < constant.MaxNumb; {
		val := start + int64(i * step)
		if val > end || i > limit {
			break
		}

		arr = append(arr, val)
		i++
	}

	return arr
}

func GenerateIntItemsRand(start int64, end int64, step int, limit int) []interface{} {
	arr := make([]interface{}, 0)

	count := (end - start) / int64(step) + 1
	if count > int64(limit) {
		count = int64(limit)
	}
	if count > int64(constant.MaxNumb) {
		count = int64(constant.MaxNumb)
	}

	for i := int64(0); i < count; {
		val := start + rand.Int63n(count)

		arr = append(arr, val)
		i++
	}

	return arr
}