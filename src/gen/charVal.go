package gen

import (

)

func GenerateByteItems(start byte, end byte, step int, index int, total int) []interface{} {
	arr := make([]interface{}, 0)

	count := index
	for i := 0; i < total - index; {
		val := start + byte(i * step)

		if val > end {
			if count < total { // loop if it's last item and not enough
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