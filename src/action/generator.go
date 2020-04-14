package action

import (
	"fmt"
	"github.com/easysoft/zendata/src/gen"
	"github.com/easysoft/zendata/src/model"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
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

	rows := gen.Generate(&definition, total, fields, out, table)
	Print(rows)
}


func Print(rows [][]string) {
	width := stringUtils.GetNumbWidth(len(rows))

	for i, cols := range rows {
		line := ""
		for j, col := range cols {
			if j >0 {
				line = line + ", "
			}
			line = line + col
		}

		idStr := fmt.Sprintf("%" + strconv.Itoa(width) + "d", i+1)
		logUtils.Screen(fmt.Sprintf("%s: %s", idStr, line))
	}
}