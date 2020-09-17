package service

import (
	"fmt"
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v3"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	strLeft = "“"
	strRight = "”"

	expLeft = "（"
	expRight = "）"

	table = "words.v1"
)
var (
	compares = []string{"=", "!=", ">", "<"}
)

func ConvertArticle(src, dist string) {
	files := make([]string, 0)
	if !fileUtils.IsDir(src) { //  file
		pth, _ := filepath.Abs(src)
		files = append(files, pth)

		if dist == "" { dist = fileUtils.AddSepIfNeeded(path.Dir(pth)) }
	} else {
		fileUtils.GetFilesInDir(src, ".txt", &files)
		if dist == "" { dist = fileUtils.AddSepIfNeeded(src) }
	}

	for _, filePath := range files {
		yamlPaths := convertSentYaml(filePath, dist)
		convertMainYaml(yamlPaths, filePath, dist)
	}
}

func convertSentYaml(filePath, dist string) (yamlPaths []string) {
	article := fileUtils.ReadFile(filePath)
	sections := parseSections(article)
	paragraphs := groupSections(sections)

	for paragIndex, parag := range paragraphs {

		for sentIndex, sent := range parag {
			fileSeq := fmt.Sprintf("p%02d-s%02d", paragIndex + 1, sentIndex + 1)

			conf := createDef(constant.ConfigTypeArticle, table, fileUtils.GetRelatPath(filePath))

			prefix := ""
			for sectIndex, sect := range sent { // each sent saved as a yaml file
				fieldSeq := fmt.Sprintf("%d-%d-%d", paragIndex + 1, sentIndex + 1, sectIndex + 1)
				if sect.Type == "exp" {
					fields := createFields(fieldSeq, prefix, sect.Val)
					conf.XFields = append(conf.XFields, fields...)

					prefix = ""
				} else {
					prefix += sect.Val

					if prefix != "" && sectIndex == len(sent) - 1 { // last section
						//if strings.LastIndex(prefix, "\n\n") == len(prefix) - 2 {
						//	prefix = prefix[:len(prefix) - 1]
						//}

						field := model.DefFieldExport{Field: fieldSeq, Prefix: prefix}
						conf.XFields = append(conf.XFields, field)
						prefix = ""
					}
				}
			}

			bytes, _ := yaml.Marshal(&conf)
			content := string(bytes)

			// convert yaml format by using a map
			m := make(map[string]interface{})
			yaml.Unmarshal([]byte(content), &m)
			bytes, _ = yaml.Marshal(&m)
			content = string(bytes)
			content = strings.Replace(content, "xfields", "\nfields", -1)

			yamlPath := fileUtils.AddSepIfNeeded(dist) +
				fileUtils.ChangeFileExt(path.Base(filePath), "-") + fileSeq + ".yaml"
			fileUtils.WriteFile(yamlPath, content)

			yamlPaths = append(yamlPaths, yamlPath)
		}
	}

	return
}

func convertMainYaml(yamlPaths []string, filePath, dist string) {
	conf := createArticle(constant.ConfigTypeArticle, fileUtils.GetRelatPath(filePath))

	for index, file := range yamlPaths {
		path := strings.TrimPrefix(file, dist)
		field := model.ArticleField{Field: strconv.Itoa(index + 1), Range: path}
		conf.XFields = append(conf.XFields, field)
	}

	bytes, _ := yaml.Marshal(&conf)
	content := string(bytes)

	// convert yaml format by using a map
	m := make(map[string]interface{})
	yaml.Unmarshal([]byte(content), &m)
	bytes, _ = yaml.Marshal(&m)
	content = string(bytes)
	content = strings.Replace(content, "xfields", "\nfields", -1)

	yamlPath := fileUtils.AddSepIfNeeded(dist) + fileUtils.ChangeFileExt(path.Base(filePath), ".yaml")
	fileUtils.WriteFile(yamlPath, content)

	relatPath := fileUtils.GetRelatPath(yamlPath)
	yamlPaths = append(yamlPaths, relatPath)
}

func createDef(typ, table, filePath string) (conf model.DefExport) {
	conf.Title = "automation"
	conf.Author = "ZenData"
	conf.Type = typ
	conf.Desc = "Generated from article " + filePath

	if table != "" {
		conf.From = table
	}

	return
}

func createArticle(typ, filePath string) (conf model.Article) {
	conf.Title = "automation"
	conf.Author = "ZenData"
	conf.Type = typ
	conf.Desc = "Generated from article " + filePath

	return
}

func createFields(seq string, prefix, exp string) (fields []model.DefFieldExport) {
	field := model.DefFieldExport{}
	field.Field = seq
	field.Prefix = prefix
	field.Rand = true
	field.Limit = 1

	// deal with exp like S：名词-姓+名词-名字=F
	exp = strings.ToLower(strings.TrimSpace(exp))
	expArr := []rune(exp)

	if string(expArr[0]) == "s" && (string(expArr[1]) == ":" || string(expArr[1]) == "：") {
		exp = string(expArr[2:])
		expArr = expArr[2:]
	}

	if strings.Index(exp, "=") == len(exp) - 2 {
		exp = string(expArr[:len(expArr) - 2])
		field.Select = exp
		field.Where = string(expArr[len(expArr) - 1])
	} else {
		field.Select = exp
		field.Where = ""
	}

	if strings.Index(field.Select, "+") < 0 {
		fields = append(fields, field)
	} else if strings.Index(field.Select, "+") > 0 { // include more than one field, split to two
		items := strings.Split(field.Select, "+")
		for _, item := range items {
			var objClone interface{} = field
			fieldClone := objClone.(model.DefFieldExport)
			fieldClone.Select = item
			fieldClone.Where = field.Where

			fields = append(fields, fieldClone)
		}
	}

	return
}

func parseSections(content string) (sections []model.ArticleSent) {
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

			if i == len(runeArr) - 1 {
				addSection(section, "str", &sections)
			}

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

			if str == "。" {
				if i < len(runeArr) && string(runeArr[i+1]) == strRight {
					i += 1
				}

				addSection(section, "str", &sections)

				strStart = false
				expStart = false
				section = ""
			} else if str == "\n" {
				// get all \n
				for j := i+1; j < len(runeArr); j++ {
					if string(runeArr[j]) == "\n" {
						section += str
						i = j
					} else {
						break
					}
				}

				addSection(section, "str", &sections)

				strStart = false
				expStart = false
				section = ""
			} else if i == len(runeArr) - 1 {
				addSection(section, "str", &sections)
			}
		}
	}

	return
}

func groupSections(sectionArr []model.ArticleSent) (paragraphs [][][]model.ArticleSent) {
	sections := make([]model.ArticleSent, 0)
	sentences := make([][]model.ArticleSent, 0)

	for index := 0; index < len(sectionArr); index++ {
		section := sectionArr[index]
		sections = append(sections, section)

		if section.IsParag {
			sentences = append(sentences, sections)
			paragraphs = append(paragraphs, sentences)

			sentences = make([][]model.ArticleSent, 0)
			sections = make([]model.ArticleSent, 0)
		} else if section.IsSent {
			if index < len(sectionArr) - 1 && sectionArr[index+1].IsParag {
				sections = append(sections, sectionArr[index+1])
				sentences = append(sentences, sections)
				paragraphs = append(paragraphs, sentences)

				sections = make([]model.ArticleSent, 0)
				sentences = make([][]model.ArticleSent, 0)

				index += 1
			} else {
				sentences = append(sentences, sections)
				if index == len(sectionArr) - 1 {
					paragraphs = append(paragraphs, sentences)
				}

				sections = make([]model.ArticleSent, 0)
			}

		}
	}

	return
}

func addSection(str, typ string, arr *[]model.ArticleSent) {
	sent := model.ArticleSent{}
	sent.Type = typ
	sent.Val = str

	runeArr := []rune(str)
	end := runeArr[len(runeArr) - 1]
	if string(end) == "\n" {
		sent.IsParag = true
	} else if string(end) == "。" {
		sent.IsSent = true
	}

	*arr = append(*arr, sent)
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

