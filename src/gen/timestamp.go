package gen

import (
	"fmt"
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"strconv"
	"strings"
	"time"
)

func CreateTimestampField(field *model.DefField, fieldWithValue *model.FieldWithValues) {
	convertTmFormat(field)

	fieldWithValue.Field = field.Field

	rang := strings.Trim(strings.TrimSpace(field.Range), ",")
	rangeSections := strings.Split(rang, ",")

	values := make([]interface{}, 0)
	for _, section := range rangeSections {
		createTimestampSectionValue(section, &values)
	}

	if len(values) == 0 {
		values = append(values, "N/A")
	}

	fieldWithValue.Values = values
}

func convertTmFormat(field *model.DefField) { // to 2006-01-02 15:04:05
	format := field.Format

	if strings.Index(format, "YYYY") > -1 {
		format = strings.Replace(format, "YYYY", "2006", 1)
	} else {
		format = strings.Replace(format, "YY", "06", 1)
	}

	format = strings.Replace(format, "MM", "01", 1)
	format = strings.Replace(format, "DD", "02", 1)
	format = strings.Replace(format, "hh", "15", 1)
	format = strings.Replace(format, "mm", "04", 1)
	format = strings.Replace(format, "ss", "05", 1)

	field.Format = format
}

func createTimestampSectionValue(section string, values *[]interface{}) {
	desc, step, unit := parseTsSection(section)
	start, end := parseTsDesc(desc)

	if (step > 0 && start > end) || (step < 0 && start < end) {
		step *= -1
	}

	// get index numbers for data retrieve
	numbs := GenerateTimeItems(start, end, step, unit, 1, "")

	// generate data by index
	index := 0
	for _, numb := range numbs {
		if index >= constant.MaxNumb {
			break
		}

		*values = append(*values, numb)
		index = index + 1
	}
}

func parseTsSection(section string) (desc string, step int, unit string) {
	section = strings.Trim(strings.TrimSpace(section), ":")

	sectionArr := strings.Split(section, ":")
	desc = sectionArr[0]
	if len(sectionArr) > 1 {
		stepStr := sectionArr[1]
		step, unit = parseTsExpr(stepStr)
		if step == 0 {
			step = 1
		}
	}

	return
}

func parseTsDesc(desc string) (start, end int64) {
	desc = strings.TrimSpace(desc)

	if strings.Contains(desc, "today") {
		start, end = getTodayTs()
		return
	}

	startStr, endStr := splitTmDesc(desc)

	start = parseTsValue(startStr, true)
	end = parseTsValue(endStr, false)

	//logUtils.PrintTo(
	//	fmt.Sprintf("From %s to %s",
	//		time.Unix(start, 0).Format("2006-01-02 15:04:05"),
	//		time.Unix(end, 0).Format("2006-01-02 15:04:05")))

	return
}

func splitTmDesc(desc string) (start, end string) {
	runeArr := []rune(desc)

	index := -1
	bracketsOpen := false
	for i := 0; i < len(runeArr); i++ {
		c := runeArr[i]

		if c == constant.RightBrackets {
			bracketsOpen = false
		} else if c == constant.LeftBrackets {
			bracketsOpen = true
		}

		str := fmt.Sprintf("%c", c)
		if !bracketsOpen && str == "-" {
			index = i
			break
		}
	}

	if index == -1 {
		start = desc
	} else if index == 0 {
		end = desc[1:]
	} else if index == len(desc)-1 {
		start = desc[:index]
	} else {
		start = desc[:index]
		end = desc[index+1:]
	}

	if len(start) > 0 {
		start = strings.TrimPrefix(start, string(constant.LeftBrackets))
		start = strings.TrimSuffix(start, string(constant.RightBrackets))
	}
	if len(end) > 0 {
		end = strings.TrimPrefix(end, string(constant.LeftBrackets))
		end = strings.TrimSuffix(end, string(constant.RightBrackets))
	}

	return
}

func parseTsValue(str string, isStart bool) (value int64) {
	str = strings.TrimSpace(str)

	if strings.HasPrefix(str, "+") || strings.HasPrefix(str, "-") {
		value = time.Now().Unix()

		value = increment(value, str)

		return
	}

	loc, _ := time.LoadLocation("Local")
	tm, err := time.ParseInLocation("20060102 150405", str, loc)
	if err != nil {
		todayStart, todayEnd := getTodayTs()
		if isStart {
			value = todayStart
		} else {
			value = todayEnd
		}
	} else {
		value = tm.Unix()
	}

	return
}

func getTodayTs() (start, end int64) {
	now := time.Now()

	start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()
	end = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location()).Unix()

	return
}

// parse time expression to number and unit
// example: parseTsExpr("+3w") ->  int(3) and "w"
func parseTsExpr(incrementStr string) (num int, unit string) {
	num = 0
	unit = "s"
	incrementStr = strings.TrimSpace(incrementStr)
	if len(incrementStr) <= 0 {
		return
	}

	units := []string{"Y", "M", "D", "w", "h", "m", "s"}

	unit = string(incrementStr[len(incrementStr)-1])
	found, _ := stringUtils.FindInArr(unit, units)
	end := len(incrementStr) - 1
	if !found {
		unit = "s"
		end = len(incrementStr)
	}

	numStr := incrementStr[:end]
	num, _ = strconv.Atoi(numStr)
	return
}

func incrementUnit(originalVal int64, incNum int, unit string) (ret int64) {
	ret = originalVal

	start := time.Unix(ret, 0)
	if incNum == 0 {
		return
	}
	switch unit {
	case "Y":
		start = start.AddDate(incNum, 0, 0)
	case "M":
		start = start.AddDate(0, incNum, 0)
	case "D":
		start = start.AddDate(0, 0, incNum)
	case "w":
		start = start.AddDate(0, 0, incNum * 7)
	case "h":
		start = start.Add(time.Hour * time.Duration(incNum))
	case "m":
		start = start.Add(time.Minute * time.Duration(incNum))
	case "s":
		start = start.Add(time.Second * time.Duration(incNum))
	default:
	}
	ret = start.Unix()
	return
}

func increment(originalVal int64, incrementStr string) (ret int64) {
	ret = originalVal

	num, unit := parseTsExpr(incrementStr)
	if num == 0 {
		return
	}

	ret = incrementUnit(originalVal, num, unit)

	return
}

func GenerateTimeItems(start int64, end int64, step int, unit string, repeat int, repeatTag string) []interface{}{
	arr := make([]interface{}, 0)

	total := 0
	if repeatTag == "" {
		for i := 0; true; {
			val := incrementUnit(start, i * step, unit)
			//val := start + int64(i*step)
			if (val >= end && step > 0) || (val <= end && step < 0) {
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
				val := start + incrementUnit(start, i*step, unit)
				//val := start + int64(i*step)
				if (val >= end && step > 0) || (val <= end && step < 0) {
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
