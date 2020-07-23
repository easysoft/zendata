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

func GenerateIntItemsByStep(start int64, end int64, step int, repeat int) []interface{} {
	arr := make([]interface{}, 0)

	total := 0
	for i := 0; true; {
		val := start + int64(i*step)
		if (val > end && step > 0) || (val < end && step < 0)  {
			break
		}

		for round := 0; round < repeat; round++ {
			arr = append(arr, val)

			total++
			if total > constant.MaxNumb {
				break
			}
		}

		if total >= constant.MaxNumb {
			break
		}
		i++
	}

	return arr
}

func GenerateIntItemsRand(start int64, end int64, step int, repeat int) []interface{} {
	arr := make([]interface{}, 0)

	countInRound := (end - start) / int64(step) + 1

	total := 0
	for round := 0; round < repeat; round++ {
		for i := int64(0); i < countInRound; {
			val := start + rand.Int63n(countInRound)

			arr = append(arr, val)
			i++
			total++

			if total > constant.MaxNumb {
				break
			}
		}

		if total > constant.MaxNumb {
			break
		}
	}

	return arr
}