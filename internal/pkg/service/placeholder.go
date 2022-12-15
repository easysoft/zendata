package service

import (
	"fmt"
	"github.com/easysoft/zendata/internal/pkg/gen/helper"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"regexp"
	"strconv"
	"strings"
)

type PlaceholderService struct {
}

func (s *PlaceholderService) ReplacePlaceholder(val string) (ret string) {
	ret = val

	re := regexp.MustCompile("(?siU)\\${(.*)}")
	matchResultArr := re.FindAllStringSubmatch(ret, -1)
	matchTimes := len(matchResultArr)

	for _, childArr := range matchResultArr {
		placeholderStr := childArr[1]
		values := s.getValForPlaceholder(placeholderStr, matchTimes)

		for _, str := range values {
			key, _ := strconv.Atoi(placeholderStr)
			temp := s.PlaceholderStr(key)
			ret = strings.Replace(ret, temp, str, 1)
		}
	}

	return
}

func (s *PlaceholderService) getValForPlaceholder(placeholderStr string, count int) []string {
	placeholderNo, _ := strconv.Atoi(placeholderStr)
	mp := vari.GlobalVars.RandFieldSectionPathToValuesMap[placeholderNo]

	tp := mp["type"].(string)
	repeatObj := mp["repeat"]

	repeat := 1
	if repeatObj != nil {
		repeat = repeatObj.(int)
	}

	strArr := make([]string, 0)
	repeatTag := mp["repeatTag"].(string)
	if tp == "int" {
		start := mp["start"].(string)
		end := mp["end"].(string)
		precision := mp["precision"].(string)
		format := mp["format"].(string)

		strArr = helper.GetRandFromRange("int", start, end, "1",
			repeat, repeatTag, precision, format, count)

	} else if tp == "float" {
		start := mp["start"].(string)
		end := mp["end"].(string)
		stepStr := fmt.Sprintf("%v", mp["step"])

		precision := mp["precision"].(string)
		format := mp["format"].(string)

		strArr = helper.GetRandFromRange("float", start, end, stepStr,
			repeat, repeatTag, precision, format, count)

	} else if tp == "char" {
		start := mp["start"].(string)
		end := mp["end"].(string)
		precision := mp["precision"].(string)
		format := mp["format"].(string)

		strArr = helper.GetRandFromRange("char", start, end, "1",
			repeat, repeatTag, precision, format, count)

	} else if tp == "list" {
		list := mp["list"].([]string)
		strArr = helper.GetRandFromList(list, repeat, count)

	}

	strArr = strArr[:count]
	return strArr
}

func (s *PlaceholderService) PlaceholderStr(key int) string {
	return fmt.Sprintf("${%d}", key)
}

func (s *PlaceholderService) PlaceholderMapForRandValues(tp string, list []string, start, end, step, precision, format string,
	repeat int, repeatTag string) map[string]interface{} {
	ret := map[string]interface{}{}

	ret["type"] = tp

	// for literal values
	ret["list"] = list

	// for interval values
	ret["start"] = start
	ret["end"] = end
	ret["step"] = step
	ret["precision"] = precision
	ret["format"] = format

	ret["repeat"] = repeat
	ret["repeatTag"] = repeatTag

	return ret
}

func (s *PlaceholderService) GetRandFieldSectionKey(pth string) (key int) {
	max := 0

	for k, v := range vari.GlobalVars.RandFieldSectionShortKeysToPathMap {
		if pth == v {
			key = k
			return
		}

		if k > max {
			max = k
		}
	}

	if key == 0 {
		key = max + 1
	}

	return
}
