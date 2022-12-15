package service

import (
	"fmt"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"regexp"
	"strconv"
	"strings"
)

type FieldService struct {
	ResService     *ResService     `inject:""`
	TextService    *TextService    `inject:""`
	ValueService   *ValueService   `inject:""`
	ArticleService *ArticleService `inject:""`

	FixService    *FixService    `inject:""`
	LoopService   *LoopService   `inject:""`
	ListService   *ListService   `inject:""`
	RangeService  *RangeService  `inject:""`
	RandomService *RandomService `inject:""`
}

func (s *FieldService) Generate(field *model.DefField, parentIsJoin bool) {
	field.Join = field.Join || parentIsJoin

	s.RangeService.DealwithFixRange(field)

	// iterate children
	if len(field.Fields) > 0 {
		for i, _ := range field.Fields {
			if field.Fields[i].From == "" {
				field.Fields[i].From = field.From
			}
			field.Fields[i].FileDir = field.FileDir

			s.Generate(&field.Fields[i], field.Join)
		}
		return
	}

	// generate values
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
		field.Values = s.RandomService.RandomValues(field.Values)
	}

	if field.Use != "" && field.From == "" {
		field.From = vari.GlobalVars.DefData.From
	}
}

func (s *FieldService) GenerateValuesForNoReferField(field *model.DefField) {
	s.CreateField(field)

	s.LoopService.ComputerLoopTimes(field) // change LoopStart, LoopEnd for conf like loop:  1-10             # 循环1次，2次……

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
		if count >= vari.GlobalVars.Total || count >= uniqueTotal || isRandomAndLoopEnd {
			for _, v := range field.Values {
				str := fmt.Sprintf("%v", v)
				str = s.FixService.AddFix(str, field, count, true)
				values = append(values, str)
			}
			break
		}

		// 处理格式、前后缀、loop等
		val := s.LoopService.LoopFieldValueToSingleStr(field, &indexOfRow, count, true)
		values = append(values, val)

		count++

		if count >= vari.GlobalVars.Total || count >= uniqueTotal {
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
		s.ListService.CreateListFieldValues(field)
	} else if field.Type == consts.FieldTypeArticle {
		s.ArticleService.CreateArticleField(field)

	} else if field.Type == consts.FieldTypeTimestamp {
		s.ValueService.CreateTimestampField(field)
	} else if field.Type == consts.FieldTypeUlid {
		s.ValueService.CreateUlidField(field)
	}

	return
}

func (s *FieldService) GenValuesForConfig(field *model.DefField) (values []interface{}) {
	groupValues := vari.GlobalVars.ResData[field.Config]

	field.Values = groupValues["all"]

	s.LoopService.LoopFieldValues(field, true)

	return
}

func (s *FieldService) GenValuesForSingleRes(field *model.DefField) {
	if field.Use != "" { // refer to ranges or instance
		groupValues := vari.GlobalVars.ResData[s.ResService.getFromKey(field)]

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
		groupValues := vari.GlobalVars.ResData[s.ResService.getFromKey(field)]
		resKey := field.Select

		// deal with the key
		if vari.GlobalVars.DefData.Type == consts.DefTypeArticle {
			resKey = resKey + "_" + field.Field
		}

		field.Values = append(field.Values, groupValues[resKey]...)
	}

	s.LoopService.LoopFieldValues(field, true)

	return
}

func (s *FieldService) GenValuesForMultiRes(field *model.DefField, withFix bool) {
	unionValues := make([]interface{}, 0) // 2 dimension arr for from, [ [a,b,c], [1,2,3] ]

	for _, from := range field.Froms {
		if from.From == "" {
			from.From = field.From
		}

		from.FileDir = field.FileDir

		s.Generate(&from, true)

		unionValues = append(unionValues, from.Values...)
	}

	count := len(unionValues)
	if count > vari.GlobalVars.Total {
		count = vari.GlobalVars.Total
	}

	field.Values = unionValues

	s.LoopService.LoopFieldValues(field, true)

	return
}

func (s *FieldService) getRepeatValuesFromAll(groupValues map[string][]interface{}, numLimit, repeat int) (ret []interface{}) {
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

func (s *FieldService) getRepeatValuesFromGroups(groupValues map[string][]interface{}, info []retsInfo) (ret []interface{}) {
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
			for i := 0; (v.numLimit > 0 && i < v.numLimit) && i < len(arr) && i < vari.GlobalVars.Total; i++ {
				index := i / v.repeat
				ret = append(ret, arr[index])
				count++
			}
		} else { // privateA
			for i := 0; i < len(arr) && i < vari.GlobalVars.Total; i++ {
				index := i / v.repeat % len(arr)
				ret = append(ret, arr[index])
				count++
			}
		}

		if count >= vari.GlobalVars.Total {
			break exit
		}

	}
	return
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
