package gen

import (
	"fmt"
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	"regexp"
	"strconv"
	"strings"
)

func GenerateIP(field *model.Field, total int, fieldMap map[string][]interface{}) {
	name := strings.TrimSpace(field.Note)
	rang := strings.TrimSpace(field.Range)

	// 192.168.[0-9].[1-254]
	//regx := regexp.MustCompile(`\[([\d\-]*)\]`)
	//sections := regx.FindAllSubmatch(rang)
	//if len(sections) == 1 {
	//	i := 0
	//	i++
	//}

	rangeItems := strings.Split(rang, ",")

	index := 0
	for itemIndex, item := range rangeItems {
		if index >= total {
			break
		}
		if strings.TrimSpace(item) == "" { continue }

		elemArr := strings.Split(item, "-")
		start := elemArr[0]
		end := ""
		if len(elemArr) > 1 {
			end = elemArr[1]
		}

		items := make([]interface{}, 0)
		isLast := itemIndex == len(rangeItems) - 1

		if CheckIPFormat(start, false) && CheckIPFormat(end, true) {
			items = GenerateIpItems(start, end, index, total, isLast)
		}

		fieldMap[name] = append(fieldMap[name], items...)
		index = index + len(items)
	}
}

func CheckIPFormat(ip string, canBeEmpty bool) bool {
	if canBeEmpty && ip == "" {
		return true
	}

	regx := `^(([1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.)(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){2}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`
	pass, _ := regexp.MatchString(regx, ip)
	return pass
}

func GenerateIpItems(start string, end string, index int, total int, isLast bool) []interface{} {
	startNumb := ConvertIPToNumb(start)
	endNumb := ConvertIPToNumb(end)

	arr := make([]interface{}, 0)

	count := index
	for i := 0; i < total - index; {
		if count >= total {
			break
		}

		nextNumb := AddIPNumb(startNumb, &i)
		if nextNumb > endNumb {
			if isLast && count < total { // loop if it's last item and not enough
				i = 0
				continue
			} else {
				break
			}
		}

		ipStr := ConvertNumbToIP(nextNumb)
		arr = append(arr, ipStr)
		count++
		i++
	}

	return arr
}

func AddIPNumb(ipNumb int, increment *int) int {
	nextNumb := ipNumb + *increment * 1
	if nextNumb % constant.Power3 == 0 {
		*increment = *increment + 1
		nextNumb = nextNumb + 1
	}

	return nextNumb
}

func ConvertNumbToIP(ipNumb int) string {
	part4 := ipNumb / constant.Power3
	part3 := (ipNumb - part4 * constant.Power3) / constant.Power2
	part2 := (ipNumb - part4 * constant.Power3 - part3 * constant.Power2) / constant.Power2
	part1 := ipNumb % constant.Power1

	return fmt.Sprintf("%d.%d.%d.%d", part4, part3, part2, part1)
}

/**
	convert ip to a 255 radix integer
 */
func ConvertIPToNumb(ip string) int  {
	if ip == "" {
		return 2147483647
	}

	arr := strings.Split(ip, ".")
	numb := 0

	base, _ := strconv.Atoi(arr[3])
	numb = numb + base * (255 * 255 * 255)

	base, _ = strconv.Atoi(arr[2])
	numb = numb + base * (255 * 255)

	base, _ = strconv.Atoi(arr[1])
	numb = numb + base * (255)

	base, _ = strconv.Atoi(arr[0])
	numb = numb + base

	return numb
}