package service

import (
	"fmt"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/gen/helper"
	valueGen "github.com/easysoft/zendata/internal/pkg/gen/value"
	"github.com/easysoft/zendata/internal/pkg/model"
	commonUtils "github.com/easysoft/zendata/pkg/utils/common"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	stringUtils "github.com/easysoft/zendata/pkg/utils/string"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type RangeService struct {
	PlaceholderService *PlaceholderService `inject:""`
	ListService        *ListService        `inject:""`
	RandomService      *RandomService      `inject:""`

	DefService     *DefService     `inject:""`
	PrintService   *PrintService   `inject:""`
	CombineService *CombineService `inject:""`
	OutputService  *OutputService  `inject:""`
	FileService    *FileService    `inject:""`

	RangeService *RangeService `inject:""`
	MainService  *MainService  `inject:""`
}

func (s *RangeService) CreateFieldValuesFromRange(field *model.DefField) {
	rang := field.Range

	// gen empty values
	if rang == "" {
		for i := 0; i < vari.GlobalVars.Total; i++ {
			field.Values = append(field.Values, "")
			if strings.Index(field.Format, "uuid") == -1 {
				break
			}
		}

		return
	}

	// gen from field's range
	rangeSections := s.ParseRangeProperty(rang) // parse 1

	index := 0
	for _, rangeSection := range rangeSections {
		if index >= consts.MaxNumb {
			break
		}
		if rangeSection == "" {
			continue
		}

		descStr, stepStr, count, countTag := s.ParseRangeSection(rangeSection) // parse 2
		if strings.ToLower(stepStr) == "r" {
			(*field).IsRand = true
		}

		typ, desc := s.ParseRangeSectionDesc(descStr) // parse 3

		items := make([]interface{}, 0)
		if typ == "literal" {
			items = s.CreateValuesFromLiteral(field, desc, stepStr, count, countTag)
		} else if typ == "interval" {
			items = s.CreateValuesFromInterval(field, desc, stepStr, count, countTag)
		} else if typ == "yaml" { // load from a yaml
			items = s.CreateValuesFromYaml(field, desc, stepStr, count, countTag)
			field.ReferToAnotherYaml = true
		}

		for _, item := range items {
			field.Values = append(field.Values, item)
		}

		index = index + len(items)
	}

	if len(field.Values) == 0 {
		field.Values = append(field.Values, "N/A")
	}
}

func (s *RangeService) CreateValuesFromLiteral(field *model.DefField, desc string, stepStr string, repeat int, repeatTag string) (items []interface{}) {
	elemArr := s.ParseDesc(desc)
	step, _ := strconv.Atoi(stepStr)
	if step == 0 {
		step = 1
	}
	total := 0

	if field.Path != "" && stepStr == "r" {
		pth := field.Path
		key := helper.GetRandFieldSection(pth)

		items = append(items, s.PlaceholderService.PlaceholderStr(key))
		mp := s.PlaceholderService.PlaceholderMapForRandValues("list", elemArr, "", "", "", "",
			field.Format, repeat, repeatTag)

		vari.GlobalVars.RandFieldSectionShortKeysToPathMap[key] = pth
		vari.GlobalVars.RandFieldSectionPathToValuesMap[key] = mp
		return
	}

	if repeatTag == "" {
		for i := 0; i < len(elemArr); {
			idx := i
			if field.Path == "" && stepStr == "r" {
				idx = commonUtils.RandNum(len(elemArr)) // should set random here too
			}

			val := elemArr[idx]
			total = s.ListService.AppendValues(&items, val, repeat, total)

			if total >= consts.MaxNumb {
				break
			}
			i += step
		}
	} else if repeatTag == "!" {
		isRand := field.Path == "" && stepStr == "r"
		for i := 0; i < repeat; {
			total = s.ListService.AppendArrItems(&items, elemArr, total, isRand)

			if total >= consts.MaxNumb {
				break
			}
			i += step
		}
	}

	if field.Path == "" && stepStr == "r" { // for ranges and instances, random
		items = s.RandomService.RandomInterfaces(items)
	}

	return
}

func (s *RangeService) CreateValuesFromInterval(field *model.DefField, desc, stepStr string, repeat int, repeatTag string) (items []interface{}) {
	elemArr := strings.Split(desc, "-")
	startStr := elemArr[0]
	endStr := startStr
	if len(elemArr) > 1 {
		endStr = elemArr[1]
	}

	dataType, step, precision, rand, _ := s.CheckRangeType(startStr, endStr, stepStr)
	field.Precision = precision

	// 1. random replacement
	if field.Path != "" && dataType != "string" && rand { // random. for res, field.Path == ""
		pth := field.Path + "->" + desc
		key := helper.GetRandFieldSection(pth)

		val := s.PlaceholderService.PlaceholderStr(key)
		strItems := make([]string, 0)

		//for i := 0; i < repeat*count; i++ { // chang to add only one placeholder item
		items = append(items, val)
		strItems = append(strItems, val)
		//}

		mp := s.PlaceholderService.PlaceholderMapForRandValues(dataType, strItems, startStr, endStr, fmt.Sprintf("%v", step),
			strconv.Itoa(precision), field.Format, repeat, repeatTag)

		vari.GlobalVars.RandFieldSectionShortKeysToPathMap[key] = pth
		vari.GlobalVars.RandFieldSectionPathToValuesMap[key] = mp

		return
	}

	if dataType == "int" {
		startInt, _ := strconv.ParseInt(startStr, 0, 64)
		endInt, _ := strconv.ParseInt(endStr, 0, 64)

		items = valueGen.GenerateItems(startInt, endInt, int64(step.(int)), 0, rand, repeat, repeatTag, 0)

	} else if dataType == "char" {
		items = valueGen.GenerateItems(startStr[0], endStr[0], int64(step.(int)), 0, rand, repeat, repeatTag, 0)

	} else if dataType == "float" {
		startFloat, _ := strconv.ParseFloat(startStr, 64)
		endFloat, _ := strconv.ParseFloat(endStr, 64)
		field.Precision = precision

		items = valueGen.GenerateItems(startFloat, endFloat, step.(float64), field.Precision, rand, repeat, repeatTag, 0)

	} else if dataType == "string" {
		if repeat == 0 {
			repeat = 1
		}
		for i := 0; i < repeat; i++ {
			items = append(items, desc)
		}
	}

	if field.Path == "" && stepStr == "r" { // for ranges and instances, random again
		items = s.RandomService.RandomInterfaces(items)
	}

	return
}

func (s *RangeService) CreateValuesFromYaml(field *model.DefField, yamlFile, stepStr string, repeat int, repeatTag string) (items []interface{}) {
	// keep root def, since vari.ZdDef will be overwrite by refer yaml file
	rootDef := vari.GlobalVars.DefData
	configDir := vari.GlobalVars.ConfigFileDir
	exportFields := vari.GlobalVars.ExportFields
	res := vari.GlobalVars.ResData

	vari.GlobalVars.ExportFields = make([]string, 0) // set to empty to use all fields

	configFile := s.FileService.ComputerReferFilePath(yamlFile, field)

	vari.GlobalVars.ConfigFileDir = fileUtils.GetAbsDir(configFile)
	s.MainService.GenerateData([]string{configFile})

	// get the data
	items = s.OutputService.GenText(true)

	if field.Rand {
		items = s.RandomService.RandomValues(items)
	}

	if repeat > 0 {
		if repeat > len(items) {
			repeat = len(items)
		}
		items = items[:repeat]
	}

	// rollback root def when finish to deal with refer yaml file
	vari.GlobalVars.DefData = rootDef
	vari.GlobalVars.ConfigFileDir = configDir
	vari.GlobalVars.ExportFields = exportFields
	vari.GlobalVars.ResData = res

	return
}

func (s *RangeService) DealwithFixRange(field *model.DefField) {
	if s.isRangeFix(field.Prefix) {
		field.PrefixRange = s.CreateFieldFixValuesFromList(field.Prefix, field)
	} else {
		var tmp interface{}
		tmp = field.Prefix
		field.PrefixRange = &model.Range{Values: []interface{}{tmp}}
	}

	if s.isRangeFix(field.Postfix) {
		field.PostfixRange = s.CreateFieldFixValuesFromList(field.Postfix, field)
	} else {
		var tmp interface{}
		tmp = field.Postfix
		field.PostfixRange = &model.Range{Values: []interface{}{tmp}}
	}
}

func (s *RangeService) CreateFieldFixValuesFromList(strRang string, field *model.DefField) (rang *model.Range) {
	rang = &model.Range{}

	if strRang == "" {
		return
	}

	rangeSections := s.ParseRangeProperty(strRang) // parse 1

	index := 0
	for _, rangeSection := range rangeSections {
		if index >= vari.GlobalVars.Total {
			break
		}

		if rangeSection == "" {
			continue
		}

		descStr, stepStr, count, countTag := s.ParseRangeSection(rangeSection) // parse 2
		if strings.ToLower(stepStr) == "r" {
			rang.IsRand = true
		}

		typ, desc := s.ParseRangeSectionDesc(descStr) // parse 3

		items := make([]interface{}, 0)
		if typ == "literal" {
			items = s.CreateValuesFromLiteral(field, desc, stepStr, count, countTag)

		} else if typ == "interval" {
			items = s.CreateValuesFromInterval(field, desc, stepStr, count, countTag)

		} else if typ == "yaml" { // load from a yaml
			items = s.CreateValuesFromYaml(field, desc, stepStr, count, countTag)
			field.ReferToAnotherYaml = true

		}

		rang.Values = append(rang.Values, items...)
		index = index + len(items)
	}

	if len(rang.Values) == 0 {
		rang.Values = append(rang.Values, "N/A")
	}

	return
}

func (s *RangeService) isRangeFix(fix string) bool {
	index := strings.Index(fix, "-")

	return index > 0 && index < len(fix)-1
}

func (s *RangeService) ParseRangeProperty(rang string) []string {
	items := make([]string, 0)

	bracketsOpen := false
	backtickOpen := false
	temp := ""

	rang = strings.Trim(rang, ",")
	runeArr := []rune(rang)

	for i := 0; i < len(runeArr); i++ {
		c := runeArr[i]

		if c == consts.RightBrackets {
			bracketsOpen = false
		} else if c == consts.LeftBrackets {
			bracketsOpen = true
		} else if !backtickOpen && c == consts.Backtick {
			backtickOpen = true
		} else if backtickOpen && c == consts.Backtick {
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
func (s *RangeService) ParseDesc(desc string) (items []string) {
	desc = strings.TrimSpace(desc)
	desc = strings.Trim(desc, ",")

	if desc == "" {
		items = append(items, desc)
		return
	}

	runeArr := []rune(desc)

	if runeArr[0] == consts.Backtick && runeArr[len(runeArr)-1] == consts.Backtick { // `xxx`
		desc = string(runeArr[1 : len(runeArr)-1])
		items = append(items, desc)

	} else if runeArr[0] == consts.LeftBrackets && runeArr[len(runeArr)-1] == consts.RightBrackets { // [abc,123]
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
func (s *RangeService) ParseRangeSection(rang string) (entry string, step string, repeat int, repeatTag string) {
	rang = strings.TrimSpace(rang)

	if rang == "" {
		repeat = 1
		return
	}

	runeArr := []rune(rang)
	if (runeArr[0] == consts.Backtick && runeArr[len(runeArr)-1] == consts.Backtick) || // `xxx`
		(string(runeArr[0]) == string(consts.LeftBrackets) && // (xxx)
			string(runeArr[len(runeArr)-1]) == string(consts.RightBrackets)) {

		entry = rang
		if repeat == 0 {
			repeat = 1
		}
		return
	}

	repeat, repeatTag, rangWithoutRepeat := s.ParseRepeat(rang)

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
func (s *RangeService) ParseRangeSectionDesc(str string) (typ string, desc string) {
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

	if string(runeArr[0]) == string(consts.LeftBrackets) && // [a-z,1-9,userA,UserB]
		string(runeArr[len(runeArr)-1]) == string(consts.RightBrackets) {

		desc = s.removeBoundary(desc)
		arr := strings.Split(desc, ",")

		temp := ""
		for _, item := range arr {
			if isScopeStr(item) && s.isCharOrNumberScope(item) { // only support a-z and 0-9 in []
				tempField := model.DefField{}
				values := s.CreateValuesFromInterval(&tempField, item, "", 1, "")

				for _, val := range values {
					temp += valueGen.InterfaceToStr(val) + ","
				}
			} else {
				temp += item + ","
			}
		}

		temp = strings.TrimSuffix(temp, ",")
		desc = string(consts.LeftBrackets) + temp + string(consts.RightBrackets)
		typ = "literal"

		return
	}

	if strings.Contains(desc, ",") || strings.Contains(desc, "`") || !strings.Contains(desc, "-") {
		typ = "literal"
	} else {
		temp := s.removeBoundary(desc)

		if isScopeStr(temp) {
			typ = "interval"
			desc = temp
		} else {
			typ = "literal"
		}
	}

	return
}

func (s *RangeService) removeBoundary(str string) string {
	str = strings.TrimLeft(str, string(consts.LeftBrackets))
	str = strings.TrimRight(str, string(consts.RightBrackets))

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

func (s *RangeService) isCharOrNumberScope(str string) bool {
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

func (s *RangeService) ParseRepeat(rang string) (repeat int, repeatTag, rangeWithoutRepeat string) {
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

func (s *RangeService) CheckRangeType(startStr string, endStr string, stepStr string) (dataType string, step interface{}, precision int,
	rand bool, count int) {

	stepStr = strings.ToLower(strings.TrimSpace(stepStr))

	if start, end, stepi, ok := s.checkRangeTypeIsInt(startStr, endStr, stepStr); ok { // is int
		step = 1
		if stepStr == "r" {
			rand = true
		}

		count = (int)(start-end) / int(stepi)
		if count < 0 {
			count = count * -1
		}

		dataType = "int"
		step = int(stepi)
		return
	} else if start, end, stepf, ok := s.checkRangeTypeIsFloat(startStr, endStr, stepStr); ok { // is float
		step = stepf

		if stepStr == "r" {
			rand = true
		}

		precision1, step1 := valueGen.GetPrecision(start, step)
		precision2, step2 := valueGen.GetPrecision(end, step)
		if precision1 < precision2 {
			precision = precision2
			step = step2
		} else {
			precision = precision1
			step = step1
		}

		if (start > end && stepf > 0) || (start < end && stepf < 0) {
			step = -1 * stepf
		}

		dataType = "float"
		count = int(math.Floor(math.Abs(start-end)/stepf + 0.5))
		return

	} else if len(startStr) == 1 && len(endStr) == 1 { // is char
		step = 1

		if stepStr != "r" {
			stepChar, errChar3 := strconv.Atoi(stepStr)
			if errChar3 == nil {
				step = stepChar
			}
		} else if stepStr == "r" {
			rand = true
		}

		if (strings.Compare(startStr, endStr) > 0 && step.(int) > 0) ||
			(strings.Compare(startStr, endStr) < 0 && step.(int) < 0) {
			step = -1 * step.(int)
		}

		dataType = "char"

		startChar := startStr[0]
		endChar := endStr[0]

		gap := 0
		if endChar > startChar { // 负数转uint单独处理
			gap = int(endChar - startChar)
		} else {
			gap = int(startChar - endChar)
		}
		count = gap / step.(int)
		if count < 0 {
			count = count * -1
		}
		return
	}

	// else is string
	dataType = "string"
	step = 1
	return
}

func (s *RangeService) checkRangeTypeIsInt(startStr string, endStr string, stepStr string) (
	start int64, end int64, step int64, ok bool) {
	step = 1

	stepStr = strings.ToLower(strings.TrimSpace(stepStr))

	start, errInt1 := strconv.ParseInt(startStr, 0, 64)
	end, errInt2 := strconv.ParseInt(endStr, 0, 64)
	var errInt3 error

	if stepStr != "" && stepStr != "r" {
		step, errInt3 = strconv.ParseInt(stepStr, 0, 64)
	}

	if errInt1 == nil && errInt2 == nil && errInt3 == nil { // is int
		if (start > end && step > 0) || (start < end && step < 0) {
			step = -1 * step
		}

		ok = true
		return

	}

	return
}

func (s *RangeService) checkRangeTypeIsFloat(startStr string, endStr string, stepStr string) (
	start float64, end float64, step float64, ok bool) {

	stepStr = strings.ToLower(strings.TrimSpace(stepStr))

	start, errFloat1 := strconv.ParseFloat(startStr, 64)
	end, errFloat2 := strconv.ParseFloat(endStr, 64)
	var errFloat3 error

	if stepStr != "" && stepStr != "r" {
		step, errFloat3 = strconv.ParseFloat(stepStr, 64)
	} else {
		step = 0
	}

	if errFloat1 == nil && errFloat2 == nil && errFloat3 == nil { // is float
		if (start > end && step > 0) || (start < end && step < 0) {
			step = -1 * step
		}

		ok = true
		return
	}

	return
}
