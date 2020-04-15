package gen

import (

)

func GenerateIntItems(start int64, end int64, step int64, index int, total int, isLast bool) []interface{} {
	arr := make([]interface{}, 0)

	count := index
	for i := 0; i < total - index; {
		if count >= total {
			break
		}

		val := start + int64(i) * step
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