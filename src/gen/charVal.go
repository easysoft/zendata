package gen

import (
	constant "github.com/easysoft/zendata/src/utils/const"
	"math/rand"
)

func GenerateByteItems(start byte, end byte, step interface{}, rand bool) []interface{} {
	if !rand {
		return GenerateByteItemsByStep(start, end, step.(int))
	} else {
		return GenerateByteItemsRand(start, end, step.(int))
	}
}

func GenerateByteItemsByStep(start byte, end byte, step int) []interface{} {
	arr := make([]interface{}, 0)

	count := constant.MaxNumb
	for i := 0; i < constant.MaxNumb; {
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

func GenerateByteItemsRand(start byte, end byte, step int) []interface{} {
	arr := make([]interface{}, 0)

	genCount := int(end - start) / step + 1
	if genCount > constant.MaxNumb {
		genCount = constant.MaxNumb
	}

	for i := 0; i < genCount; {
		ran := rand.Intn(genCount)
		val := start + byte(ran)

		arr = append(arr, val)
		i++
	}

	return arr
}