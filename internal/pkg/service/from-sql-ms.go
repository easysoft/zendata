package service

import (
	"fmt"
	"github.com/easysoft/zendata/internal/pkg/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"gopkg.in/yaml.v2"
	"strings"
)

type MsSqlParseService struct {
	SqlParseService *SqlParseService `inject:""`
}

func (s *MsSqlParseService) genKeysYaml(pkMap map[string]string) {
	// gen key yaml files
	inst := model.ZdInstances{Title: "keys", Desc: "automated export"}

	for tableName, keyCol := range pkMap {
		item := model.ZdInstancesItem{}

		item.Instance = fmt.Sprintf("%s_%s", tableName, keyCol)
		item.Range = "1-100000"

		inst.Instances = append(inst.Instances, item)
	}

	bytes, _ := yaml.Marshal(&inst)
	content := strings.ReplaceAll(string(bytes), "'-'", "\"\"")

	if vari.GlobalVars.Output != "" {
		vari.GlobalVars.Output = fileUtils.AddSepIfNeeded(vari.GlobalVars.Output)
		outFile := vari.GlobalVars.Output + "keys.yaml"
		fileUtils.WriteFile(outFile, content)

	} else {
		logUtils.PrintTo(content)
	}
}
