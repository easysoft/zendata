package service

import (
	"github.com/easysoft/zendata/internal/pkg/model"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
)

func (s *OutputService) GenRows(def *model.DefData) {
	simulatedFieldFromDef := model.DefField{
		Fields: def.Fields,
		Union:  true,
	}

	s.CombineService.CombineChildrenIfNeeded(&simulatedFieldFromDef)

	for _, item := range simulatedFieldFromDef.Values {
		line := s.PlaceholderService.ReplacePlaceholder(item.(string))

		logUtils.PrintLine(line + "\n")
	}

	return
}
