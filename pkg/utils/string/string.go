package stringUtils

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"strconv"
	"strings"
	"unicode"

	constant "github.com/easysoft/zendata/internal/pkg/const"
)

func TrimAll(str string) string {
	str = strings.Trim(str, "\n")
	str = strings.TrimSpace(str)

	return str
}

func Ucfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}
func UcAll(str string) string {
	ret := ""

	for _, v := range str {
		ret += string(unicode.ToUpper(v))
	}
	return ret
}

func Lcfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

func U2s(form string) (to string, err error) {
	bs, err := hex.DecodeString(strings.Replace(form, `\u`, ``, -1))
	if err != nil {
		return
	}
	for i, bl, br, r := 0, len(bs), bytes.NewReader(bs), uint16(0); i < bl; i += 2 {
		binary.Read(br, binary.BigEndian, &r)
		to += string(r)
	}
	return
}

func BoolToPass(b bool) string {
	if b {
		return constant.PASS.String()
	} else {
		return constant.FAIL.String()
	}
}

func FindInArrBool(str string, arr []string) bool {
	found, _ := FindInArr(str, arr)
	return found
}
func FindInArr(str string, arr []string) (bool, int) {
	for index, s := range arr {
		if str == s {
			return true, index
		}
	}

	return false, -1
}

func StrInArr(str string, arr []string) bool {
	found, _ := FindInArr(str, arr)
	return found
}
func InArray(need interface{}, arr []string) bool {
	for _, v := range arr {
		if need == v {
			return true
		}
	}
	return false
}

func GetNumbWidth(numb int) int {
	str := strconv.Itoa(numb)
	width := len(str)

	return width
}

func StartWith(str, sub string) bool {
	return strings.Index(str, sub) == 0
}
func EndWith(str, sub string) bool {
	return strings.Contains(str, sub) && strings.LastIndex(str, sub) == len(str)-len(sub)
}
