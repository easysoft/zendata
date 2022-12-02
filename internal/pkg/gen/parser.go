package gen

import (
	"fmt"
	"github.com/easysoft/zendata/internal/pkg/const"
	valueGen "github.com/easysoft/zendata/internal/pkg/gen/value"
	"github.com/easysoft/zendata/internal/pkg/model"
	stringUtils "github.com/easysoft/zendata/pkg/utils/string"
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

		if i == len(runeArr)-1 {
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

// for Literal only
func ParseDesc(desc string) (items []string) {
	desc = strings.TrimSpace(desc)
	desc = strings.Trim(desc, ",")

	if desc == "" {
		items = append(items, desc)
		return
	}

	runeArr := []rune(desc)

	if runeArr[0] == constant.Backtick && runeArr[len(runeArr)-1] == constant.Backtick { // `xxx`
		desc = string(runeArr[1 : len(runeArr)-1])
		items = append(items, desc)

	} else if runeArr[0] == constant.LeftBrackets && runeArr[len(runeArr)-1] == constant.RightBrackets { // [abc,123]
		desc = string(runeArr[1 : len(runeArr)-1])
		items = strings.Split(desc, ",")

	} else {
		items = append(items, desc)
	}

	return
}

/*
*

		convert range item to entity, step, repeat
		[user1,user2]{2} -> entry  =>[user1,user2]
	                        step   =>1
	                        repeat =>2
*/
func ParseRangeSection(rang string) (entry string, step string, repeat int, repeatTag string) {
	rang = strings.TrimSpace(rang)

	if rang == "" {
		repeat = 1
		return
	}

	runeArr := []rune(rang)
	if (runeArr[0] == constant.Backtick && runeArr[len(runeArr)-1] == constant.Backtick) || // `xxx`
		(string(runeArr[0]) == string(constant.LeftBrackets) && // (xxx)
			string(runeArr[len(runeArr)-1]) == string(constant.RightBrackets)) {

		entry = rang
		if repeat == 0 {
			repeat = 1
		}
		return
	}

	repeat, repeatTag, rangWithoutRepeat := ParseRepeat(rang)

	sectionArr := strings.Split(rangWithoutRepeat, ":")
	entry = sectionArr[0]
	if len(sectionArr) == 2 {
		step = strings.TrimSpace(strings.ToLower(sectionArr[1]))
	}

	if step != "" {
		pattern := "\\d+"
		isNum, _ := regexp.MatchString(pattern, step)
		if !isNum && step != "r" {
			entry = rang
			step = ""
		}
	}

	if repeat == 0 {
		repeat = 1
	}
	return entry, step, repeat, repeatTag
}

/*
*

		get range item entity's type and desc
		1-9 or [1-9]  -> type => interval
	                     desc => 1-9 or [1-9]
		[user1,user2] -> type => literal
	                     desc => user2,user3
*/
func ParseRangeSectionDesc(str string) (typ string, desc string) {
	desc = strings.TrimSpace(str)
	runeArr := []rune(desc)

	if desc == "" {
		typ = "literal"
		return
	}

	if stringUtils.EndWith(desc, ".yaml") { // refer to another yaml file
		typ = "yaml"
		return
	}

	if string(runeArr[0]) == string(constant.LeftBrackets) && // [a-z,1-9,userA,UserB]
		string(runeArr[len(runeArr)-1]) == string(constant.RightBrackets) {

		desc = removeBoundary(desc)
		arr := strings.Split(desc, ",")

		temp := ""
		for _, item := range arr {
			if isScopeStr(item) && isCharOrNumberScope(item) { // only support a-z and 0-9 in []
				tempField := model.DefField{}
				values := CreateValuesFromInterval(&tempField, item, "", 1, "")

				for _, val := range values {
					temp += valueGen.InterfaceToStr(val) + ","
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
	} else {
		temp := removeBoundary(desc)

		if isScopeStr(temp) {
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

func isScopeStr(str string) bool {
	arr := strings.Split(str, "-")
	if len(arr) < 2 || strings.TrimSpace(str) == "-" {
		return false
	}

	left := strings.TrimSpace(arr[0])
	right := strings.TrimSpace(arr[1])

	if len(left) != 1 || len(right) != 1 { // more than on char, must be number
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

func isCharOrNumberScope(str string) bool {
	arr := strings.Split(str, "-")
	if len(arr) < 2 {
		return false
	}

	left := strings.TrimSpace(arr[0])
	right := strings.TrimSpace(arr[1])

	if len(left) == 1 && len(right) == 1 {
		return true
	}

	return false
}

func ParseRepeat(rang string) (repeat int, repeatTag, rangeWithoutRepeat string) {
	repeat = 1

	regx := regexp.MustCompile(`\{(.*)!?\}`)
	arr := regx.FindStringSubmatch(rang)
	tag := ""
	if len(arr) == 2 {
		str := strings.TrimSpace(arr[1])
		if str[len(str)-1:] == "!" {
			tag = str[len(str)-1:]
			str = strings.TrimSpace(str[:len(str)-1])
		}
		repeat, _ = strconv.Atoi(str)
	}
	repeatTag = tag
	rangeWithoutRepeat = regx.ReplaceAllString(rang, "")

	return
}
