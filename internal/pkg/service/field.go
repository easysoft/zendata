package service

import (
	"errors"
	"fmt"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/gen"
	"github.com/easysoft/zendata/internal/pkg/gen/helper"
	valueGen "github.com/easysoft/zendata/internal/pkg/gen/value"
	"github.com/easysoft/zendata/internal/pkg/model"
	commonUtils "github.com/easysoft/zendata/pkg/utils/common"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	stringUtils "github.com/easysoft/zendata/pkg/utils/string"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type FieldService struct {
	ResService     *ResService
	TextService    *TextService
	ValueService   *ValueService
	ArticleService *ArticleService
}

func NewFieldService() *FieldService {
	return &FieldService{
		ResService:   NewResService(),
		TextService:  NewTextService(),
		ValueService: NewValueService(),
	}
}

func (s *FieldService) Generate(field *model.DefField) {
	gen.DealwithFixRange(field)

	// has children
	if len(field.Fields) > 0 {
		for i, _ := range field.Fields {
			s.Generate(&field.Fields[i])
		}
		return
	}

	//
	if len(field.Froms) > 0 { // refer to multi res
		s.GenValuesForMultiRes(field, true)

	} else if field.From != "" && field.Type != consts.FieldTypeArticle { // refer to res
		s.GenValuesForSingleRes(field)

	} else if field.Config != "" { // refer to config
		s.GenValuesForConfig(field)

	} else { // not a refer
		s.GenerateValuesForNoReferField(field)
	}

	// random values
	if field.Rand && field.Type != consts.FieldTypeArticle {
		field.Values = gen.RandomValues(field.Values)
	}

	if field.Use != "" && field.From == "" {
		field.From = vari.GenVars.DefData.From
	}
}

func (s *FieldService) GenerateValuesForNoReferField(field *model.DefField) {
	s.CreateField(field)

	s.computerLoopTimes(field) // change LoopStart, LoopEnd for conf like loop:  1-10             # 循环1次，2次……

	uniqueTotal := s.computerUniqueTotal(field) // computer total for conf like prefix: 1-3, postfix: 1-3

	indexOfRow := 0
	count := 0
	values := make([]interface{}, 0)

	for {
		// 2. random replacement
		isRandomAndLoopEnd := !vari.ResLoading && //  ignore rand in resource
			!(*field).ReferToAnotherYaml &&
			(*field).IsRand && (*field).LoopIndex > (*field).LoopEnd
		// isNotRandomAndValOver := !(*field).IsRand && indexOfRow >= len(fieldWithValues.Values)
		if count >= vari.GenVars.Total || count >= uniqueTotal || isRandomAndLoopEnd {
			for _, v := range field.Values {
				str := fmt.Sprintf("%v", v)
				str = s.addFix(str, field, count, true)
				values = append(values, str)
			}
			break
		}

		// 处理格式、前后缀、loop等
		val := s.loopFieldValWithFix(field, &indexOfRow, count, true)
		values = append(values, val)

		count++

		if count >= vari.GenVars.Total || count >= uniqueTotal {
			break
		}

		(*field).LoopIndex = (*field).LoopIndex + 1
		if (*field).LoopIndex > (*field).LoopEnd {
			(*field).LoopIndex = (*field).LoopStart
		}
	}

	field.Values = values

	return
}

func (s *FieldService) CreateField(field *model.DefField) {
	if field.Type == "" { // set default
		field.Type = consts.FieldTypeList
	}
	if field.Length > 0 {
		field.Length = field.Length - len(field.Prefix) - len(field.Postfix)
		if field.Length < 0 {
			field.Length = 0
		}
	}

	if field.Type == consts.FieldTypeList {
		s.CreateListFieldValues(field)
	} else if field.Type == consts.FieldTypeArticle {
		s.ArticleService.CreateArticleField(field)

	} else if field.Type == consts.FieldTypeTimestamp {
		s.ValueService.CreateTimestampField(field)
	} else if field.Type == consts.FieldTypeUlid {
		s.ValueService.CreateUlidField(field)
	}

	return
}

func (s *FieldService) CreateListFieldValues(field *model.DefField) {
	if strings.Index(field.Range, ".txt") > -1 {
		s.TextService.CreateFieldValuesFromText(field)
	} else {
		s.CreateFieldValuesFromRange(field)
	}
}

func (s *FieldService) computerLoopTimes(field *model.DefField) {
	if (*field).LoopIndex != 0 {
		return
	}

	arr := strings.Split(field.Loop, "-")
	(*field).LoopStart, _ = strconv.Atoi(arr[0])
	if len(arr) > 1 {
		field.LoopEnd, _ = strconv.Atoi(arr[1])
	}

	if (*field).LoopStart == 0 {
		(*field).LoopStart = 1
	}
	if (*field).LoopEnd == 0 {
		(*field).LoopEnd = 1
	}

	(*field).LoopIndex = (*field).LoopStart
}

func (s *FieldService) CreateFieldValuesFromRange(field *model.DefField) {
	rang := field.Range

	// gen empty values
	if rang == "" {
		for i := 0; i < vari.GenVars.Total; i++ {
			field.Values = append(field.Values, "")
			if strings.Index(field.Format, "uuid") == -1 {
				break
			}
		}

		return
	}

	// gen from field's range
	rangeSections := gen.ParseRangeProperty(rang) // parse 1

	index := 0
	for _, rangeSection := range rangeSections {
		if index >= consts.MaxNumb {
			break
		}
		if rangeSection == "" {
			continue
		}

		descStr, stepStr, count, countTag := gen.ParseRangeSection(rangeSection) // parse 2
		if strings.ToLower(stepStr) == "r" {
			(*field).IsRand = true
		}

		typ, desc := gen.ParseRangeSectionDesc(descStr) // parse 3

		items := make([]interface{}, 0)
		if typ == "literal" {
			items = s.CreateValuesFromLiteral(field, desc, stepStr, count, countTag)
		} else if typ == "interval" {
			items = s.CreateValuesFromInterval(field, desc, stepStr, count, countTag)
		} else if typ == "yaml" { // load from a yaml
			items = s.CreateValuesFromYaml(field, desc, stepStr, count, countTag)
			field.ReferToAnotherYaml = true
		}

		field.Values = append(field.Values, items...)
		index = index + len(items)
	}

	if len(field.Values) == 0 {
		field.Values = append(field.Values, "N/A")
	}
}

func (s *FieldService) CreateFieldFixValuesFromList(strRang string, field *model.DefField) (rang *model.Range) {
	rang = &model.Range{}

	if strRang == "" {
		return
	}

	rangeSections := gen.ParseRangeProperty(strRang) // parse 1

	index := 0
	for _, rangeSection := range rangeSections {
		if index >= consts.MaxNumb {
			break
		}
		if rangeSection == "" {
			continue
		}

		descStr, stepStr, count, countTag := gen.ParseRangeSection(rangeSection) // parse 2
		if strings.ToLower(stepStr) == "r" {
			rang.IsRand = true
		}

		typ, desc := gen.ParseRangeSectionDesc(descStr) // parse 3

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

func (s *FieldService) CreateValuesFromLiteral(field *model.DefField, desc string, stepStr string, repeat int, repeatTag string) (items []interface{}) {
	elemArr := gen.ParseDesc(desc)
	step, _ := strconv.Atoi(stepStr)
	if step == 0 {
		step = 1
	}
	total := 0

	if field.Path != "" && stepStr == "r" {
		pth := field.Path
		key := helper.GetRandFieldSection(pth)

		items = append(items, gen.Placeholder(key))
		mp := gen.PlaceholderMapForRandValues("list", elemArr, "", "", "", "",
			field.Format, repeat, repeatTag)

		vari.RandFieldSectionShortKeysToPathMap[key] = pth
		vari.RandFieldSectionPathToValuesMap[key] = mp
		return
	}

	if repeatTag == "" {
		for i := 0; i < len(elemArr); {
			idx := i
			if field.Path == "" && stepStr == "r" {
				idx = commonUtils.RandNum(len(elemArr)) // should set random here too
			}

			val := elemArr[idx]
			total = gen.AppendValues(&items, val, repeat, total)

			if total >= consts.MaxNumb {
				break
			}
			i += step
		}
	} else if repeatTag == "!" {
		isRand := field.Path == "" && stepStr == "r"
		for i := 0; i < repeat; {
			total = gen.AppendArrItems(&items, elemArr, total, isRand)

			if total >= consts.MaxNumb {
				break
			}
			i += step
		}
	}

	if field.Path == "" && stepStr == "r" { // for ranges and instances, random
		items = gen.RandomInterfaces(items)
	}

	return
}

func (s *FieldService) CreateValuesFromInterval(field *model.DefField, desc, stepStr string, repeat int, repeatTag string) (items []interface{}) {
	elemArr := strings.Split(desc, "-")
	startStr := elemArr[0]
	endStr := startStr
	if len(elemArr) > 1 {
		endStr = elemArr[1]
	}

	dataType, step, precision, rand, _ := gen.CheckRangeType(startStr, endStr, stepStr)
	field.Precision = precision

	// 1. random replacement
	if field.Path != "" && dataType != "string" && rand { // random. for res, field.Path == ""
		pth := field.Path + "->" + desc
		key := helper.GetRandFieldSection(pth)

		val := gen.Placeholder(key)
		strItems := make([]string, 0)

		//for i := 0; i < repeat*count; i++ { // chang to add only one placeholder item
		items = append(items, val)
		strItems = append(strItems, val)
		//}

		mp := gen.PlaceholderMapForRandValues(dataType, strItems, startStr, endStr, fmt.Sprintf("%v", step),
			strconv.Itoa(precision), field.Format, repeat, repeatTag)

		vari.RandFieldSectionShortKeysToPathMap[key] = pth
		vari.RandFieldSectionPathToValuesMap[key] = mp

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
		items = gen.RandomInterfaces(items)
	}

	return
}

func (s *FieldService) CreateValuesFromYaml(field *model.DefField, yamlFile, stepStr string, repeat int, repeatTag string) (items []interface{}) {
	// keep root def, since vari.ZdDef will be overwrite by refer yaml file
	rootDef := vari.GenVars.DefData
	configDir := vari.GenVars.ConfigFileDir
	res := vari.GenVars.ResData

	configFile := fileUtils.ComputerReferFilePath(yamlFile, field)
	fieldsToExport := make([]string, 0) // set to empty to use all fields
	rows, colIsNumArr, _ := gen.GenerateFromYaml([]string{configFile}, &fieldsToExport)
	if field.Rand {
		rows = gen.RandomValuesArr(rows)
	}

	items = gen.PrintLines(rows, consts.FormatData, "", colIsNumArr, fieldsToExport)

	if repeat > 0 {
		if repeat > len(items) {
			repeat = len(items)
		}
		items = items[:repeat]
	}

	// rollback root def when finish to deal with refer yaml file
	vari.GenVars.DefData = rootDef
	vari.GenVars.ConfigFileDir = configDir
	vari.GenVars.ResData = res

	return
}

func (s *FieldService) GenValuesForConfig(field *model.DefField) (values []interface{}) {
	groupValues := vari.GenVars.ResData[field.Config]

	field.Values = groupValues["all"]

	s.loopFieldValues(field, true)

	return
}

func (s *FieldService) GenValuesForSingleRes(field *model.DefField) {
	if field.Use != "" { // refer to ranges or instance
		groupValues := vari.GenVars.ResData[s.ResService.getFromKey(field)]

		uses := strings.TrimSpace(field.Use) // like group{limit:repeat}
		use, numLimit, repeat := s.getNum(uses)
		if strings.Index(use, "all") == 0 {
			valuesForAdd := s.getRepeatValuesFromAll(groupValues, numLimit, repeat)
			field.Values = append(field.Values, valuesForAdd...)
		} else {
			infos := s.parseUse(uses)
			valuesForAdd := s.getRepeatValuesFromGroups(groupValues, infos)
			field.Values = append(field.Values, valuesForAdd...)
		}
	} else if field.Select != "" { // refer to excel
		groupValues := vari.GenVars.ResData[s.ResService.getFromKey(field)]
		resKey := field.Select

		// deal with the key
		if vari.GenVars.DefData.Type == consts.DefTypeArticle {
			resKey = resKey + "_" + field.Field
		}

		field.Values = append(field.Values, groupValues[resKey]...)
	}

	s.loopFieldValues(field, true)

	return
}

func (s *FieldService) GenValuesForMultiRes(field *model.DefField, withFix bool) {
	unionValues := make([]interface{}, 0) // 2 dimension arr for child, [ [a,b,c], [1,2,3] ]

	for _, child := range field.Froms {
		if child.From == "" {
			child.From = field.From
		}

		child.FileDir = field.FileDir
		childValues := s.GenerateForFieldRecursive(&child, withFix, vari.GenVars.Total)
		unionValues = append(unionValues, childValues...)
	}

	count := len(unionValues)
	if count > vari.GenVars.Total {
		count = vari.GenVars.Total
	}

	field.Values = unionValues

	s.loopFieldValues(field, true)

	return
}

func (s *FieldService) addFix(str string, field *model.DefField, count int, withFix bool) (ret string) {
	prefix := s.GetStrValueFromRange(field.PrefixRange, count)
	postfix := s.GetStrValueFromRange(field.PostfixRange, count)
	divider := field.Divider

	if field.Length > runewidth.StringWidth(str) {
		str = stringUtils.AddPad(str, *field)
	}
	if withFix && !vari.Trim {
		str = prefix + str + postfix
	}
	if vari.Format == consts.FormatText && !vari.Trim {
		str += divider
	}

	ret = str
	return
}

func (s *FieldService) loopFieldValues(field *model.DefField, withFix bool) {
	s.computerLoopTimes(field)

	values := make([]interface{}, 0)

	indexOfRow := 0
	count := 0
	for {
		// 处理格式、前后缀、loop等
		str := s.loopFieldValWithFix(field, &indexOfRow, count, withFix)
		values = append(values, str)

		count++
		isRandomAndLoopEnd := (*field).IsRand && (*field).LoopIndex == (*field).LoopEnd
		isNotRandomAndValOver := !(*field).IsRand && indexOfRow >= len(values)
		if count >= vari.GenVars.Total || isRandomAndLoopEnd || isNotRandomAndValOver {
			break
		}

		(*field).LoopIndex = (*field).LoopIndex + 1
		if (*field).LoopIndex > (*field).LoopEnd {
			(*field).LoopIndex = (*field).LoopStart
		}
	}

	return
}

func (s *FieldService) loopFieldValWithFix(field *model.DefField, indexOfRow *int, count int, withFix bool) (loopStr string) {

	for j := 0; j < (*field).LoopIndex; j++ {
		if loopStr != "" {
			loopStr = loopStr + field.Loopfix
		}

		str, err := s.GetFieldVal(*field, indexOfRow)
		if err != nil {
			str = "N/A"
		}
		loopStr = loopStr + str

		*indexOfRow++
	}

	loopStr = s.addFix(loopStr, field, count, withFix)

	return
}

func (s *FieldService) GetStrValueFromRange(rang *model.Range, index int) string {
	if rang == nil || len(rang.Values) == 0 {
		return ""
	}

	idx := index % len(rang.Values)
	x := rang.Values[idx]
	return s.convPrefixVal2Str(x, "")
}

func (s *FieldService) convPrefixVal2Str(val interface{}, format string) string {
	str := "n/a"
	success := false

	switch val.(type) {
	case int64:
		if format != "" {
			str, success = stringUtils.FormatStr(format, val.(int64), 0)
		}
		if !success {
			str = strconv.FormatInt(val.(int64), 10)
		}
	case float64:
		precision := 0
		if format != "" {
			str, success = stringUtils.FormatStr(format, val.(float64), precision)
		}
		if !success {
			str = strconv.FormatFloat(val.(float64), 'f', precision, 64)
		}
	case byte:
		str = string(val.(byte))
		if format != "" {
			str, success = stringUtils.FormatStr(format, str, 0)
		}
		if !success {
			str = string(val.(byte))
		}
	case string:
		str = val.(string)

		match, _ := regexp.MatchString("%[0-9]*d", format)
		if match {
			valInt, err := strconv.Atoi(str)
			if err == nil {
				str, success = stringUtils.FormatStr(format, valInt, 0)
			}
		} else {
			str, success = stringUtils.FormatStr(format, str, 0)
		}
	default:
	}

	return str
}

func (s *FieldService) GetFieldVal(field model.DefField, index *int) (val string, err error) {
	// 叶节点
	if len(field.Values) == 0 {
		if helper.SelectExcelWithExpr(field) {
			logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("fail_to_generate_field", field.Field), color.FgCyan)
			err = errors.New("")
		}
		return
	}

	idx := *index % len(field.Values)
	str := field.Values[idx]
	val = s.GetFieldValStr(field, str)

	return
}

func (s *FieldService) GetFieldValStr(field model.DefField, val interface{}) string {
	str := "n/a"
	success := false

	format := strings.TrimSpace(field.Format)

	if field.Type == consts.FieldTypeTimestamp && field.Format != "" {
		str = time.Unix(val.(int64), 0).Format(field.Format)
		return str
	}

	switch val.(type) {
	case int64:
		if format != "" {
			str, success = stringUtils.FormatStr(format, val.(int64), 0)
		}
		if !success {
			str = strconv.FormatInt(val.(int64), 10)
		}
	case float64:
		precision := 0
		if field.Precision > 0 {
			precision = field.Precision
		}
		if format != "" {
			str, success = stringUtils.FormatStr(format, val.(float64), precision)
		}
		if !success {
			str = strconv.FormatFloat(val.(float64), 'f', precision, 64)
		}
	case byte:
		str = string(val.(byte))
		if format != "" {
			str, success = stringUtils.FormatStr(format, str, 0)
		}
		if !success {
			str = string(val.(byte))
		}
	case string:
		str = val.(string)

		match, _ := regexp.MatchString("%[0-9]*d", format)
		if match {
			valInt, err := strconv.Atoi(str)
			if err == nil {
				str, success = stringUtils.FormatStr(format, valInt, 0)
			}
		} else {
			str, success = stringUtils.FormatStr(format, str, 0)
		}
	default:
	}

	return str
}

func (s *FieldService) computerUniqueTotal(field *model.DefField) (ret int) {
	ret = len(field.Values)

	if field.PostfixRange != nil && len(field.PostfixRange.Values) > 0 {
		ret *= len(field.PostfixRange.Values)
	}

	if field.PrefixRange != nil && len(field.PrefixRange.Values) > 0 {
		ret *= len(field.PrefixRange.Values)
	}

	return
}

func (s *FieldService) getRepeatValuesFromAll(groupValues map[string][]string, numLimit, repeat int) (ret []string) {
	if repeat == 0 {
		repeat = 1
	}

	count := 0
exit:
	for _, arr := range groupValues {
		for _, item := range arr {
			for i := 0; i < repeat; i++ {
				ret = append(ret, item)
				count++

				if numLimit > 0 && count >= numLimit {
					break exit
				}
			}
		}
	}

	return
}

func (s *FieldService) getRepeatValuesFromGroups(groupValues map[string][]string, info []retsInfo) (ret []string) {
	count := 0

exit:
	for _, v := range info {
		if v.repeat == 0 {
			v.repeat = 1
		}

		arr := groupValues[v.ret]
		if len(arr) == 0 {
			break exit
		}
		if v.numLimit != 0 { // privateB{n}
			for i := 0; (v.numLimit > 0 && i < v.numLimit) && i < len(arr) && i < vari.GenVars.Total; i++ {
				index := i / v.repeat
				ret = append(ret, arr[index])
				count++
			}
		} else { // privateA
			for i := 0; i < len(arr) && i < vari.GenVars.Total; i++ {
				index := i / v.repeat % len(arr)
				ret = append(ret, arr[index])
				count++
			}
		}

		if count >= vari.GenVars.Total {
			break exit
		}

	}
	return
}

// pars Uses
type retsInfo struct {
	ret      string
	numLimit int
	repeat   int
}

func (s *FieldService) parseUse(groups string) (results []retsInfo) {
	rets := strings.Split(groups, ",")
	results = make([]retsInfo, len(rets))
	regx := regexp.MustCompile(`\{([^:]*):?([^:]*)\}`)
	for k, v := range rets {
		results[k].ret = regx.ReplaceAllString(v, "")
		arr := regx.FindStringSubmatch(v)
		if len(arr) >= 2 {
			results[k].numLimit, _ = strconv.Atoi(arr[1])
		}
		if len(arr) >= 3 {
			results[k].repeat, _ = strconv.Atoi(arr[2])
			if results[k].repeat == 0 {
				results[k].repeat = 1
			}
		}
	}
	return
}

func (s *FieldService) getNum(group string) (ret string, numLimit, repeat int) {
	regx := regexp.MustCompile(`\{([^:]*):?([^:]*)\}`)
	arr := regx.FindStringSubmatch(group)
	if len(arr) >= 2 {
		numLimit, _ = strconv.Atoi(arr[1])
	}
	if len(arr) >= 3 {
		repeat, _ = strconv.Atoi(arr[2])
	}

	ret = regx.ReplaceAllString(group, "")

	return
}
