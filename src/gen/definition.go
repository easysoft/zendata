package gen

import (
	"github.com/easysoft/zendata/src/model"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func LoadRootDef(file string) model.DefData {
	def := model.DefData{}

	yamlContent, err := ioutil.ReadFile(file)
	if err != nil {
		logUtils.Screen(i118Utils.I118Prt.Sprintf("fail_to_read_file", file))
		return def
	}

	err = yaml.Unmarshal(yamlContent, &def)
	if err != nil {
		logUtils.Screen(i118Utils.I118Prt.Sprintf("fail_to_read_file", file))
		return def
	}

	return def
}
