package gen

import (
	"github.com/easysoft/zendata/src/model"
	dateUtils "github.com/easysoft/zendata/src/utils/date"
	"strconv"
	"strings"
)

func GenerateTimestamp(field *model.DefField, total int, fieldMap map[string][]interface{}) {
	name := strings.TrimSpace(field.Note)
	rang := strings.TrimSpace(field.Range)
	stepStr := strings.TrimSpace("field.Step")

	rangeItems := strings.Split(rang, ",")

	index := 0
	for _, item := range rangeItems {
		if index >= total {
			break
		}
		if strings.TrimSpace(item) == "" { continue }

		elemArr := strings.Split(item, "-")
		startStr := elemArr[0]
		endStr := ""
		if len(elemArr) > 1 {
			endStr = elemArr[1]
		}

		items := make([]interface{}, 0)

		start, err1 := dateUtils.DateStrToTimestamp(startStr)
		end, err2 := dateUtils.DateStrToTimestamp(endStr)
		if err1 == nil && err2 == nil {
			step, err3 := strconv.Atoi(stepStr)
			if err3 != nil {
				step = 1000 // default 1 sec
			}

			items = GenerateIntItems(start, end, int64(step), true)
		}

		fieldMap[name] = append(fieldMap[name], items...)
		index = index + len(items)
	}
}