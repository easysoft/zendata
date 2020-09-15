package service

import (
	"fmt"
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
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
	//src = "data/words"
	//dist = "demo"
)
var (
	compares = []string{"=", "!=", ">", "<"}
)

func ConvertArticle(src, dist string) {
	files := make([]string, 0)
	if !fileUtils.IsDir(src) {
		pth, _ := filepath.Abs(src)
		files = append(files, pth)

		if dist == "" {
			dist = path.Dir(pth)
		}
	} else {
		fileUtils.GetFilesInDir(src, ".txt", &files)
		if dist == "" {
			dist = src
		}
	}

	for _, filePath := range files {
		article := fileUtils.ReadFile(filePath)
		content := convertToYaml(article, filePath)

		newPath := fileUtils.AddSepIfNeeded(dist) + fileUtils.ChangeFileExt(path.Base(filePath), ".yaml")
		fileUtils.WriteFile(newPath, content)
	}
}

func convertToYaml(article, filePath string) (content string) {
	sections := parseSections(article)

	conf := createDef(constant.ConfigTypeArticle, table, filePath)

	prefix := ""
	for index, section := range sections {
		tye := section["type"]
		val := section["val"]

		if tye == "exp" {
			fields := createFields(index, prefix, val)
			conf.XFields = append(conf.XFields, fields...)

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

func createDef(typ, table, filePath string) (conf model.DefExport) {
	conf.Title = "automation"
	conf.Author = "ZenData"
	conf.From = table
	conf.Type = typ
	conf.Desc = "Generated from article " + filePath

	return
}

func createFields(index int, prefix, exp string) (fields []model.DefFieldExport) {
	field := model.DefFieldExport{}
	field.Field = strconv.Itoa(index)
	field.Prefix = prefix
	field.Rand = true
	field.Limit = 1

	// deal with exp like S：名词-姓+名词-名字=F
	exp = strings.ToLower(strings.TrimSpace(exp))
	expArr := []rune(exp)

	if string(expArr[0]) == "s" && (string(expArr[1]) == ":" || string(expArr[1]) == "：") {
		exp = string(expArr[2:])
		expArr = expArr[2:]
		field.UseLastSameValue = true
	}

	if strings.Index(exp, "=") == len(exp) - 2 {
		exp = string(expArr[:len(expArr) - 2])
		field.Select = stringUtils.GetPinyin(exp)
		field.Where = fmt.Sprintf("%s = '%s'", field.Select, string(expArr[len(expArr) - 1]))
	} else {
		field.Select = stringUtils.GetPinyin(exp)
		field.Where = "true"
		//field.Where = getPinyin(exp) + " = 'y'"
	}

	if strings.Index(field.Select, "+") < 0 {
		fields = append(fields, field)
	} else if strings.Index(field.Select, "+") > 0 { // include more than one field, split to two
		arr := strings.Split(field.Where, "=")
		right := ""
		if len(arr) > 1 {
			right = arr[1]
		}

		items := strings.Split(field.Select, "+")
		for _, item := range items {
			var objClone interface{} = field
			fieldClone := objClone.(model.DefFieldExport)
			fieldClone.Select = item

			if len(arr) > 1 { // has conditions
				fieldClone.Where = item + " = " + right
			}

			fields = append(fields, fieldClone)
		}
	}

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

