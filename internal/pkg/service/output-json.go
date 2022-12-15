package service

import (
	"encoding/json"
	"github.com/easysoft/zendata/internal/pkg/model"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
)

type OutputService struct {
	CombineService *CombineService `inject:""`
}

func (s *OutputService) GenJson(def *model.DefData) {
	records := make([]map[string]interface{}, 0)

	for i := 0; i < vari.GlobalVars.Total; i++ {
		record := map[string]interface{}{}

		for _, field := range def.Fields {
			s.GenFieldMap(&field, &record, i)
		}

		records = append(records, record)
	}

	bytes, err := json.MarshalIndent(records, "", "\t")
	if err != nil {
		logUtils.PrintTo("json marshal failed")
	}

	jsonStr := string(bytes)
	logUtils.PrintBlock(jsonStr)

	return
}

func (s *OutputService) GenFieldMap(field *model.DefField, mp *map[string]interface{}, i int) {
	if field.Union || len(field.Fields) == 0 {
		(*mp)[field.Field] = field.Values[i%len(field.Values)]

	} else {
		childMap := map[string]interface{}{}

		for _, child := range field.Fields {
			s.GenFieldMap(&child, &childMap, i)
		}

		(*mp)[field.Field] = childMap
	}

	return
}

func (s *OutputService) ConnectValues(values []interface{}) (ret string) {
	for _, item := range values {
		ret += item.(string)
	}

	return
}
