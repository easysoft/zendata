package service

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/domain"
	"github.com/easysoft/zendata/internal/pkg/helper"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/oklog/ulid/v2"
	"math/rand"
	"strings"
	"time"
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

func (s *ValueService) CreateUlidField(field *domain.DefField) {
	count := 0

	t := time.Unix(1000000, 0)
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)

	for true {
		val := ulid.MustNew(ulid.Timestamp(t), entropy).String()
		field.Values = append(field.Values, val)

		count++
		if count >= consts.MaxNumb || count > vari.GlobalVars.Total {
			break
		}
	}
}
