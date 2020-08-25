package gen

import (
	"fmt"
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/utils/const"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func ParseRangeProperty(rang string) []string {
	items := make([]string, 0)

	bracketsOpen := false
	backtickOpen := false
	temp := ""

	rang = strings.Trim(rang, ",")
	runeArr := []rune(rang)

	for i := 0; i < len(runeArr); i++ {
		c := runeArr[i]

		if c == constant.RightBrackets {
			bracketsOpen = false
		} else if c == constant.LeftBrackets {
			bracketsOpen = true
		} else if !backtickOpen && c == constant.Backtick {
			backtickOpen = true
		} else if backtickOpen && c == constant.Backtick {
			backtickOpen = false
		}

		if i == len(runeArr) - 1 {
			temp += fmt.Sprintf("%c", c)
			items = append(items, temp)
		} else if !bracketsOpen && !backtickOpen && c == ',' {
			items = append(items, temp)
			temp = ""
			bracketsOpen = false
			backtickOpen = false
		} else {
			temp += fmt.Sprintf("%c", c)
		}
	}

	return items
}

func ParseDesc(desc string) (items []string) {
	desc = strings.TrimSpace(desc)
	desc = strings.Trim(desc, ",")
	runeArr := []rune(desc)

	if runeArr[0] == constant.Backtick &&  runeArr[len(runeArr) - 1] == constant.Backtick { 	// `xxx`
		desc = string(runeArr[1 : len(desc) - 1])
		items = append(items, desc)

	} else if runeArr[0] == constant.LeftBrackets &&  runeArr[len(runeArr) - 1] == constant.RightBrackets { // [abc,123]
		desc = string(runeArr[1 : len(desc) - 1])
		items = strings.Split(desc, ",")

	} else {
		items = append(items, desc)
	}

	return
}

/**
	convert range item to entity, step, repeat
	[user1,user2]{2} -> entry  =>[user1,user2]
                        step   =>1
                        repeat =>2
*/
func ParseRangeSection(item string) (entry string, step string, repeat int) {
	item = strings.TrimSpace(item)
	runeArr := []rune(item)
	if (runeArr[0] == constant.Backtick &&  runeArr[len(runeArr) - 1] == constant.Backtick) || // `xxx`
		(string(item[0]) == string(constant.LeftBrackets) && // (xxx)
			string(item[len(item) - 1]) == string(constant.RightBrackets)) {

		entry = item
		if repeat == 0 { repeat = 1 }
		if step == "" { step = "1" }
		return
	}

	regx := regexp.MustCompile(`\{(.*)\}`)
	arr := regx.FindStringSubmatch(item)
	if len(arr) == 2 {
		repeat, _ = strconv.Atoi(arr[1])
	}
	itemWithoutRepeat := regx.ReplaceAllString(item, "")

	sectionArr := strings.Split(itemWithoutRepeat, ":")
	entry = sectionArr[0]
	if len(sectionArr) == 2 {
		step = strings.ToLower(sectionArr[1])
	}

	if repeat == 0 { repeat = 1 }
	if step == "" { step = "1" }
	return entry, step, repeat
}

/**
	get range item entity's type and desc
	1-9 or [1-9]  -> type => interval
                     desc => 1-9 or [1-9]
	[user1,user2] -> type => literal
                     desc => user2,user3
*/
func ParseRangeSectionDesc(str string) (typ string, desc string) {
	desc = strings.TrimSpace(str)

	if stringUtils.EndWith(desc, ".yaml") { // refer to another yaml file
		typ = "yaml"
		return
	}

	if string(desc[0]) == string(constant.LeftBrackets) && // [a-z,1-9,userA,UserB]
		string(desc[len(desc)-1]) == string(constant.RightBrackets) {

		desc = removeBoundary(desc)
		arr := strings.Split(desc, ",")

		temp := ""
		for _, item := range arr {
			if isBoundaryStr(item) {
				tempField := model.DefField{}
				values := CreateValuesFromInterval(&tempField, item, "1", 1)

				for _, val := range values {
					temp += InterfaceToStr(val) + ","
				}
			} else {
				temp += item + ","
			}
		}

		temp = strings.TrimSuffix(temp, ",")
		desc = string(constant.LeftBrackets) + temp + string(constant.RightBrackets)
		typ = "literal"

		return
	}

	if strings.Contains(desc, ",") || strings.Contains(desc, "`") || !strings.Contains(desc, "-") {
		typ = "literal"
	} else  {
		temp := removeBoundary(desc)

		if isBoundaryStr(temp) {
			typ = "interval"
			desc = temp
		} else {
			typ = "literal"
		}
	}

	return
}

func removeBoundary(str string) string {
	str = strings.TrimLeft(str, string(constant.LeftBrackets))
	str = strings.TrimRight(str, string(constant.RightBrackets))

	return str
}

func isBoundaryStr(str string) bool {
	arr := strings.Split(str, "-")
	if len(arr) < 2 {
		return false
	}

	left := strings.TrimSpace(arr[0])
	right := strings.TrimSpace(arr[1])

	if len(left) != 1 || len(right) != 1 {
		leftRune := []rune(string(left[0]))[0]
		rightRune := []rune(string(right[0]))[0]

		if unicode.IsNumber(leftRune) && unicode.IsNumber(rightRune) {
			return true
		} else {
			return false
		}
	} else {
		leftRune := []rune(string(left[0]))[0]
		rightRune := []rune(string(right[0]))[0]

		if (unicode.IsLetter(leftRune) && unicode.IsLetter(rightRune)) ||
			(unicode.IsNumber(leftRune) && unicode.IsNumber(rightRune)) {
			return true
		} else {
			return false
		}
	}
}
