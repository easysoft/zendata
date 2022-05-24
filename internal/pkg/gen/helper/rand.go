package helper

import (
	commonUtils "github.com/easysoft/zendata/pkg/utils/common"
	stringUtils "github.com/easysoft/zendata/pkg/utils/string"
	"strconv"
	"strings"
)

func GetRandFromList(list []string, repeat, count int) []string {
	ret := make([]string, 0)

	for i := 0; i < count; i++ {
		rand := commonUtils.RandNum(len(list))
		val := list[rand]

		items := make([]string, 0)
		for round := 0; round < repeat; round++ {
			items = append(items, val)
		}

		ret = append(ret, strings.Join(items, ""))
	}

	return ret
}

func GetRandFromRange(dataType, start, end, step string, repeat int, repeatTag, precisionStr string, count int, format string) []string {
	precision, _ := strconv.Atoi(precisionStr)

	ret := make([]string, 0)

	if dataType == "int" {
		startInt, _ := strconv.ParseInt(start, 0, 64)
		endInt, _ := strconv.ParseInt(end, 0, 64)
		stepInt, _ := strconv.ParseInt(step, 0, 64)

		if endInt < startInt && stepInt > 0 {
			stepInt = stepInt * -1
		}

		countInRound := (endInt-startInt)/stepInt + 1 // stepInt should be 1

		for i := 0; i < count; i++ {
			rand := commonUtils.RandNum64(countInRound)
			if stepInt < 0 {
				rand = rand * -1
			}
			val := startInt + rand

			items := make([]string, 0)
			item := strconv.FormatInt(val, 10)
			if format != "" {
				formatVal, success := stringUtils.FormatStr(format, val, 0)
				if success {
					item = formatVal
				}
			}

			for round := 0; round < repeat; round++ {
				items = append(items, item)
			}

			ret = append(ret, items...)
		}
	} else if dataType == "char" {
		startChar := start[0]
		endChar := end[0]
		stepInt, _ := strconv.ParseInt(step, 10, 64)

		if int64(endChar) < int64(startChar) && stepInt > 0 {
			stepInt = stepInt * -1
		}

		countInRound := (int64(endChar)-int64(startChar))/stepInt + 1 // stepInt should be 1

		for i := 0; i < count; i++ {
			rand := commonUtils.RandNum64(countInRound)
			if stepInt < 0 {
				rand = rand * -1
			}
			val := startChar + byte(rand)
			items := make([]string, 0)

			item := string(val)
			if format != "" {
				formatVal, success := stringUtils.FormatStr(format, val, 0)
				if success {
					item = formatVal
				}
			}

			for round := 0; round < repeat; round++ {
				items = append(items, item)
			}

			ret = append(ret, items...)
		}

	} else if dataType == "float" {
		startFloat, _ := strconv.ParseFloat(start, 64)
		endFloat, _ := strconv.ParseFloat(end, 64)
		stepFloat, _ := strconv.ParseFloat(step, 64)

		if endFloat < startFloat && stepFloat > 0 {
			stepFloat = stepFloat * -1
		}

		countInRound := (endFloat-startFloat)/stepFloat + 1 // stepInt should be 1

		for i := 0; i < count; i++ {
			rand := commonUtils.RandNum64(int64(countInRound))
			if stepFloat < 0 {
				rand = rand * -1
			}

			val := startFloat + float64(rand)*stepFloat

			items := make([]string, 0)

			item := strconv.FormatFloat(val, 'f', precision, 64)
			if format != "" {
				formatVal, success := stringUtils.FormatStr(format, val, precision)
				if success {
					item = formatVal
				}
			}

			for round := 0; round < repeat; round++ {
				items = append(items, item)
			}

			ret = append(ret, items...)
		}
	}

	return ret
}
