package gen

import (
	"github.com/easysoft/zendata/src/model"
	"strings"
)

func GenerateList(field model.Field, count int, fieldMap map[string][]interface{}) {
	//for i := 0; i < count; i++ {
	//}

	name := strings.TrimSpace(field.Name)
	rang := strings.TrimSpace(field.Range)
	prefix := field.Prefix
	postfix := field.Postfix

	rangeItems := strings.Split(rang, ",")

	index := 0
	for _, item := range rangeItems {
		if index >= count {
			break
		}
		if strings.TrimSpace(item) == "" {
			continue
		}

		elemArr := strings.Split(item, "-")
		start := elemArr[0]
		if len(elemArr) > 1 {
			end := elemArr[0]
		}

		index++
	}
}