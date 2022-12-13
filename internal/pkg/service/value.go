package service

import (
	constant "github.com/easysoft/zendata/internal/pkg/const"
	valueGen "github.com/easysoft/zendata/internal/pkg/gen/value"
	"github.com/easysoft/zendata/internal/pkg/model"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/oklog/ulid/v2"
	"math/rand"
	"strings"
	"time"
)

type ValueService struct {
}

func NewValueService() *ValueService {
	return &ValueService{}
}

func (s *ValueService) CreateTimestampField(field *model.DefField) {
	valueGen.ConvertTmFormat(field)

	rang := strings.Trim(strings.TrimSpace(field.Range), ",")
	rangeSections := strings.Split(rang, ",")

	values := make([]interface{}, 0)
	for _, section := range rangeSections {
		valueGen.CreateTimestampSectionValue(section, &values)
	}

	if len(values) == 0 {
		values = append(values, "N/A")
	}

	field.Values = values
}

func (s *ValueService) CreateUlidField(field *model.DefField) {
	count := 0

	t := time.Unix(1000000, 0)
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)

	for true {
		val := ulid.MustNew(ulid.Timestamp(t), entropy).String()
		field.Values = append(field.Values, val)

		count++
		if count >= constant.MaxNumb || count > vari.GenVars.Total {
			break
		}
	}
}
