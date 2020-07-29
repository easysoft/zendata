package gen

import (
	commonUtils "github.com/easysoft/zendata/src/utils/common"
	"strconv"
	"strings"
)

func GetRandFromList(list []string, repeat int) string {
	rand := commonUtils.RandNum(len(list))
	val := list[rand]

	items := make([]string, 0)
	for round := 0; round < repeat; round++ {
		items = append(items, val)
	}

	return strings.Join(items, "")
}

func GetRandFromRange(dataType, start, end, step string, repeat, precision int) string {
	if dataType == "int" {
		startInt, _ := strconv.ParseInt(start, 0, 64)
		endInt, _ := strconv.ParseInt(end, 0, 64)
		stepInt, _ := strconv.ParseInt(step, 0, 64)

		countInRound := (startInt - endInt) / stepInt
		rand := commonUtils.RandNum64(countInRound)
		if stepInt < 0 {
			rand = rand * -1
		}
		val := startInt + rand

		items := make([]string, 0)
		for round := 0; round < repeat; round++ {
			items = append(items, string(val))
		}

		return strings.Join(items, "")

	} else if dataType == "char" {
		startChar := start[0]
		endChar := end[0]
		stepInt, _ := strconv.ParseInt(step, 0, 64)
		countInRound := int64(int64(endChar) - int64(startChar)) / stepInt

		rand := commonUtils.RandNum64(countInRound)
		if stepInt < 0 {
			rand = rand * -1
		}
		val := startChar + byte(rand)
		items := make([]string, 0)
		for round := 0; round < repeat; round++ {
			items = append(items, string(val))
		}
		return strings.Join(items, "")

	} else if dataType == "float" {
		startFloat, _ := strconv.ParseFloat(start, 64)
		endFloat, _ := strconv.ParseFloat(end, 64)
		stepFloat, _ := strconv.ParseFloat(step, 64)

		countInRound := (startFloat - endFloat) / stepFloat

		rand := commonUtils.RandNum64(int64(countInRound))
		if stepFloat < 0 {
			rand = rand * -1
		}

		val := startFloat + float64(rand) * stepFloat

		items := make([]string, 0)
		for round := 0; round < repeat; round++ {
			str := strconv.FormatFloat(val, 'f', precision, 64)
			items = append(items, str)

			return strings.Join(items, "")

		}
	}

	return "N/A"
}