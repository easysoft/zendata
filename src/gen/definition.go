package gen

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

func LoadDefinitionFromFile(file string) {
	def := model.Definition{}

	yamlContent, err := ioutil.ReadFile(file)
	if err != nil {
		logUtils.Screen("fail to read " + file)
		return
	}

	err = yaml.Unmarshal(yamlContent, &def)
	if err != nil {
		logUtils.Screen("fail to parse " + file)
		return
	}

	if strings.Index(file, "def") != 0 && constant.Definition.Title == "" { // only add the fields in first level yaml file
		constant.Definition = def
	}

	for _, field := range def.Fields {
		constant.LoadedFields[field.Name] = field // add to a map
	}
}
