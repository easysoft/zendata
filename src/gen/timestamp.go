package gen

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	"strconv"
	"strings"
	"time"
)

func CreateTimestampField(field *model.DefField, fieldWithValue *model.FieldWithValues) {
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
		if index >= constant.MaxNumb { break }

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
	if endStr == "" {
		endStr = startStr
	}

	start = parseTsValue(startStr, true)
	end = parseTsValue(endStr, false)

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

	tm, err := time.Parse("20060102", str)
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