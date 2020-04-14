package action

import (
	"fmt"
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/gen"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func Generate(file string, total int, fields string, out string, table string) {
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

	rows := gen.Generate(definition, total, fields, out, table)
	Print(rows)
}


func Print(rows [][]string) {
	for i, cols := range rows {
		for j, col := range cols {
			logUtils.Screen(fmt.Sprintf("%d-%d: %s", i, j, col))
		}
	}
}