package service

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/easysoft/zendata/internal/pkg/helper"
	commonUtils "github.com/easysoft/zendata/pkg/utils/common"
	"github.com/easysoft/zendata/pkg/utils/vari"
)

type PlaceholderService struct {
	FixService *FixService `inject:""`
}

func (s *PlaceholderService) ReplacePlaceholder(val string) (ret interface{}) {
	ret = val

	re := regexp.MustCompile("(?siU)\\${(.*)}")
	matchResultArr := re.FindAllStringSubmatch(fmt.Sprintf("%v", ret), -1)
	matchTimes := len(matchResultArr)

	for _, childArr := range matchResultArr {
		placeholderStr := childArr[1]
		values := s.getValForPlaceholder(placeholderStr, matchTimes)

		for _, val := range values {
			key, _ := strconv.Atoi(placeholderStr)
			temp := s.PlaceholderStr(key)

			str := fmt.Sprintf("%v", ret)
			isNotStr := matchTimes == 1 && commonUtils.GetType(val) != "string" &&
				strings.Index(str, "${") == 0 && str[len(str)-1] == '}'

			if isNotStr {
				ret = val
			} else {
				ret = strings.Replace(fmt.Sprintf("%v", ret), temp, fmt.Sprintf("%v", val), 1)
				ret = s.FixService.TrimIfFormatIsNotText(fmt.Sprintf("%v", ret))
			}
		}
	}

	return
}

func (s *PlaceholderService) getValForPlaceholder(placeholderStr string, count int) (ret []interface{}) {
	placeholderNo, _ := strconv.Atoi(placeholderStr)
	mp := vari.GlobalVars.RandFieldSectionPathToValuesMap[placeholderNo]

	tp := mp["type"].(string)
	repeatObj := mp["repeat"]

	repeat := 1
	if repeatObj != nil {
		repeat = repeatObj.(int)
	}

	repeatTag := mp["repeatTag"].(string)
	if tp == "int" {
		start := mp["start"].(string)
		end := mp["end"].(string)
		precision := mp["precision"].(string)
		format := mp["format"].(string)

		ret = helper.GetRandValuesFromRange("int", start, end, "1", repeat, repeatTag, precision, format, count)

	} else if tp == "float" {
		start := mp["start"].(string)
		end := mp["end"].(string)
		stepStr := fmt.Sprintf("%v", mp["step"])

		precision := mp["precision"].(string)
		format := mp["format"].(string)

		ret = helper.GetRandValuesFromRange("float", start, end, stepStr,
			repeat, repeatTag, precision, format, count)

	} else if tp == "char" {
		start := mp["start"].(string)
		end := mp["end"].(string)
		precision := mp["precision"].(string)
		format := mp["format"].(string)

		ret = helper.GetRandValuesFromRange("char", start, end, "1",
			repeat, repeatTag, precision, format, count)

	} else if tp == "list" {
		list := mp["list"].([]string)
		ret = helper.GetRandFromList(list, repeat, count)

	}

	ret = ret[:count]

	return
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
