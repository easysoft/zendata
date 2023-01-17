package gen

import (
	"errors"
	"fmt"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	genHelper "github.com/easysoft/zendata/internal/pkg/gen/helper"
	"github.com/easysoft/zendata/internal/pkg/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	stringUtils "github.com/easysoft/zendata/pkg/utils/string"
	"github.com/easysoft/zendata/pkg/utils/vari"
)

func GenerateFromContent(fileContents [][]byte, fieldsToExport *[]string) (
	rows [][]string, colIsNumArr []bool, err error) {

	vari.GlobalVars.DefData = LoadDataContentDef(fileContents, fieldsToExport)

	if err = CheckParams(); err != nil {
		return
	}

	FixTotalNum()
	genResData(fieldsToExport)

	topLevelFieldNameToValuesMap := genFieldsData(fieldsToExport, &colIsNumArr, vari.GlobalVars.Total)
	twoDimArr := genDataTwoDimArr(topLevelFieldNameToValuesMap, fieldsToExport, vari.GlobalVars.Total)
	rows = populateRowsFromTwoDimArr(twoDimArr, vari.GlobalVars.Recursive, true, vari.GlobalVars.Total)

	return
}

func GenerateFromYaml(files []string, fieldsToExport *[]string) (
	rows [][]string, colIsNumArr []bool, err error) {

	vari.GlobalVars.ConfigFileDir = fileUtils.GetAbsDir(files[0])

	contents := LoadFilesContents(files)
	rows, colIsNumArr, err = GenerateFromContent(contents, fieldsToExport)

	return
}

func GenerateForFieldRecursive(field *model.DefField, withFix bool, total int) (values []string) {
	DealwithFixRange(field)

	if len(field.Fields) > 0 { // has child fields
		values = genValuesForChildFields(field, withFix, total)

	} else if len(field.Froms) > 0 { // refer to multi res
		values = GenValuesForMultiRes(field, withFix, total)

	} else if field.From != "" && field.Type != consts.FieldTypeArticle { // refer to res
		values = GenValuesForSingleRes(field, total)

	} else if field.Config != "" { // refer to config
		values = GenValuesForConfig(field, total)

	} else { // leaf field
		values = GenerateValuesForField(field, total)
	}

	if field.Rand && field.Type != consts.FieldTypeArticle {
		values = RandomStrValues(values)
	}

	return values
}

func GenerateValuesForField(field *model.DefField, total int) []string {
	values := make([]string, 0)

	fieldWithValues := CreateField(field)

	computerLoop(field)
	indexOfRow := 0
	count := 0

	uniqueTotal := len(fieldWithValues.Values)

	if l := len(field.PostfixRange.Values); l > 0 {
		uniqueTotal *= l
	}
	if l := len(field.PrefixRange.Values); l > 0 {
		uniqueTotal *= l
	}

	for {
		// 2. random replacement
		isRandomAndLoopEnd := !vari.ResLoading && //  ignore rand in resource
			!(*field).ReferToAnotherYaml &&
			(*field).IsRand && (*field).LoopIndex > (*field).LoopEnd
		// isNotRandomAndValOver := !(*field).IsRand && indexOfRow >= len(fieldWithValues.Values)
		if count >= total || count >= uniqueTotal || isRandomAndLoopEnd {
			for _, v := range fieldWithValues.Values {
				str := fmt.Sprintf("%v", v)
				str = addFix(str, field, count, true)
				values = append(values, str)
			}
			break
		}

		// 处理格式、前后缀、loop等
		val := loopFieldValWithFix(field, fieldWithValues, &indexOfRow, count, true)
		values = append(values, val)

		count++

		if count >= total || count >= uniqueTotal {
			break
		}

		(*field).LoopIndex = (*field).LoopIndex + 1
		if (*field).LoopIndex > (*field).LoopEnd {
			(*field).LoopIndex = (*field).LoopStart
		}
	}

	return values
}

func CheckParams() (err error) {
	if len(vari.GlobalVars.DefData.Fields) == 0 {
		err = errors.New("")
	} else if vari.GlobalVars.DefData.Type == consts.DefTypeArticle && vari.GlobalVars.Output == "" { // gen article
		errMsg := i118Utils.I118Prt.Sprintf("gen_article_must_has_out_param")
		logUtils.PrintErrMsg(errMsg)
		err = errors.New(errMsg)
	}

	return
}

func FixTotalNum() {
	if vari.GlobalVars.Total < 0 {
		if vari.GlobalVars.DefData.Type == consts.DefTypeArticle {
			vari.GlobalVars.Total = 1
		} else {
			vari.GlobalVars.Total = consts.DefaultNumber
		}
	}
}

func genResData(fieldsToExport *[]string) {
	// 为被引用的资源生成数据
	vari.ResLoading = true // not to use placeholder when loading res
	vari.Res = LoadResDef(*fieldsToExport)
	vari.ResLoading = false
}

func genFieldsData(fieldsToExport *[]string, colIsNumArr *[]bool, total int) (topLevelFieldNameToValuesMap map[string][]string) {
	topLevelFieldNameToValuesMap = map[string][]string{}

	for index, field := range vari.GlobalVars.DefData.Fields {
		if !stringUtils.StrInArr(field.Field, *fieldsToExport) {
			continue
		}

		if field.Use != "" && field.From == "" {
			field.From = vari.GlobalVars.DefData.From
		}
		values := GenerateForFieldRecursive(&field, true, total)

		if index > len(vari.GlobalVars.DefData.Fields)-1 {
			logUtils.PrintLine("")
		}

		vari.GlobalVars.DefData.Fields[index].Precision = field.Precision

		topLevelFieldNameToValuesMap[field.Field] = values
		*colIsNumArr = append(*colIsNumArr, field.IsNumb)
	}

	return
}

func genDataTwoDimArr(topLevelFieldNameToValuesMap map[string][]string, fieldsToExport *[]string, total int) (
	arrOfArr [][]string) { // 2 dimension arr for child, [ [a,b,c], [1,2,3] ]

	for _, child := range vari.GlobalVars.DefData.Fields {
		if !stringUtils.StrInArr(child.Field, *fieldsToExport) {
			continue
		}

		childValues := topLevelFieldNameToValuesMap[child.Field]

		// is value expression
		if child.Value != "" {
			childValues = genHelper.GenExpressionValues(child, topLevelFieldNameToValuesMap, vari.GlobalVars.TopFieldMap)
		}

		// select from excel with expr
		if genHelper.IsSelectExcelWithExpr(child) {
			selects := genHelper.ReplaceVariableValues(child.Select, topLevelFieldNameToValuesMap)
			wheres := genHelper.ReplaceVariableValues(child.Where, topLevelFieldNameToValuesMap)

			childValues = make([]string, 0)
			childMapValues := make([][]string, 0)
			for index, slct := range selects {
				temp := child
				temp.Select = slct
				temp.Where = wheres[index%len(wheres)]

				resFile, _, sheet := fileUtils.GetResProp(temp.From, temp.FileDir)

				//	问题描述：
				//	原代码为：`selectCount := vari.Toal / len(selects)`
				//	因为整除的向下取整，如果`len(selects)`为3，`total`为8，则`selectCount`为2
				//	对于每一个`selects`的元素来说，都只会查两个元素，这样加起来一共只有6个结果，
				//	导致另外两个结果只能通过重复查到的数据的方式补充。
				//	解决方案：
				//  将代码改为: `selectCount := total / len(selects) + 1`,以达到使用人员的真正想要的
				//	即查到足够的数量，而不是通过重复补齐
				selectCount := total/len(selects) + 1
				mp := generateFieldValuesFromExcel(resFile, sheet, &temp, selectCount) // re-generate values
				for _, items := range mp {
					childMapValues = append(childMapValues, items)
				}
			}
			for index := 0; len(childValues) < total; {
				for i, _ := range selects {
					childValues = append(childValues, childMapValues[i][index%len(childMapValues[i])])
				}
				index++
			}
		}

		arrOfArr = append(arrOfArr, childValues)
	}

	return
}
