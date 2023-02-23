package genHelper

import (
	"fmt"
	valueGen "github.com/easysoft/zendata/internal/pkg/gen/value"
	"github.com/easysoft/zendata/internal/pkg/helper"
	commonUtils "github.com/easysoft/zendata/pkg/utils/common"
	"strconv"
	"strings"
)

func GetRandFromList(list []string, repeat, count int) (ret []interface{}) {
	for i := 0; i < count; i++ {
		rand := commonUtils.RandNum(len(list))
		val := list[rand]

		items := make([]string, 0)
		for round := 0; round < repeat; round++ {
			items = append(items, val)
		}

		ret = append(ret, strings.Join(items, ""))
	}

	return ret
}

func GetRandValuesFromRange(dataType, start, end, step string, repeat int, repeatTag, precisionStr string,
	format string, count int) (ret []interface{}) {

	precision, _ := strconv.Atoi(precisionStr)

	items := make([]interface{}, 0)

	if dataType == "int" {
		startInt, _ := strconv.ParseInt(start, 0, 64)
		endInt, _ := strconv.ParseInt(end, 0, 64)
		stepInt, _ := strconv.ParseInt(step, 0, 64)

		if endInt < startInt && stepInt > 0 {
			stepInt = stepInt * -1
		}

		items = valueGen.GenerateItems(startInt, endInt, stepInt, 0, true, repeat, repeatTag, count)

	} else if dataType == "char" {
		startChar := start[0]
		endChar := end[0]
		stepInt, _ := strconv.ParseInt(step, 10, 64)

		if int64(endChar) < int64(startChar) && stepInt > 0 {
			stepInt = stepInt * -1
		}

		items = valueGen.GenerateItems(startChar, endChar, stepInt, 0, true, repeat, repeatTag, count)

	} else if dataType == "float" {

		startFloat, _ := strconv.ParseFloat(start, 64)
		endFloat, _ := strconv.ParseFloat(end, 64)
		stepFloat, _ := strconv.ParseFloat(step, 64)

		if endFloat < startFloat && stepFloat > 0 {
			stepFloat = stepFloat * -1
		}

		items = valueGen.GenerateItems(startFloat, endFloat, stepFloat, precision, true, repeat, repeatTag, count)

	}

	for _, item := range items {
		typ := commonUtils.GetType(item)
		val := item

		if format != "" { // need to format for string
			val = getFormatStr(item, precision, format)
		} else if typ == "float" && precision != 0 {
			val = commonUtils.ChanePrecision(item.(float64), precision)
		}

		ret = append(ret, val)
	}

	return
}

func getFormatStr(val interface{}, precision int, format string) (ret string) {
	typ := valueGen.GetType(val)

	if format == "" {
		if typ == "int" {
			ret = fmt.Sprintf("%d", val)
		} else if typ == "char" {
			ret = string(val.(uint8))
		} else if typ == "float" {
			ret = strconv.FormatFloat(val.(float64), 'f', precision, 64)
		}

	} else {
		formatVal, success := helper.FormatStr(format, val, 0)
		if success {
			ret = formatVal
		} else {
			ret = fmt.Sprintf("%d", val)
		}
	}

	return
}
