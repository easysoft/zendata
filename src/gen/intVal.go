package gen

import (
	commonUtils "github.com/easysoft/zendata/src/utils/common"
	constant "github.com/easysoft/zendata/src/utils/const"
	"strconv"
	"strings"
)

func GenerateIntItems(start int64, end int64, step interface{}, rand bool, repeat int, repeatTag string) []interface{} {
	if !rand {
		return generateIntItemsByStep(start, end, step.(int), repeat, repeatTag)
	} else {
		return generateIntItemsRand(start, end, step.(int), repeat, repeatTag)
	}
}

func generateIntItemsByStep(start int64, end int64, step int, repeat int, repeatTag string) []interface{} {
	arr := make([]interface{}, 0)

	total := 0
	if repeatTag == "" {
		for i := 0; true; {
			val := start + int64(i*step)
			if (val > end && step > 0) || (val < end && step < 0) {
				break
			}

			for round := 0; round < repeat; round++ {
				arr = append(arr, val)

				total++
				if total > constant.MaxNumbForLangeRange {
					break
				}
			}

			if total >= constant.MaxNumbForLangeRange {
				break
			}
			i++
		}
	} else if repeatTag == "!" {
		for round := 0; round < repeat; round++ {
			for i := 0; true; {
				val := start + int64(i*step)
				if (val > end && step > 0) || (val < end && step < 0) {
					break
				}

				arr = append(arr, val)

				if total >= constant.MaxNumbForLangeRange {
					break
				}
				i++
			}

			if total >= constant.MaxNumbForLangeRange {
				break
			}
		}
	}

	return arr
}

func generateIntItemsRand(start int64, end int64, step int, repeat int, repeatTag string) []interface{} {
	arr := make([]interface{}, 0)

	countInRound := (end - start) / int64(step)
	total := 0

	str := strconv.FormatInt(end, 10)
	if strings.Count(str, "9") == len(str) && start == 0 {
		for i := int64(0); i < countInRound; {
			rand := commonUtils.RandStrAndNum(countInRound)
			for round := 0; round < repeat; round++ {
				arr = append(arr, rand)

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
	} else {
		if repeatTag == "" {
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
		} else if repeatTag == "!" {
			for round := 0; round < repeat; round++ {
				for i := int64(0); i < countInRound; {
					rand := commonUtils.RandNum64(countInRound)
					if step < 0 {
						rand = rand * -1
					}

					val := start + rand
					arr = append(arr, val)

					if total > constant.MaxNumb {
						break
					}
					i++
				}

				if total > constant.MaxNumb {
					break
				}
			}

		}
	}

	return arr
}
