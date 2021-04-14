package action

import (
	"github.com/easysoft/zendata/src/gen/helper"
	"github.com/easysoft/zendata/src/model"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"gopkg.in/yaml.v3"
	"path/filepath"
	"time"
)

var (
	MaxLen = 10

	IgnoreWords      = []string{"了", "的"}
	IgnoreCategories = []string{"介词"}
)

func ParseArticle(file string, out string) {
	startTime := time.Now().Unix()

	content := fileUtils.ReadFile(file)
	words := helper.LoadAllWords()

	templ := replaceWords(content, words)
	yamlObj := model.DefArticle{Type: "article", Content: templ, Author: "zendata",
		From: "words.v1", Title: "Template", Version: "1.1"}
	bytes, _ := yaml.Marshal(&yamlObj)
	yamlStr := string(bytes)

	outFile := ""
	if out != "" {
		out = fileUtils.AddSepIfNeeded(out)
		outFile = filepath.Join(out, fileUtils.ChangeFileExt(filepath.Base(file), ".yaml"))
		WriteToFile(outFile, yamlStr)
	} else {
		logUtils.PrintTo(yamlStr)
	}

	entTime := time.Now().Unix()
	logUtils.PrintTo(i118Utils.I118Prt.Sprintf("generate_article_templ", outFile, entTime-startTime))
}

func replaceWords(content string, words map[string]string) (ret string) {
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

				if !ignoreToReplace(str, val) {
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
		}

		if !found {
			newRuneArr = append(newRuneArr, runeArr[i])

			i++
		}
	}

	ret = string(newRuneArr)

	return
}

func ignoreToReplace(val, category string) bool {
	if stringUtils.StrInArr(val, IgnoreWords) {
		return true
	}
	if stringUtils.StrInArr(category, IgnoreCategories) {
		return true
	}

	//if category == "姓" || category == "名字" {
	//	return true
	//}

	return false
}
