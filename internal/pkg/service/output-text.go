package service

import (
	"github.com/easysoft/zendata/internal/pkg/model"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
)

func (s *OutputService) GenText(returnedForYamlRefer bool) (lines []interface{}) {
	simulatedFieldFromDef := model.DefField{
		Fields: vari.GlobalVars.DefData.Fields,
		Join:   true,
	}

	s.CombineService.CombineChildrenIfNeeded(&simulatedFieldFromDef, true)

	if !returnedForYamlRefer {
		s.PrintHumanHeaderIfNeeded()
	}

	for _, item := range simulatedFieldFromDef.Values {
		line := s.PlaceholderService.ReplacePlaceholder(item.(string))

		if returnedForYamlRefer {
			lines = append(lines, line)
		} else {
			logUtils.PrintLine(line + "\n")
		}
	}

	return
}
