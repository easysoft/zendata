package gen

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

func LoadConfigDef(defaultFile, configFile string, fieldsToExport *[]string) model.DefData {
	defaultDef := model.DefData{}
	configDef := model.DefData{}

	// load defaultDef
	if defaultFile != "" {
		pathDefaultFile := defaultFile
		if !fileUtils.IsAbosutePath(pathDefaultFile) {
			pathDefaultFile = vari.WorkDir + pathDefaultFile
		}

		defaultContent, err := ioutil.ReadFile(pathDefaultFile)
		defaultContent = ReplaceSpecialChars(defaultContent)
		if err != nil {
			logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("fail_to_read_file", pathDefaultFile), color.FgCyan)
			return defaultDef
		}
		err = yaml.Unmarshal(defaultContent, &defaultDef)
		if err != nil {
			logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("fail_to_read_file", pathDefaultFile), color.FgCyan)
			return defaultDef
		}
	}

	// load configDef
	pathConfigFile := configFile
	if !fileUtils.IsAbosutePath(pathConfigFile) {
		pathConfigFile = vari.WorkDir + pathConfigFile
	}

	yamlContent, err := ioutil.ReadFile(pathConfigFile)
	yamlContent = ReplaceSpecialChars(yamlContent)
	if err != nil {
		logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("fail_to_read_file", pathConfigFile), color.FgCyan)
		return configDef
	}
	err = yaml.Unmarshal(yamlContent, &configDef)
	if err != nil {
		logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("fail_to_parse_file", pathConfigFile), color.FgCyan)
		return configDef
	}

	mergerDefine(&defaultDef, &configDef, fieldsToExport)
	orderFields(&defaultDef, *fieldsToExport)

	for _, field := range defaultDef.Fields {
		if vari.Trim {
			field.Prefix = ""
			field.Postfix = ""
		}
	}

	return defaultDef
}

func mergerDefine(defaultDef, configDef *model.DefData, fieldsToExport *[]string) {
	isSetFieldsToExport := false
	if len(*fieldsToExport) > 0 {
		isSetFieldsToExport = true
	}

	defaultFieldMap := map[string]*model.DefField{}
	configFieldMap := map[string]*model.DefField{}
	sortedKeys := make([]string, 0)


	if configDef.Type != "" {
		vari.Type = configDef.Type
	} else if defaultDef.Type != "" {
		vari.Type = defaultDef.Type
	} else {
		vari.Type = constant.ConfigTypeText
	}

	if configDef.From != "" && defaultDef.From == "" {
		defaultDef.From = configDef.From
	}

	for i, field := range defaultDef.Fields {
		if !isSetFieldsToExport {
			*fieldsToExport = append(*fieldsToExport, field.Field)
		}

		CreatePathToFieldMap(&defaultDef.Fields[i], defaultFieldMap, nil)
	}
	for i, field := range configDef.Fields {
		vari.TopFiledMap[field.Field] = field
		if !isSetFieldsToExport {
			if !stringUtils.FindInArr(field.Field, *fieldsToExport) {
				*fieldsToExport = append(*fieldsToExport, field.Field)
			}
		}

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
		if field == nil || strings.Index(field.Path, "~~") > -1 { continue } // ignore no-top fields

		_, exist := defaultFieldMap[field.Path]
		if !exist {
			defaultDef.Fields = append(defaultDef.Fields, *field)
		}
	}

	for _, field := range defaultDef.Fields {
		vari.TopFiledMap[field.Field] = field
	}
}

func orderFields(defaultDef *model.DefData, fieldsToExport []string) {
	mp := map[string]model.DefField{}
	for _, field := range defaultDef.Fields {
		mp[field.Field] = field
	}

	fields := make([]model.DefField, 0)
	for _, fieldName := range fieldsToExport {
		fields = append(fields, mp[fieldName])
	}

	defaultDef.Fields = fields
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

	//if child.Prefix != "" {
		(*parent).Prefix = child.Prefix
	//}
	//if child.Postfix != "" {
		(*parent).Postfix = child.Postfix
	//}

	if child.Loop != "" {
		(*parent).Loop = child.Loop
	}
	//if child.Loopfix != "" {
		(*parent).Loopfix = child.Loopfix
	//}
	//if child.Format != "" {
		(*parent).Format = child.Format
	//}

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
	if child.From != "" {
		(*parent).From = child.From
	}

	if child.Type != "" {
		(*parent).Type = child.Type
	}

	if child.Precision != 0 {
		(*parent).Precision = child.Precision
	}
	if child.Width != 0 {
		(*parent).Width = child.Width
	}
}

func ReplaceSpecialChars(bytes []byte) []byte {
	str := string(bytes)

	inRanges := false // for ranges yaml only
	ret := ""
	for _, line := range strings.Split(str, "\n") {
		if strings.Index(strings.TrimSpace(line), "ranges") == 0 {
			inRanges = true
		} else if len(line) > 0 && string(line[0]) != " " { // not begin with space, ranges end
			inRanges = false
		}

		if strings.Index(strings.TrimSpace(line), "range") == 0 || inRanges {
			line = strings.ReplaceAll(line,"[", string(constant.LeftBrackets))
			line = strings.ReplaceAll(line,"]", string(constant.RightBrackets))
		}

		ret += line + "\n"
	}

	return []byte(ret)
}

