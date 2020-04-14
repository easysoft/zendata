package gen

import (
	"github.com/easysoft/zendata/src/model"
	"strconv"
	"strings"
)

func GenerateList(field model.Field, total int, fieldMap map[string][]interface{}) {
	name := strings.TrimSpace(field.Name)
	rang := strings.TrimSpace(field.Range)

	rangeItems := strings.Split(rang, ",")

	index := 0
	for itemIndex, item := range rangeItems {
		if index >= total { break }
		if strings.TrimSpace(item) == "" { continue }

		elemArr := strings.Split(item, "-")
		startStr := elemArr[0]
		endStr := ""
		if len(elemArr) > 1 {
			endStr = elemArr[1]
		}

		items := make([]interface{}, 0)
		isLast := itemIndex == len(rangeItems) - 1

		startInt, err1 := strconv.ParseInt(startStr, 0, 64)
		endInt, err2 := strconv.ParseInt(endStr, 0, 64)
		if err1 == nil && err2 == nil { // int
			items = GenerateIntItems(startInt, endInt, index, total, isLast)
		} else {
			//startFloat, err1 := strconv.ParseFloat(startStr, 64)
			//endFloat, err2 := strconv.ParseFloat(endStr, 64)
			//if err1 == nil && err2 == nil { // float
			//
			//} else if len(startStr) > 1 && len(endStr) > 1 { // single character
			//	items = GenerateByteItems(byte(startStr[0]), byte(startStr[0]), index, total)
			//}
		}

		fieldMap[name] = append(fieldMap[name], items...)
		index = index + len(items)
	}


}

func GenerateIntItems(start int64, end int64, index int, total int, isLast bool) []interface{} {
	arr := make([]interface{}, 0)

	count := index
	for i := 0; i < total - index; {
		val := start + int64(i)

		if val > end {
			if isLast && count < total { // loop if it's last item and not enough
				i = 0
				continue
			} else {
				break
			}
		}

		arr = append(arr, val)
		count++
		i++
	}

	return arr
}

func GenerateByteItems(start byte, end byte, index int, total int) []byte {
	arr := make([]byte, 0)

	for i := 0; i < total - index; i++ {
		bt := start + byte(i)
		if bt >= end {
			break
		}

		arr = append(arr, bt)
	}

	return arr
}