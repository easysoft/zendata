package gen

import (
	"fmt"
	"github.com/easysoft/zendata/src/utils/const"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

/**
	split field range string with comma to a array, ignore the comma in []
	1-2:R,[user1,user2]{2} -> 1-2:R
                              [user1,user2]{2}
 */
func ParseRange(rang string) []string {
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
func ParseRangeItem(item string) (entry string, step string, repeat int) {
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
func ParseEntry(str string) (typ string, desc string) {
	desc = strings.TrimSpace(str)

	if strings.Contains(desc, ",") || strings.Contains(desc, "`") || !strings.Contains(desc, "-") {
		typ = "literal"
	} else  {
		temp := strings.ReplaceAll(desc, string(constant.LeftBrackets), "")
		temp = strings.ReplaceAll(temp,string(constant.RightBrackets), "")

		if isBoundaryStr(temp) {
			typ = "interval"
			desc = temp
		} else {
			typ = "literal"
		}
	}

	return
}

func isBoundaryStr(str string) bool {
	arr := strings.Split(str, "-")
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
