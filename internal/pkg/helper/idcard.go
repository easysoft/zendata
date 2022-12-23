package helper

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const ()

func GenerateIdCard() (ret string) {
	areaCode := AreaCode[rand.Intn(len(AreaCode))] +
		fmt.Sprintf("%0*d", 4, randInt(1, 9999))

	birthday := RandDate().Format("20060102")

	randomCode := fmt.Sprintf("%0*d", 3, randInt(0, 999))
	prefix := areaCode + birthday + randomCode

	ret = prefix + VerifyCode(prefix)

	return
}

var AreaCode = []string{
	"11",
	"12",
	"13",
	"14",
	"15",
	"21",
	"22",
	"23",
	"31",
	"32",
	"33",
	"34",
	"35",
	"36",
	"37",
	"41",
	"42",
	"43",
	"44",
	"45",
	"46",
	"50",
	"51",
	"52",
	"53",
	"54",
	"61",
	"62",
	"63",
	"64",
	"65",
	"71",
	"81",
	"82",
	"91",
}

var AreaCodeMap = map[string]string{
	"11": "北京",
	"12": "天津",
	"13": "河北",
	"14": "山西",
	"15": "内蒙古",
	"21": "辽宁",
	"22": "吉林",
	"23": "黑龙江",
	"31": "上海",
	"32": "江苏",
	"33": "浙江",
	"34": "安徽",
	"35": "福建",
	"36": "江西",
	"37": "山东",
	"41": "河南",
	"42": "湖北",
	"43": "湖南",
	"44": "广东",
	"45": "广西",
	"46": "海南",
	"50": "重庆",
	"51": "四川",
	"52": "贵州",
	"53": "云南",
	"54": "西藏",
	"61": "陕西",
	"62": "甘肃",
	"63": "青海",
	"64": "宁夏",
	"65": "新疆",
	"71": "台湾",
	"81": "香港",
	"82": "澳门",
	"91": "国外",
}

// VerifyCode 通过给定的身份证号生成最后一位的 VerifyCode
func VerifyCode(cardId string) string {
	tmp := 0
	for i, v := range Wi {
		t, _ := strconv.Atoi(string(cardId[i]))
		tmp += t * v
	}
	return ValCodeArr[tmp%11]
}

var ValCodeArr = []string{
	"1", "0", "X", "9", "8", "7", "6", "5", "4", "3", "2",
}

var Wi = []int{
	7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2,
}

// 指定范围随机 int
func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}

// 指定范围随机 int64
func randInt64(min, max int64) int64 {
	return min + rand.Int63n(max-min)
}

// RandDate 返回随机时间，时间区间从 1970 年 ~ 2020 年
func RandDate() time.Time {
	begin, _ := time.Parse("2006-01-02 15:04:05", "1970-01-01 00:00:00")
	end, _ := time.Parse("2006-01-02 15:04:05", "2020-01-01 00:00:00")
	return time.Unix(randInt64(begin.Unix(), end.Unix()), 0)
}
