package service

import (
	"fmt"
	"github.com/easysoft/zendata/internal/pkg/domain"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
)

func (s *OutputService) GenText(returnedForYamlRefer bool) (lines []interface{}) {
	simulatedFieldFromDef := domain.DefField{
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
			logUtils.PrintLine(fmt.Sprintf("%v", line) + "\n")
		}
	}

	return
}
