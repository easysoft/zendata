package gen

import (
	constant "github.com/easysoft/zendata/src/utils/const"
	"math/rand"
)

func GenerateByteItems(start byte, end byte, step interface{}, rand bool, limit int) []interface{} {
	if !rand {
		return GenerateByteItemsByStep(start, end, step.(int), limit)
	} else {
		return GenerateByteItemsRand(start, end, step.(int), limit)
	}
}

func GenerateByteItemsByStep(start byte, end byte, step int, limit int) []interface{} {
	arr := make([]interface{}, 0)

	count := constant.MaxNumb
	for i := 0; i < constant.MaxNumb; {
		val := start + byte(int(i) * step)
		if val > end || i > limit {
			break
		}

		arr = append(arr, val)
		count++
		i++
	}

	return arr
}

func GenerateByteItemsRand(start byte, end byte, step int, limit int) []interface{} {
	arr := make([]interface{}, 0)

	count := int(end - start) / step + 1
	if count > limit {
		count = limit
	}
	if count > constant.MaxNumb {
		count = constant.MaxNumb
	}

	for i := 0; i < count; {
		ran := rand.Intn(count)
		val := start + byte(ran)

		arr = append(arr, val)
		i++
	}

	return arr
}