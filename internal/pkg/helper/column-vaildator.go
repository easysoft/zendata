package helper

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

func IsMobilePhone(str string) (ret bool) {
	regx := regexp.MustCompile(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`)
	ret = regx.MatchString(str)

	return
}

func IsTelPhone(str string) (ret bool) {
	regx := regexp.MustCompile(`[0-9\-\(\)（）]*[0-9]{4,}$`)
	ret = regx.MatchString(str)

	return
}

func IsIDCard(id string) bool {
	id = strings.ToUpper(id)
	if len(id) != 15 && len(id) != 18 {
		return false
	}
	r := regexp.MustCompile("(\\d{15})|(\\d{17}([0-9]|X))")
	if !r.MatchString(id) {
		return false
	}
	if len(id) == 15 {
		tm2, _ := time.Parse("01/02/2006", string([]byte(id)[8:10])+"/"+string([]byte(id)[10:12])+"/"+"19"+string([]byte(id)[6:8]))
		if tm2.Unix() <= 0 {
			return false
		}
		return true
	} else {
		tm2, _ := time.Parse("01/02/2006", string([]byte(id)[10:12])+"/"+string([]byte(id)[12:14])+"/"+string([]byte(id)[6:10]))
		if tm2.Unix() <= 0 {
			return false
		}
		// 检验18位身份证的校验码是否正确。
		// 校验位按照ISO 7064:1983.MOD 11-2的规定生成，X可以认为是数字10。
		arrInt := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
		arrCh := []string{"1", "0", "X", "9", "8", "7", "6", "5", "4", "3", "2"}
		sign := 0
		for k, v := range arrInt {
			intTemp, _ := strconv.Atoi(string([]byte(id)[k : k+1]))
			sign += intTemp * v
		}
		n := sign % 11
		valNum := arrCh[n]
		if valNum != string([]byte(id)[17:18]) {
			return false
		}
		return true
	}
}
