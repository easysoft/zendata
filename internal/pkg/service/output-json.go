package service

import (
	"encoding/json"
	"github.com/easysoft/zendata/internal/pkg/model"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"regexp"
)

func (s *OutputService) GenJson(def *model.DefData) {
	records := s.GenObjs(def)

	s.PrintJsonHeader()

	for i, record := range records {
		bytes, err := json.MarshalIndent(record, "", "\t")
		if err != nil {
			logUtils.PrintTo("json marshal failed")
			break
		}

		jsonStr := string(bytes)

		regx := regexp.MustCompile(`(\n+)`)
		jsonStr = "\t" + regx.ReplaceAllString(jsonStr, "${1}\t")

		postStr := "\n"
		if i < len(records)-1 {
			postStr = "," + postStr
		}

		logUtils.PrintRecord(jsonStr + postStr)
	}

	s.PrintJsonFooter()

	return
}

func (s *OutputService) PrintJsonHeader() {
	logUtils.PrintRecord("[\n")
}
func (s *OutputService) PrintJsonFooter() {
	logUtils.PrintRecord("]\n")
}
