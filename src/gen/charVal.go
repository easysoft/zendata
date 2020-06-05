package gen

import (
	constant "github.com/easysoft/zendata/src/utils/const"
	"math/rand"
)

func GenerateByteItems(start byte, end byte, step interface{}, rand bool, repeat int) []interface{} {
	if !rand {
		return GenerateByteItemsByStep(start, end, step.(int), repeat)
	} else {
		return GenerateByteItemsRand(start, end, step.(int), repeat)
	}
}

func GenerateByteItemsByStep(start byte, end byte, step int, repeat int) []interface{} {
	arr := make([]interface{}, 0)

	total := 0
	for round := 0; round < repeat; round++ {
		for i := 0; true; {
			val := start + byte(int(i)*step)
			if val > end {
				break
			}

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

func GenerateByteItemsRand(start byte, end byte, step int, repeat int) []interface{} {
	arr := make([]interface{}, 0)

	countInRound := int(end - start) / step + 1

	total := 0
	for round := 0; round < repeat; round++ {
		for i := 0; i < countInRound; {
			ran := rand.Intn(countInRound)
			val := start + byte(ran)

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