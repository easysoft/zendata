package service

import (
	"fmt"
	constant "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	commonUtils "github.com/easysoft/zendata/pkg/utils/common"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type ArticleService struct {
	ResService *ResService `inject:""`
}

func (s *ArticleService) CreateArticleField(field *model.DefField) {
	values := make([]interface{}, 0)

	numMap, nameMap, indexMap, contentWithoutComments := s.getNumMap(field.Range)
	resFile, resType, sheet := fileUtils.GetResProp(field.From, "")
	dataMap := s.getDataMap(numMap, nameMap, field, resFile, resType, sheet)

	for i := 0; i < vari.GlobalVars.Total; i++ {
		content := s.genArticle(contentWithoutComments, dataMap, indexMap) + "\n"
		values = append(values, content)
	}

	field.Values = values
}

func (s *ArticleService) genArticle(content string, dataMap map[string][]interface{},
	indexMap map[string]int) (ret string) {
	ret = content

	//for key, arr := range dataMap {
	//	for _, item := range arr {
	//		placeholder := fmt.Sprintf("{%s}", key)
	//		ret = strings.Replace(ret, placeholder, item, 1)
	//	}
	//}

	regx := regexp.MustCompile(`[\(\[\{]((?U).*)[\)\]\}]`)
	arr := regx.FindAllStringSubmatch(content, -1)

	for _, child := range arr {

		slotStr := child[0]
		slotName := child[1]

		tag := slotStr[0]
		var value interface{}

		if string(tag) == "(" { // fixed
			if indexMap[slotName] < 0 {
				indexMap[slotName] = 0
			}

			mode := len(dataMap[slotName])
			if mode == 0 {
				mode = 1
			}
			index := indexMap[slotName] % mode

			dt, ok := dataMap[slotName]
			if ok && len(dt) > 0 {
				value = dt[index]
			}

		} else if string(tag) == "[" { // seq
			indexMap[slotName] = indexMap[slotName] + 1

			mode := len(dataMap[slotName])
			if mode == 0 {
				mode = 1
			}
			index := indexMap[slotName] % mode

			dt, ok := dataMap[slotName]
			if ok && len(dt) > 0 {
				value = dt[index]
			}

		} else if string(tag) == "{" { // random
			mode := len(dataMap[slotName])
			if mode == 0 {
				mode = 1
			}

			dt, ok := dataMap[slotName]
			if ok && len(dt) > 0 {
				value = dt[commonUtils.RandNum(mode)]
			}
		}

		if value == nil {
			value = ""
		}

		ret = strings.Replace(ret, slotStr, value.(string), 1)
	}

	return
}

func (s *ArticleService) getDataMap(numMap map[string]int, nameMap map[string]string, field *model.DefField,
	resFile string, resType string, sheet string) (ret map[string][]interface{}) {
	ret = map[string][]interface{}{}

	field.Rand = false
	for key, _ := range numMap {
		originTotal := vari.GlobalVars.Total
		vari.GlobalVars.Total = constant.MaxNumb // load all words

		slct, ok := nameMap[key]
		if ok {
			field.Select = slct
		} else {
			field.Select = key
		}

		valueMap, _ := s.ResService.getResValue(resFile, resType, sheet, field)
		ret[key] = valueMap[field.Select]

		vari.GlobalVars.Total = originTotal // rollback
	}

	return
}

func (s *ArticleService) getNumMap(content string) (numMap map[string]int, nameMap map[string]string, indexMap map[string]int, contentWithoutComments string) {
	numMap = map[string]int{}
	nameMap = map[string]string{}
	indexMap = map[string]int{}
	arrWithoutComments := make([]string, 0)

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

		arrWithoutComments = append(arrWithoutComments, line)

		regxSeq := regexp.MustCompile(`\[((?U).*)\]`)
		arrSeq := regxSeq.FindAllStringSubmatch(line, -1)

		regxRand := regexp.MustCompile(`\{((?U).*)\}`)
		arrRand := regxRand.FindAllStringSubmatch(line, -1)

		arr := append(arrSeq, arrRand...)

		for _, child := range arr {
			name := child[1]
			i, ok := numMap[name]
			if !ok {
				numMap[name] = 1
			} else {
				numMap[name] = i + 1
			}

			_, ok2 := indexMap[name]
			if !ok2 {
				indexMap[name] = -1
			}
		}
	}

	contentWithoutComments = strings.Join(arrWithoutComments, "\n")

	return
}

func (s *ArticleService) GenArticle(lines []interface{}) {
	var filePath = logUtils.OutputFileWriter.Name()
	defer logUtils.OutputFileWriter.Close()
	fileUtils.RmFile(filePath)

	for index, line := range lines {
		articlePath := s.genArticleFiles(filePath, index)
		fileWriter, _ := os.OpenFile(articlePath, os.O_RDWR|os.O_CREATE, 0777)
		fmt.Fprint(fileWriter, line)
		fileWriter.Close()
	}
}

func (s *ArticleService) genArticleFiles(pth string, index int) (ret string) {
	pfix := fmt.Sprintf("%03d", index+1)

	ret = strings.TrimSuffix(pth, filepath.Ext(pth))
	ret += "-" + pfix + filepath.Ext(pth)

	return
}
