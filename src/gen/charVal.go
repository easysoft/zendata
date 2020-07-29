package gen

import (
	constant "github.com/easysoft/zendata/src/utils/const"
)

//func GenerateByteItems(start byte, end byte, step interface{}, rand bool, repeat int) []interface{} {
//	if !rand {
//		return GenerateByteItemsByStep(start, end, step.(int), repeat)
//	} else {
//		return GenerateByteItemsRand(start, end, step.(int), repeat)
//	}
//}

func GenerateByteItemsByStep(start byte, end byte, step int, repeat int) []interface{} {
	arr := make([]interface{}, 0)

	total := 0
	for i := 0; true; {
		val := start + byte(int(i)*step)
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

		if total > constant.MaxNumb {
			break
		}
		i++
	}
	return arr
}

//func GenerateByteItemsRand(start byte, end byte, step int, repeat int) []interface{} {
//	arr := make([]interface{}, 0)
//
//	countInRound := int(int(end) - int(start)) / step
//	total := 0
//	for i := 0; i < countInRound; {
//		rand := commonUtils.RandNum(countInRound)
//		if step < 0 {
//			rand = rand * -1
//		}
//
//		val := start + byte(rand)
//
//		for round := 0; round < repeat; round++ {
//			arr = append(arr, val)
//
//			total++
//			if total > constant.MaxNumb {
//				break
//			}
//		}
//
//		if total > constant.MaxNumb {
//			break
//		}
//		i++
//	}
//
//	return arr
//}