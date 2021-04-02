package gen

import (
	"fmt"
	"github.com/easysoft/zendata/src/model"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	"github.com/easysoft/zendata/src/utils/vari"
	"regexp"
	"strings"
)

func CreateArticleField(field *model.DefField, fieldWithValue *model.FieldWithValues) {
	fieldWithValue.Field = field.Field

	content := genArticleContent(field) + "\n"

	fieldWithValue.Values = append(fieldWithValue.Values, content)
}

func genArticleContent(field *model.DefField) (ret string) {

	numMap, nameMap, contentWithOutDefine := getNumMap(field.Range)

	resFile, resType, sheet := fileUtils.GetResProp(field.From, "")
	dataMap := getDataMap(numMap, nameMap, field, resFile, resType, sheet)

	ret = genArticle(contentWithOutDefine, dataMap, nameMap)

	return
}

func genArticle(content string, dataMap map[string][]string, nameMap map[string]string) (ret string) {
	ret = content

	for key, arr := range dataMap {
		for _, item := range arr {
			placeholder := fmt.Sprintf("{%s}", key)
			ret = strings.Replace(ret, placeholder, item, 1)
		}
	}

	return
}

func getDataMap(numMap map[string]int, nameMap map[string]string, field *model.DefField,
	resFile string, resType string, sheet string) (ret map[string][]string) {
	ret = map[string][]string{}

	field.Rand = true
	for key, num := range numMap {
		originTotal := vari.Total
		vari.Total = num

		slct, ok := nameMap[key]
		if ok {
			field.Select = slct
		} else {
			field.Select = key
		}

		valueMap, _ := getResValue(resFile, resType, sheet, field)
		ret[key] = valueMap[field.Select]

		vari.Total = originTotal // rollback
	}

	return
}

func getNumMap(content string) (numMap map[string]int, nameMap map[string]string, contentWithOutDefine string) {
	numMap = map[string]int{}
	nameMap = map[string]string{}
	arrWithOutDefine := make([]string, 0)

	content = strings.Trim(content, "`")
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.Index(line, "#") == 0 {
			line = strings.TrimLeft(line, "#")
			arr := strings.Split(line, "=")
			if len(arr) < 2 {
				continue
			}

			leftArr := strings.Split(arr[0], " ")
			vari := leftArr[len(leftArr)-1]
			expr := strings.Split(arr[1], " ")[0]

			nameMap[vari] = expr
			continue
		}

		arrWithOutDefine = append(arrWithOutDefine, line)

		regx := regexp.MustCompile(`\{((?U).*)\}`)
		arr := regx.FindAllStringSubmatch(line, -1)

		for _, child := range arr {
			name := child[1]
			i, ok := numMap[name]
			if !ok {
				numMap[name] = 1
			} else {
				numMap[name] = i + 1
			}
		}
	}

	contentWithOutDefine = strings.Join(arrWithOutDefine, "\n")

	return
}
