package main

import (
	"fmt"
	"strconv"
)

func main() {
	arrOfArr := [3][]string{}
	arrOfArr[0] = []string{"a","b","c","d","e"}
	arrOfArr[1] = []string{"1","2"}
	arrOfArr[2] = []string{"X","Y","Z"}

	indexArr := make([]int, 0)
	for _, _ = range arrOfArr {
		indexArr = append(indexArr, 0)
	}

	for i := 0; i < len(arrOfArr); i++ {
		loop := 1
		for j := i + 1; j < len(arrOfArr); j++ {
			loop = loop * len(arrOfArr[j])
		}

		indexArr[i] = loop
	}

	for i := 0; i < 100; i ++ {
		str := strconv.Itoa(i) + ": "
		for j := 0; j < len(arrOfArr); j++ {
			child := arrOfArr[j]

			mod := indexArr[j]
			remainder := i / mod % len(child)
			str = str + child[remainder]
		}
		fmt.Println(str)

		//addCur(&indexArr, arrOfArr)
	}
}