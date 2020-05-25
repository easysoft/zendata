package gen

import (
	"github.com/easysoft/zendata/src/model"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

func LoadRootDef(defaultFile, ymlFile string, fieldsToExport *[]string) model.DefData {
	defaultDef := model.DefData{}
	ymlDef := model.DefData{}

	if defaultFile != "" {
		defaultContent, err := ioutil.ReadFile(defaultFile)
		if err != nil {
			logUtils.Screen(i118Utils.I118Prt.Sprintf("fail_to_read_file", defaultFile))
			return defaultDef
		}
		err = yaml.Unmarshal(defaultContent, &defaultDef)
		if err != nil {
			logUtils.Screen(i118Utils.I118Prt.Sprintf("fail_to_read_file", defaultFile))
			return defaultDef
		}
	}

	yamlContent, err := ioutil.ReadFile(ymlFile)
	if err != nil {
		logUtils.Screen(i118Utils.I118Prt.Sprintf("fail_to_read_file", ymlFile))
		return ymlDef
	}
	err = yaml.Unmarshal(yamlContent, &ymlDef)
	if err != nil {
		logUtils.Screen(i118Utils.I118Prt.Sprintf("fail_to_read_file", ymlFile))
		return ymlDef
	}

	MergerDefine(&defaultDef, &ymlDef)

	fieldsToExportIsEmpty := len(*fieldsToExport) == 0
	for _, field := range ymlDef.Fields {
		if fieldsToExportIsEmpty {
			*fieldsToExport = append(*fieldsToExport, field.Field)
		}
	}

	return defaultDef
}

func MergerDefine(defaultDef, ymlDef *model.DefData) {
	defaultFieldMap := map[string]*model.DefField{}
	ymlFieldMap := map[string]*model.DefField{}

	for _, field := range defaultDef.Fields {
		CreatePathToFieldMap(&field, defaultFieldMap)
	}

	for _, field := range ymlDef.Fields {
		CreatePathToFieldMap(&field, ymlFieldMap)
	}

	for path, field := range ymlFieldMap {
		parent, exist := defaultFieldMap[path]
		if exist {
			CopyField(*field, parent)
		}
	}

	for _, field := range ymlFieldMap {
		if strings.Index(field.Path, "~~") > -1 { continue } // only for top fields

		_, exist := defaultFieldMap[field.Path]
		if !exist {
			defaultDef.Fields = append(defaultDef.Fields, *field)
		}
	}
}

func CreatePathToFieldMap(field *model.DefField, mp map[string]*model.DefField) {
	if field.Path == "" { // root
		field.Path = field.Field
	}

	if len(field.Fields) > 0 {
		for _, child := range field.Fields {
			child.Path = field.Path + "~~" + child.Field

			CreatePathToFieldMap(&child, mp)
		}
	} else {
		mp[field.Path] = field
	}
}

func CopyField(child model.DefField, parent *model.DefField) {
	parent.Range = child.Range

	if child.Note != "" {
		parent.Note = child.Note
	}
	if child.Prefix != "" {
		parent.Prefix = child.Prefix
	}
	if child.Postfix != "" {
		parent.Postfix = child.Postfix
	}
	if child.Loop != 0 {
		parent.Loop = child.Loop
	}
	if child.Loopfix != "" {
		parent.Loopfix = child.Loopfix
	}
	if child.Format != "" {
		parent.Format = child.Format
	}

	if child.From != "" {
		parent.From = child.From
	}
	if child.Select != "" {
		parent.Select = child.Select
	}
	if child.Where != "" {
		parent.Where = child.Where
	}
	if child.Use != "" {
		parent.Use = child.Use
	}

	if child.Precision != 0 {
		parent.Precision = child.Precision
	}
}
