package service

import (
	"fmt"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/domain"
	commonUtils "github.com/easysoft/zendata/pkg/utils/common"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	stringUtils "github.com/easysoft/zendata/pkg/utils/string"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	MaxLen           = 10
	IgnoreWords      = []string{"了", "的"}
	IgnoreCategories = []string{"姓", "名字", "介词"}
)

type ArticleService struct {
	ResService *ResService `inject:""`
}

func (s *ArticleService) CreateArticleField(field *domain.DefField) {
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

func (s *ArticleService) getDataMap(numMap map[string]int, nameMap map[string]string, field *domain.DefField,
	resFile string, resType string, sheet string) (ret map[string][]interface{}) {
	ret = map[string][]interface{}{}

	field.Rand = false
	for key, _ := range numMap {
		originTotal := vari.GlobalVars.Total
		vari.GlobalVars.Total = consts.MaxNumb // load all words

		slct, ok := nameMap[key]
		if ok {
			field.Select = slct
		} else {
			field.Select = key
		}

		valueMap, _ := s.ResService.GetResValueFromExcelOrYaml(resFile, resType, sheet, field)
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

func (s *ArticleService) GenYamlFromArticle(file string) {
	startTime := time.Now().Unix()

	content := fileUtils.ReadFile(file)
	words := s.LoadAllWords()

	templ := s.replaceWords(content, words)
	yamlObj := domain.DefArticle{Type: "article", Content: templ, Author: "zendata",
		From: "words.v1", Title: "Template", Version: "1.1"}
	bytes, _ := yaml.Marshal(&yamlObj)
	yamlStr := string(bytes)

	outFile := ""
	if vari.GlobalVars.Output != "" {
		vari.GlobalVars.Output = fileUtils.AddSepIfNeeded(vari.GlobalVars.Output)
		outFile = filepath.Join(vari.GlobalVars.Output, fileUtils.ChangeFileExt(filepath.Base(file), ".yaml"))
		fileUtils.WriteFile(outFile, yamlStr)

	} else {
		logUtils.PrintTo(yamlStr)
	}

	entTime := time.Now().Unix()
	logUtils.PrintTo(i118Utils.I118Prt.Sprintf("generate_article_templ", outFile, entTime-startTime))
}

func (s *ArticleService) replaceWords(content string, words map[string]string) (ret string) {
	runeArr := []rune(content)
	newRuneArr := make([]rune, 0)
	lastUsedWordOfCategoryMap := map[string]string{}
	for i := 0; i < len(runeArr); {
		found := false
		for j := MaxLen; j >= 0; j-- {
			end := i + j
			if end > len(runeArr) {
				end = len(runeArr)
			}

			chars := runeArr[i:end]
			str := ""
			for _, char := range chars {
				str += string(char)
			}

			val, ok := words[str]
			if ok {
				if str == "有" {
					logUtils.PrintTo("")
				}

				lastOne, ok := lastUsedWordOfCategoryMap[val]

				new := ""
				if ok && lastOne == str {
					new = "(" + val + ")"
				} else {
					new = "{" + val + "}"
				}
				itemArr := []rune(new)
				newRuneArr = append(newRuneArr, itemArr...)

				lastUsedWordOfCategoryMap[val] = str // update

				i = end
				found = true
				break
			}
		}

		if !found {
			newRuneArr = append(newRuneArr, runeArr[i])

			i++
		}
	}

	ret = string(newRuneArr)

	return
}

func (s *ArticleService) LoadAllWords() (ret map[string]string) {
	ret = map[string]string{}

	rows, _ := vari.DB.Table("words_v1").Where("true").Select("*").Rows()
	defer rows.Close()

	columns, err := rows.Columns()
	colNum := len(columns)

	colIndexToCategoryName := map[int]string{}
	for index, col := range columns {
		colIndexToCategoryName[index] = col
	}

	// build an empty string array to retrieve row
	var record = make([]interface{}, colNum)
	for i, _ := range record {
		var itf string
		record[i] = &itf
	}

	for rows.Next() {
		err = rows.Scan(record...)
		if err != nil {
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_parse_row", err.Error()))
			return
		}

		for index := len(record) - 1; index >= 0; index-- {
			word := record[1].(*string)
			category := colIndexToCategoryName[index]
			isBelowToCategory := record[index].(*string)

			if *isBelowToCategory == "y" {
				if !stringUtils.StrInArr(category, IgnoreCategories) &&
					!stringUtils.StrInArr(*word, IgnoreWords) {

					ret[*word] = category
				}

				break
			}
		}
	}

	return
}
