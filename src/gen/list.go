package gen

import (
	"github.com/easysoft/zendata/src/model"
	"strings"
)

func GenerateList(field model.Field, content *string) {
	name := strings.TrimSpace(field.Name)
	rang := strings.TrimSpace(field.Range)
	prefix := field.Prefix
	postfix := field.Postfix

	rangeItems := strings.Split(rang, ",")

	for _, item := range rangeItems {
		if strings.TrimSpace(item) == "" {
			continue
		}

		elemArr := strings.Split(item, "-")
		start := elemArr[0]
		if len(elemArr) > 1 {
			end := elemArr[0]
		}


	}
}