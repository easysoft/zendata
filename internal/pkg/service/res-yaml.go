package service

import (
	"github.com/easysoft/zendata/internal/pkg/helper"
	"github.com/easysoft/zendata/internal/pkg/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"gopkg.in/yaml.v3"
	"os"
)

type ResYamlService struct {
	ResConfigService    *ResConfigService    `inject:""`
	ResRangesService    *ResRangesService    `inject:""`
	ResInstancesService *ResInstancesService `inject:""`
}

func (s *ResYamlService) GetResFromYaml(resFile string) (valueMap map[string][]interface{}) { // , resName string) {
	if vari.GlobalVars.CacheResFileToMap[resFile] != nil { // already cached
		valueMap = vari.GlobalVars.CacheResFileToMap[resFile]
		return
	}

	yamlContent, err := os.ReadFile(resFile)
	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_read_file", resFile))
		return
	}

	yamlContent = helper.ReplaceSpecialChars(yamlContent)

	insts := model.ResInstances{}
	err = yaml.Unmarshal(yamlContent, &insts)
	if err == nil && insts.Instances != nil && len(insts.Instances) > 0 { // instances
		insts.FileDir = fileUtils.GetAbsDir(resFile)
		valueMap = s.ResInstancesService.GetResFromInstances(insts)

	} else {
		ranges := model.ResRanges{}
		err = yaml.Unmarshal(yamlContent, &ranges)

		if err == nil && ranges.Ranges != nil && len(ranges.Ranges) > 0 { // ranges
			valueMap = s.ResRangesService.GetResFromRanges(ranges)

		} else {
			configRes := model.DefField{}
			err = yaml.Unmarshal(yamlContent, &configRes)
			if err == nil { // config
				valueMap = s.ResConfigService.GetResForConfig(configRes)

			}
		}
	}

	vari.GlobalVars.CacheResFileToMap[resFile] = valueMap

	return
}
