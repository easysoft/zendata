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

	tagOpen := false
	temp := ""

	rang = strings.Trim(rang, ",")
	runeArr := []rune(rang)

	for i := 0; i < len(runeArr); i++ {
		c := runeArr[i]

		if c == constant.RightChar {
			tagOpen = false
		} else if c == constant.LeftChar  {
			tagOpen = true
		}

		if i == len(runeArr) - 1 {
			temp += fmt.Sprintf("%c", c)
			items = append(items, temp)
		} else if !tagOpen && c == ',' {
			items = append(items, temp)
			temp = ""
			tagOpen = false
		} else {
			temp += fmt.Sprintf("%c", c)
		}
	}

	return items
}

/**
	convert range item to entity, step, repeat
	[user1,user2]{2} -> entry  =>[user1,user2]
                        step   =>1
                        repeat =>2
*/
func ParseRangeItem(item string) (entry string, step string, repeat int) {
	item = strings.TrimSpace(item)
	if string(item[0]) == string(constant.LeftChar) &&  // It's a whole when meet (xxx)
			string(item[len(item) - 1]) == string(constant.RightChar) {

		if repeat == 0 { repeat = 1 }
		if step == "" { step = "1" }
		return item, step, repeat
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
	str = strings.TrimSpace(str)
	desc = strings.ReplaceAll(str, string(constant.LeftChar), "")
	desc = strings.ReplaceAll(desc,string(constant.RightChar), "")

	if strings.Contains(desc, ",") || !strings.Contains(desc, "-") {
		typ = "literal"
	} else  {
		if !isBoundaryStr(desc) {
			typ = "literal"
		} else {
			typ = "interval"
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
