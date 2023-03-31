package helper

import (
	uuid "github.com/satori/go.uuid"
	"regexp"
	"strings"
)

func GenerateToken(format string) (ret string) {
	ret = uuid.NewV4().String()
	sep := ""

	regx := regexp.MustCompile(`token\(\s*(\S+)\s*\)`)
	arr := regx.FindStringSubmatch(format)
	if len(arr) > 1 {
		sep = strings.Trim(arr[1], "'")
	}

	ret = strings.ReplaceAll(ret, "-", sep)

	return
}
