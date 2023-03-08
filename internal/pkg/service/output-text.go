package service

import (
	"fmt"
	"github.com/easysoft/zendata/internal/pkg/domain"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
)

func (s *OutputService) FindValuesByPath(path string, defFields domain.DefField) []interface{} {
	if path == defFields.Path {
		return defFields.Values
	}
	if len(defFields.Fields) == 0 {
		return nil
	}

	for _, item := range defFields.Fields {
		values := s.FindValuesByPath(path, item)
		if values != nil {
			return values
		}
	}
	return nil
}

func (s *OutputService) GenText(returnedForYamlRefer bool) (lines []interface{}) {
	simulatedFieldFromDef := domain.DefField{
		Fields: vari.GlobalVars.DefData.Fields,
		Join:   true,
	}
	s.CombineService.CombineChildrenIfNeeded(&simulatedFieldFromDef, true)

	if !returnedForYamlRefer {
		s.PrintHumanHeaderIfNeeded()
	}

	if vari.GlobalVars.ExportChildField != "" {
		values := s.FindValuesByPath(vari.GlobalVars.ExportChildField, simulatedFieldFromDef)
		for _, item := range values {
			line := s.PlaceholderService.ReplacePlaceholder(item.(string))
			if returnedForYamlRefer {
				lines = append(lines, line)
			} else {
				logUtils.PrintLine(fmt.Sprintf("%v", line) + "\n")
			}
		}
		return
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
