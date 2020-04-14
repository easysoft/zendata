package gen

import (
	"github.com/easysoft/zendata/src/model"
	"strings"
	"unicode"
)

func GenerateList(field model.Field, total int, fieldMap map[string][]interface{}) {
	name := strings.TrimSpace(field.Name)
	rang := strings.TrimSpace(field.Range)
	prefix := field.Prefix
	postfix := field.Postfix

	rangeItems := strings.Split(rang, ",")

	index := 0
	for _, item := range rangeItems {
		if index >= total { break }
		if strings.TrimSpace(item) == "" { continue }

		elemArr := strings.Split(item, "-")
		start := elemArr[0]
		end := ""
		if len(elemArr) > 1 {
			end = elemArr[1]
		}

		if len(start) > 1 || len(end) > 1 { // must be one character
			continue
		}

		if (unicode.IsDigit(rune(start[0])) || unicode.IsLetter(rune(start[0]))) && // must be digit or letter
				 (unicode.IsDigit(rune(end[0])) || unicode.IsLetter(rune(end[0]))) {
			items := GenerateListItems(byte(start[0]), byte(end[0]), index, total)

			fieldMap[name] = append(fieldMap[name], items)
			index = index + len(items)
		}
			continue
		}

		index++
	}
}

func GenerateListItems(startCh byte, endCh byte, index int, total int) []byte {
	arr := make([]byte, 0)

	for i := 0; i < total - index; i++ {
		bt := startCh + byte(i)
		if bt >= endCh {
			break
		}

		arr = append(arr, bt)
	}

	return arr
}