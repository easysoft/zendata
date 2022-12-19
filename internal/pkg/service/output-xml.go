package service

import (
	"fmt"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
)

func (s *OutputService) GenXml() {
	records := s.GenObjs()

	s.PrintXmlHeader()

	for i, record := range records {
		line := s.getXmlLine(i, record)

		logUtils.PrintRecord(line + "\n")
	}

	s.PrintXmlFooter()

	return
}

func (s *OutputService) PrintXmlHeader() {
	logUtils.PrintLine("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<testdata>\n  <title>Test Data</title>\n")

}
func (s *OutputService) PrintXmlFooter() {
	logUtils.PrintLine("</testdata>\n")
}

func (s *OutputService) getXmlLine(index int, record map[string]interface{}) string {
	str := ""
	j := 0
	for key, val := range record {
		str += fmt.Sprintf("    <%s>%s</%s>", key, val, key)
		if j != len(record)-1 {
			str = str + "\n"
		}

		j++
	}

	text := fmt.Sprintf("  <row>\n%s\n  </row>", str)

	return text
}
