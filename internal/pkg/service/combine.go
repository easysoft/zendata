package service

import (
	constant "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	stringUtils "github.com/easysoft/zendata/pkg/utils/string"
	"github.com/easysoft/zendata/pkg/utils/vari"
)

type CombineService struct {
	ExpressionService *ExpressionService `inject:""`
	LoopService       *LoopService       `inject:""`
}

func (s *CombineService) CombineChildrenIfNeeded(field *model.DefField) {
	if len(field.Fields) == 0 || !field.Union {
		return
	}

	fieldNameToValuesMap := map[string][]interface{}{}
	fieldNameToFieldMap := map[string]model.DefField{}

	// 1. get values for child fields
	for _, child := range field.Fields {
		fieldNameToValuesMap[child.Field] = child.Values
		fieldNameToFieldMap[child.Field] = child
	}

	// 2. deal with expression
	arrOfArr := make([][]interface{}, 0) // 2 dimension arr for child, [ [a,b,c], [1,2,3] ]
	for i, child := range field.Fields {
		childValues := fieldNameToValuesMap[child.Field]

		if child.Value != "" {
			childValues = s.ExpressionService.GenExpressionValues(child, fieldNameToValuesMap, fieldNameToFieldMap)
		}

		arrOfArr = append(arrOfArr, childValues)

		// clear child values after combined
		field.Fields[i].Values = nil
	}

	// 3. get combined values for parent field
	isRecursive := vari.Recursive
	if stringUtils.InArray(field.Mode, constant.Modes) { // set on field level
		isRecursive = field.Mode == constant.ModeRecursive || field.Mode == constant.ModeRecursiveShort
	}

	field.Values = s.combineChildrenValues(arrOfArr, isRecursive)
	s.LoopService.LoopFieldValues(field, true)
}

func (s *CombineService) combineChildrenValues(arrOfArr [][]interface{}, isRecursive bool) (ret []interface{}) {
	valueArr := s.populateRowsFromTwoDimArr(arrOfArr, isRecursive, false)

	for _, arr := range valueArr {
		line := ""
		for _, item := range arr {
			line += item.(string)
		}

		ret = append(ret, line)
	}
	return
}

func (s *CombineService) populateRowsFromTwoDimArr(arrOfArr [][]interface{}, isRecursive, isOnTopLevel bool) (
	values [][]interface{}) {
	count := vari.GenVars.Total

	if !isOnTopLevel {
		if isRecursive {
			count = s.getRecordCountForRecursive(arrOfArr)
		} else {
			count = s.getRecordCountForParallel(arrOfArr)
		}
	}

	indexArr := make([]int, 0)
	if isRecursive {
		indexArr = s.getModArr(arrOfArr)
	}

	for i := 0; i < count; i++ {
		strArr := make([]interface{}, 0)
		for j := 0; j < len(arrOfArr); j++ {
			child := arrOfArr[j]

			var index int
			if isRecursive {
				mod := indexArr[j]
				index = i / mod % len(child)
			} else {
				index = i % len(child)
			}

			val := child[index]
			strArr = append(strArr, val)
		}

		values = append(values, strArr)
	}

	return
}

func (s *CombineService) getRecordCountForParallel(arrOfArr [][]interface{}) int {
	// get max count of 2nd dim arr
	count := 1
	for _, arr := range arrOfArr {
		if count < len(arr) {
			count = len(arr)
		}
	}

	if count > vari.GenVars.Total {
		count = vari.GenVars.Total
	}

	return count
}

func (s *CombineService) getRecordCountForRecursive(arrOfArr [][]interface{}) int {
	count := 1
	for i := 0; i < len(arrOfArr); i++ {
		arr := arrOfArr[i]
		count = len(arr) * count
	}
	return count
}

func (s *CombineService) getModArr(arrOfArr [][]interface{}) []int {
	indexArr := make([]int, 0)
	for _, _ = range arrOfArr {
		indexArr = append(indexArr, 0)
	}

	for i := 0; i < len(arrOfArr); i++ {
		loop := 1
		for j := i + 1; j < len(arrOfArr); j++ {
			loop = loop * len(arrOfArr[j])
		}

		indexArr[i] = loop
	}

	return indexArr
}
