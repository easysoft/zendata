package gen

import (
	commonUtils "github.com/easysoft/zendata/src/utils/common"
	"strconv"
	"strings"
)

func GetRandFromList(list []string, repeatStr string, count int) []string {
	ret := make([]string, 0)

	for i := 0; i < count; i++ {
		rand := commonUtils.RandNum(len(list))
		val := list[rand]

		repeat, _ := strconv.Atoi(repeatStr)
		items := make([]string, 0)
		for round := 0; round < repeat; round++ {
			items = append(items, val)
		}

		ret = append(ret, strings.Join(items, ""))
	}

	return ret
}

func GetRandFromRange(dataType, start, end, step, repeatStr, precisionStr string, count int) []string {
	repeat, _ := strconv.Atoi(repeatStr)
	precision, _ := strconv.Atoi(precisionStr)

	ret := make([]string, 0)

	if dataType == "int" {
		startInt, _ := strconv.ParseInt(start, 0, 64)
		endInt, _ := strconv.ParseInt(end, 0, 64)
		stepInt, _ := strconv.ParseInt(step, 0, 64)

		if endInt < startInt && stepInt > 0 {
			stepInt = stepInt * -1
		}

		countInRound := (endInt - startInt) / stepInt

		for i := 0; i < count; i++ {
			rand := commonUtils.RandNum64(countInRound)
			if stepInt < 0 {
				rand = rand * -1
			}
			val := startInt + rand

			items := make([]string, 0)
			for round := 0; round < repeat; round++ {
				items = append(items, strconv.FormatInt(val, 10))
			}

			temp := strings.Join(items, "")
			ret = append(ret, temp)
		}
	} else if dataType == "char" {
		startChar := start[0]
		endChar := end[0]
		stepInt, _ := strconv.ParseInt(step, 10, 64)

		if int64(endChar) < int64(startChar) && stepInt > 0 {
			stepInt = stepInt * -1
		}

		countInRound := (int64(endChar) - int64(startChar)) / stepInt

		for i := 0; i < count; i++ {
			rand := commonUtils.RandNum64(countInRound)
			if stepInt < 0 {
				rand = rand * -1
			}
			val := startChar + byte(rand)
			items := make([]string, 0)
			for round := 0; round < repeat; round++ {
				items = append(items, string(val))
			}

			temp := strings.Join(items, "")
			ret = append(ret, temp)
		}

	} else if dataType == "float" {
		startFloat, _ := strconv.ParseFloat(start, 64)
		endFloat, _ := strconv.ParseFloat(end, 64)
		stepFloat, _ := strconv.ParseFloat(step, 64)

		if endFloat < startFloat && stepFloat > 0 {
			stepFloat = stepFloat * -1
		}

		countInRound := (endFloat - startFloat) / stepFloat

		for i := 0; i < count; i++ {
			rand := commonUtils.RandNum64(int64(countInRound))
			if stepFloat < 0 {
				rand = rand * -1
			}

			val := startFloat + float64(rand)*stepFloat

			items := make([]string, 0)
			for round := 0; round < repeat; round++ {
				str := strconv.FormatFloat(val, 'f', precision, 64)
				items = append(items, str)
			}

			temp := strings.Join(items, "")
			ret = append(ret, temp)
		}
	}

	return ret
}