package service

import (
	"github.com/easysoft/zendata/internal/pkg/gen"
	"github.com/easysoft/zendata/internal/pkg/model"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
)

func (s *OutputService) GenText(def *model.DefData) {
	simulatedFieldFromDef := model.DefField{
		Fields: def.Fields,
		Join:   true,
	}

	s.CombineService.CombineChildrenIfNeeded(&simulatedFieldFromDef)

	gen.PrintHumanHeaderIfNeeded(vari.GlobalVars.ExportFields)

	for _, item := range simulatedFieldFromDef.Values {
		line := s.PlaceholderService.ReplacePlaceholder(item.(string))

		logUtils.PrintLine(line + "\n")
	}

	return
}
