package service

import (
	"strings"

	"github.com/easysoft/zendata/internal/pkg/domain"
	"github.com/easysoft/zendata/internal/pkg/helper"
)

type ValueService struct {
}

func (s *ValueService) CreateTimestampField(field *domain.DefField) {
	helper.ConvertTmFormat(field)

	rang := strings.Trim(strings.TrimSpace(field.Range), ",")
	rangeSections := strings.Split(rang, ",")

	values := make([]interface{}, 0)
	for _, section := range rangeSections {
		helper.CreateTimestampSectionValue(section, &values)
	}

	if len(values) == 0 {
		values = append(values, "N/A")
	}

	field.Values = values
}
