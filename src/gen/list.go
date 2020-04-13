package gen

import (
	"github.com/easysoft/zendata/src/model"
	"strings"
	"unicode"
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
			continue // TODO: log
		}

		elemArr := strings.Split(item, "-")
		start := elemArr[0]
		end := ""
		if len(elemArr) > 1 {
			end = elemArr[1]
		}

		if len(start) > 1 || len(end) > 1 {
			continue // TODO: log
		}

		if unicode.IsDigit(rune(start[0])) && unicode.IsDigit(rune(end[0])) {

		} else if unicode.IsLetter(rune(start[0])) && unicode.IsLetter(rune(end[0])) {

		}

		index++
	}
}