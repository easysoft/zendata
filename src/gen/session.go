package gen

import (
	"github.com/easysoft/zendata/src/model"
	"strconv"
	"strings"
)

func GenerateSessionId(field *model.Field, total int, fieldMap map[string][]interface{}) {
	name := strings.TrimSpace(field.Note)
	rang := strings.TrimSpace(field.Range)
	stepStr := strings.TrimSpace(field.Step)

	rangeItems := strings.Split(rang, ",")

	index := 0
	for itemIndex, item := range rangeItems {
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
		isLast := itemIndex == len(rangeItems) - 1

		startInt, errInt1 := strconv.ParseInt(startStr, 0, 64)
		endInt, errInt2 := strconv.ParseInt(endStr, 0, 64)
		if errInt1 == nil && errInt2 == nil { // int
			step, errInt3 := strconv.Atoi(stepStr)
			if errInt3 != nil {
				step = 1
			}

			items = GenerateIntItems(startInt, endInt, int64(step), index, total, isLast)
		}

		fieldMap[name] = append(fieldMap[name], items...)
		index = index + len(items)
	}
}