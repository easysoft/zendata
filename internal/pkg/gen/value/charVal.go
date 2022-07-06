package valueGen

import (
	constant "github.com/easysoft/zendata/internal/pkg/const"
	commonUtils "github.com/easysoft/zendata/pkg/utils/common"
)

func GenerateByteItems(start byte, end byte, step int, rand bool, repeat int, repeatTag string) []interface{} {
	if !rand {
		return generateByteItemsByStep(start, end, step, repeat, repeatTag)
	} else {
		return generateByteItemsRand(start, end, step, repeat, repeatTag)
	}
}

func generateByteItemsByStep(start byte, end byte, step int, repeat int, repeatTag string) []interface{} {
	arr := make([]interface{}, 0)

	total := 0
	if repeatTag == "" {
		for i := 0; true; {
			val := start + byte(int(i)*step)
			if (val > end && step > 0) || (val < end && step < 0) {
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
	} else if repeatTag == "!" {
		for round := 0; round < repeat; round++ {
			for i := 0; true; {
				val := start + byte(int(i)*step)
				if (val >= end && step > 0) || (val <= end && step < 0) {
					break
				}

				arr = append(arr, val)

				if total > constant.MaxNumb {
					break
				}
				i++
			}

			total++
			if total > constant.MaxNumb {
				break
			}
		}
	}

	return arr
}

func generateByteItemsRand(start byte, end byte, step int, repeat int, repeatTag string) []interface{} {
	arr := make([]interface{}, 0)

	countInRound := (int(end) - int(start)) / step
	total := 0

	if repeatTag == "" {
		for i := 0; i < countInRound; {
			rand := commonUtils.RandNum(countInRound)
			if step < 0 {
				rand = rand * -1
			}

			val := start + byte(rand)

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
	} else if repeatTag == "!" {
		for round := 0; round < repeat; round++ {
			for i := 0; i < countInRound; {
				rand := commonUtils.RandNum(countInRound)
				if step < 0 {
					rand = rand * -1
				}

				val := start + byte(rand)
				arr = append(arr, val)

				if total > constant.MaxNumb {
					break
				}
				i++
			}

			total++
			if total > constant.MaxNumb {
				break
			}
		}
	}

	return arr
}
