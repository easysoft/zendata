package action

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/gen"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func Generate(file string, count int, fields string, out string, table string) {
	definition := model.Definition{}

	yamlContent, err := ioutil.ReadFile(file)
	if err != nil {
		logUtils.Screen("fail to read " + file)
		return
	}

	err = yaml.Unmarshal(yamlContent, &definition)
	if err != nil {
		logUtils.Screen("fail to parse " + file)
		return
	}

	gen.Generate(definition, count, fields, out, table)
}
