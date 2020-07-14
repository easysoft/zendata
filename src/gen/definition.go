package gen

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

func LoadConfigDef(defaultFile, configFile string, fieldsToExport *[]string) model.DefData {
	defaultDef := model.DefData{}
	configDef := model.DefData{}

	// load defaultDef
	if defaultFile != "" {
		defaultContent, err := ioutil.ReadFile(defaultFile)
		defaultContent = ReplaceSpecialChars(defaultContent)
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

	// load configDef
	yamlContent, err := ioutil.ReadFile(configFile)
	yamlContent = ReplaceSpecialChars(yamlContent)
	if err != nil {
		logUtils.Screen(i118Utils.I118Prt.Sprintf("fail_to_read_file", configFile))
		return configDef
	}
	err = yaml.Unmarshal(yamlContent, &configDef)
	if err != nil {
		logUtils.Screen(i118Utils.I118Prt.Sprintf("fail_to_parse_file", configFile))
		return configDef
	}

	// use all fields as default
	if len(*fieldsToExport) == 0 {
		for _, field := range configDef.Fields {
			*fieldsToExport = append(*fieldsToExport, field.Field)
		}
	}

	MergerDefine(&defaultDef, &configDef)

	return defaultDef
}

func MergerDefine(defaultDef, configDef *model.DefData) {
	defaultFieldMap := map[string]*model.DefField{}
	configFieldMap := map[string]*model.DefField{}
	sortedKeys := make([]string, 0)

	for i := range defaultDef.Fields {
		CreatePathToFieldMap(&defaultDef.Fields[i], defaultFieldMap, nil)
	}

	for i := range configDef.Fields {
		CreatePathToFieldMap(&configDef.Fields[i], configFieldMap, &sortedKeys)
	}

	// overwrite
	for path, field := range configFieldMap {
		parent, exist := defaultFieldMap[path]
		if exist {
			CopyField(*field, parent)
			defaultFieldMap[path] = parent
		}
	}

	// append
	for _, key := range sortedKeys {
		field := configFieldMap[key]
		if strings.Index(field.Path, "~~") > -1 { continue } // ignore no-top fields

		_, exist := defaultFieldMap[field.Path]
		if !exist {
			defaultDef.Fields = append(defaultDef.Fields, *field)
		}
	}
}

func CreatePathToFieldMap(field *model.DefField, mp map[string]*model.DefField, keys *[]string) {
	if field.Path == "" { // root
		field.Path = field.Field
	}

	path := field.Path
	//logUtils.Screen(path + " -> " + field.Field)
	mp[path] = field

	if keys != nil {
		*keys = append(*keys, path)
	}

	if len(field.Fields) > 0 {
		for i := range field.Fields {
			field.Fields[i].Path = field.Path + "~~" + field.Fields[i].Field

			CreatePathToFieldMap(&field.Fields[i], mp, keys)
		}
	}
}

func CopyField(child model.DefField, parent *model.DefField) {
	if child.Note != "" {
		(*parent).Note = child.Note
	}
	if child.Range != "" {
		(*parent).Range = child.Range
	}
	if child.Prefix != "" {
		(*parent).Prefix = child.Prefix
	}
	if child.Postfix != "" {
		(*parent).Postfix = child.Postfix
	}
	if child.Loop != 0 {
		(*parent).Loop = child.Loop
	}
	if child.Loopfix != "" {
		(*parent).Loopfix = child.Loopfix
	}
	if child.Format != "" {
		(*parent).Format = child.Format
	}

	if child.From != "" {
		(*parent).From = child.From
	}
	if child.Select != "" {
		(*parent).Select = child.Select
	}
	if child.Where != "" {
		(*parent).Where = child.Where
	}
	if child.Use != "" {
		(*parent).Use = child.Use
	}

	if child.Precision != 0 {
		(*parent).Precision = child.Precision
	}
}

func ReplaceSpecialChars(bytes []byte) []byte {
	str := string(bytes)

	ret := ""
	for _, line := range strings.Split(str, "\n") {
		if strings.Index(strings.TrimSpace(line), "range") == 0 {
			line = strings.ReplaceAll(line,"[", string(constant.LeftChar))
			line = strings.ReplaceAll(line,"]", string(constant.RightChar))
		}

		ret += line + "\n"
	}

	return []byte(ret)
}

