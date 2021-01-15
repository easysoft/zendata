package gen

import (
	"fmt"
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	logUtils "github.com/easysoft/zendata/src/utils/log"
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
	desc, step := parseTsSection(section)
	start, end := parseTsDesc(desc)

	if step > 0 && start > end {
		step *= -1
	}

	// get index numbers for data retrieve
	numbs := GenerateIntItems(start, end, step, false, 1)

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

func parseTsSection(section string) (desc string, step int) {
	section = strings.TrimSpace(section)

	sectionArr := strings.Split(section, ":")
	desc = sectionArr[0]
	step = 1
	if len(sectionArr) > 1 {
		stepStr := sectionArr[1]
		stepTemp, err := strconv.Atoi(stepStr)

		if err == nil {
			step = stepTemp
		}
	}

	return
}

func parseTsDesc(desc string) (start, end int64) {
	desc = strings.TrimSpace(desc)

	if desc == "today" {
		start, end = getTodayTs()
		return
	}

	arr := strings.Split(desc, "-")
	startStr := arr[0]
	endStr := ""
	if len(arr) > 1 {
		endStr = arr[1]
	}

	start = parseTsValue(startStr, true)
	end = parseTsValue(endStr, false)

	if start > end {
		temp := start
		start = end
		end = temp
	}

	logUtils.PrintTo(
		fmt.Sprintf("From %s to %s",
			time.Unix(start, 0).Format("2006-01-02 15:04:05"),
			time.Unix(end, 0).Format("2006-01-02 15:04:05")))

	return
}

func parseTsValue(str string, isStart bool) (value int64) {
	str = strings.TrimSpace(str)

	if strings.Contains(str, "now") {
		value = time.Now().Unix()
		return
	} else if str == "today" {
		start, end := getTodayTs()
		if isStart {
			value = start
		} else {
			value = end
		}

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
		if !isStart {
			tm = time.Date(tm.Year(), tm.Month(), tm.Day(), 23, 59, 59, 0, tm.Location())
		}
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
