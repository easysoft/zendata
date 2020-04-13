package action

import (
	"github.com/easysoft/zendata/src/model"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func Generate(file string, count int, fields string, out string, table string) {
	definition := model.Definition{}

	yamlContent, err := ioutil.ReadFile(file)
	if err == nil {
		logUtils.Screen(string(yamlContent))

		err = yaml.Unmarshal(yamlContent, &definition)

		if err == nil {
			logUtils.Screen("===")
		}
	}
}
