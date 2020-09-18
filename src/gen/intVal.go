package gen

import (
	commonUtils "github.com/easysoft/zendata/src/utils/common"
	constant "github.com/easysoft/zendata/src/utils/const"
)

func GenerateIntItems(start int64, end int64, step interface{}, rand bool, repeat int) []interface{} {
	if !rand {
		return generateIntItemsByStep(start, end, step.(int), repeat)
	} else{
		return generateIntItemsRand(start, end, step.(int), repeat)
	}
}

func generateIntItemsByStep(start int64, end int64, step int, repeat int) []interface{} {
	arr := make([]interface{}, 0)

	total := 0
	for i := 0; true; {
		val := start + int64(i * step)
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

func generateIntItemsRand(start int64, end int64, step int, repeat int) []interface{} {
	arr := make([]interface{}, 0)

	countInRound := (end - start) / int64(step) + 1
	total := 0
	for i := int64(0); i < countInRound; {
		rand := commonUtils.RandNum64(countInRound)
		if step < 0 {
			rand = rand * -1
		}

		val := start + rand
		for round := 0; round < repeat; round++ {
			arr = append(arr, val)

			total++

			if total > constant.MaxNumb {
				break
			}
		}

		if total > constant.MaxNumb {
			break
		}
		i++
	}

	return arr
}