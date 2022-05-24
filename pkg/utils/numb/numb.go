package numbUtils

import (
	"math"
	"strings"
)

var num2char = "0123456789abcdefghijklmnopqrstuvwxyz"

func NumToBHex(num int) string {
	n := 36

	numStr := ""
	for num != 0 {
		yu := num % n
		numStr = string(num2char[yu]) + numStr
		num = num / n
	}

	numStr = strings.ToUpper(numStr)
	numStr = strings.Repeat("0", 5-len(numStr)) + numStr
	return numStr
}

func BHex2Num(str string) int {
	n := 36

	str = strings.ToLower(str)
	v := 0.0
	length := len(str)
	for i := 0; i < length; i++ {
		s := string(str[i])
		index := strings.Index(num2char, s)
		v += float64(index) * math.Pow(float64(n), float64(length-1-i)) // 倒序
	}
	return int(v)
}
