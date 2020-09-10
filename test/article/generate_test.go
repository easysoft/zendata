package main

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v3"
	"strconv"
	"strings"
	"testing"
)

const (
	strLeft = "“"
	strRight = "”"

	expLeft = "（"
	expRight = "）"
)

func TestGenerate(ts *testing.T) {
	files := make([]string, 0)
	getFilesInDir("xdoc/words-9.3", ".txt", &files)

	for _, filePath := range files {
		article := fileUtils.ReadFile(filePath)
		content := convertToYaml(article)

		newPath := changeFileExt(filePath, ".yaml")
		fileUtils.WriteFile(newPath, content)
	}
}

func convertToYaml(article string) (content string) {
	sections := parseSections(article)

	conf := createDef(constant.ConfigTypeArticle, "words.v1")

	prefix := ""
	for index, section := range sections {
		tye := section["type"]
		val := section["val"]

		if tye == "exp" {
			field := createField(index, prefix, val)
			conf.XFields = append(conf.XFields, field)
			prefix = ""
		} else {
			prefix += val
		}
	}

	bytes, _ := yaml.Marshal(&conf)
	content = string(bytes)

	// convert yaml format by using a map
	m := make(map[string]interface{})
	yaml.Unmarshal([]byte(content), &m)
	bytes, _ = yaml.Marshal(&m)
	content = string(bytes)
	content = strings.Replace(content, "xfields", "\nfields", -1)

	return
}

func createDef(typ, table string) (conf model.DefExport) {
	conf.Title = "automation"
	conf.Author = "ZenData"
	conf.From = table
	conf.Type = typ
	conf.Desc = "Generated from article text automatically"

	return
}

func createField(index int, prefix, exp string) (field model.DefFieldExport) {
	field.Field = strconv.Itoa(index)
	field.Prefix = prefix

	field.Select = getPinyin(exp)
	field.Where = "true"
	field.Rand = true
	field.Limit = 1

	return
}

func parseSections(content string) (sections []map[string]string) {
	strStart := false
	expStart := false

	content = strings.TrimSpace(content)
	runeArr := []rune(content)

	section := ""
	for i := 0; i < len(runeArr); i++ {
		item := runeArr[i]
		str := string(item)

		isCouple, duplicateStr := isCouple(i, runeArr)
		if isCouple {
			section += duplicateStr
			i += 1
		} else if strStart && str == strRight { // str close
			addSection(section, "str", &sections)

			strStart = false
			section = ""
		} else if expStart && str == expRight { // exp close
			addSection(section, "exp", &sections)

			expStart = false
			section = ""
		} else if !strStart && !expStart && str == strLeft { // str start
			if section != "" && strings.TrimSpace(section) != "+" {
				addSection(section, "str", &sections)
			}

			strStart = true
			section = ""
		} else if !strStart && !expStart && str == expLeft { // exp start
			if section != "" && strings.TrimSpace(section) != "+" {
				addSection(section, "str", &sections)
			}

			expStart = true
			section = ""
		} else {
			section += str
		}
	}

	return
}

func addSection(str, typ string, arr *[]map[string]string) {
	mp := map[string]string{}
	mp["type"] = typ
	mp["val"] = str

	*arr = append(*arr, mp)
}

func isCouple(i int, arr []rune) (isCouple bool, duplicateStr string) {
	if string(arr[i]) == strLeft && (i + 1 < len(arr) && string(arr[i + 1]) == strLeft) {
		isCouple = true
		duplicateStr = string(arr[i])
	} else if string(arr[i]) == strRight && (i + 1 < len(arr) && string(arr[i + 1]) == strRight) {
		isCouple = true
		duplicateStr = string(arr[i])
	} else if string(arr[i]) == expLeft && (i + 1 < len(arr) && string(arr[i + 1]) == expLeft) {
		isCouple = true
		duplicateStr = string(arr[i])
	} else if string(arr[i]) == expRight && (i + 1 < len(arr) && string(arr[i + 1]) == expRight) {
		isCouple = true
		duplicateStr = string(arr[i])
	}

	return
}