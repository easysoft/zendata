package gen

import (
	"github.com/easysoft/zendata/src/model"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func LoadRootDef(file string) model.DefData {
	def := model.DefData{}

	yamlContent, err := ioutil.ReadFile(file)
	if err != nil {
		logUtils.Screen("fail to read " + file)
		return def
	}

	err = yaml.Unmarshal(yamlContent, &def)
	if err != nil {
		logUtils.Screen("fail to parse " + file)
		return def
	}

	return def
}
