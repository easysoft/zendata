package gen

import (
	"github.com/easysoft/zendata/internal/pkg/domain"
	"github.com/easysoft/zendata/internal/pkg/helper"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	stringUtils "github.com/easysoft/zendata/pkg/utils/string"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

func LoadDataContentDef(filesContents [][]byte, fieldsToExport *[]string) (ret domain.DefData) {
	ret = domain.DefData{}
	for _, f := range filesContents {
		right := LoadContentDef(f)
		ret = MergeDef(ret, right, fieldsToExport)
	}

	return
}

func LoadContentDef(content []byte) (ret domain.DefData) {
	content = helper.ReplaceSpecialChars(content)
	err := yaml.Unmarshal(content, &ret)
	if err != nil {
		logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("fail_to_parse_file"), color.FgCyan)
		return
	}

	return
}

func LoadFilesContents(files []string) (contents [][]byte) {
	contents = make([][]byte, 0)
	for _, f := range files {
		if f == "" {
			continue
		}
		pathDefaultFile := fileUtils.GetAbsolutePath(f)
		if !fileUtils.FileExist(pathDefaultFile) {
			return
		}
		content, err := ioutil.ReadFile(pathDefaultFile)
		if err != nil {
			logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("fail_to_parse_file"), color.FgCyan)
			return
		}
		contents = append(contents, content)
	}

	return
}

func LoadDataDef(files []string, fieldsToExport *[]string) (ret domain.DefData) {
	newFiles := make([]string, 0)
	for _, f := range files {
		if f == "" {
			continue
		}
		newFiles = append(newFiles, f)
	}

	ret = domain.DefData{}
	for _, f := range newFiles {
		right := LoadDef(f)

		ret = MergeDef(ret, right, fieldsToExport)
	}

	return
}

func LoadDef(file string) (ret domain.DefData) {
	pathDefaultFile := fileUtils.GetAbsolutePath(file)

	if !fileUtils.FileExist(pathDefaultFile) {
		return
	}

	defaultContent, err := ioutil.ReadFile(pathDefaultFile)
	defaultContent = helper.ReplaceSpecialChars(defaultContent)
	if err != nil {
		logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("fail_to_read_file", pathDefaultFile), color.FgCyan)
		return
	}
	err = yaml.Unmarshal(defaultContent, &ret)
	if err != nil {
		logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("fail_to_parse_file", pathDefaultFile), color.FgCyan)
		return
	}

	return
}

func MergeDef(defaultDef domain.DefData, configDef domain.DefData, fieldsToExport *[]string) domain.DefData {
	if configDef.Type == "article" && configDef.Content != "" {
		convertArticleContent(&configDef)
	}

	mergerDefine(&defaultDef, &configDef, fieldsToExport)
	orderFields(&defaultDef, *fieldsToExport)

	for index, _ := range defaultDef.Fields {
		if vari.GlobalVars.Trim {
			defaultDef.Fields[index].Prefix = ""
			defaultDef.Fields[index].Postfix = ""
		}
	}

	return defaultDef
}

func convertArticleContent(config *domain.DefData) {
	field := domain.DefField{}
	field.Type = config.Type
	field.From = config.From
	field.Range = "`" + config.Content + "`"

	config.Fields = append(config.Fields, field)
}

func mergerDefine(defaultDef, configDef *domain.DefData, fieldsToExport *[]string) {
	isSetFieldsToExport := false
	if len(*fieldsToExport) > 0 {
		isSetFieldsToExport = true
	}

	defaultFieldMap := map[string]*domain.DefField{}
	configFieldMap := map[string]*domain.DefField{}
	sortedKeys := make([]string, 0)

	//if configDef.Type != "" {
	//	vari.GlobalVars.DefDataType = configDef.Type
	//} else if defaultDef.Type != "" {
	//	vari.GlobalVars.DefDataType = defaultDef.Type
	//} else {
	//	vari.GlobalVars.DefDataType = consts.DefTypeText
	//}

	if configDef.Content != "" && defaultDef.Content == "" {
		defaultDef.Content = configDef.Content
	}
	if configDef.From != "" && defaultDef.From == "" {
		defaultDef.From = configDef.From
	}
	if configDef.Type != "" && defaultDef.Type == "" {
		defaultDef.Type = configDef.Type
	}

	for i, field := range defaultDef.Fields {
		if !isSetFieldsToExport {
			*fieldsToExport = append(*fieldsToExport, field.Field)
		}

		defaultDef.Fields[i].FileDir = vari.GlobalVars.ConfigFileDir
		CreatePathToFieldMap(&defaultDef.Fields[i], defaultFieldMap, nil)
	}
	for i, field := range configDef.Fields {
		vari.GlobalVars.TopFieldMap[field.Field] = field
		if !isSetFieldsToExport {
			if !stringUtils.StrInArr(field.Field, *fieldsToExport) {
				*fieldsToExport = append(*fieldsToExport, field.Field)
			}
		}

		configDef.Fields[i].FileDir = vari.GlobalVars.ConfigFileDir
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
		if field == nil || strings.Index(field.Path, "~~") > -1 {
			continue
		} // ignore no-top fields

		_, exist := defaultFieldMap[field.Path]
		if !exist {
			defaultDef.Fields = append(defaultDef.Fields, *field)
		}
	}

	for _, field := range defaultDef.Fields {
		vari.GlobalVars.TopFieldMap[field.Field] = field
	}
}

func orderFields(defaultDef *domain.DefData, fieldsToExport []string) {
	mp := map[string]domain.DefField{}
	for _, field := range defaultDef.Fields {
		mp[field.Field] = field
	}

	fields := make([]domain.DefField, 0)
	for _, fieldName := range fieldsToExport {
		fields = append(fields, mp[fieldName])
	}

	defaultDef.Fields = fields
}

func CreatePathToFieldMap(field *domain.DefField, mp map[string]*domain.DefField, keys *[]string) {
	if field.Path == "" { // root
		field.Path = field.Field
	}

	path := field.Path
	//logUtils.Screen(path + " -> " + field.ZdField)
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

func CopyField(child domain.DefField, parent *domain.DefField) {
	if child.Note != "" {
		(*parent).Note = child.Note
	}
	if child.Range != "" {
		(*parent).Range = child.Range
	}

	(*parent).Prefix = child.Prefix
	(*parent).Postfix = child.Postfix
	(*parent).Divider = child.Divider

	if child.Loop != "" {
		(*parent).Loop = child.Loop
	}
	(*parent).Loopfix = child.Loopfix
	(*parent).Format = child.Format

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
	if child.Length != 0 {
		(*parent).Length = child.Length
	}
}