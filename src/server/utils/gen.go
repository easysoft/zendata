package serverUtils

import (
	constant "github.com/easysoft/zendata/src/utils/const"
	"net/http"
	"strconv"
	"strings"
)

func ParseGenParams(req *http.Request) (defaultFile, configFile, fields string, count int,
	format string, trim bool, table string, decode bool, input, output string) {
	query := req.URL.Query()

	defaultFile = ParserGetParams(query, "default", "d")
	configFile = ParserGetParams(query, "config", "c")
	trimStr := ParserGetParams(query, "trim", "T")
	countStr := ParserGetParams(query, "lines", "n")
	if countStr == "" {
		countStr = "10"
	}

	fields = ParserGetParams(query, "field", "F")

	format = constant.FormatJson
	table = ""

	trimStr = strings.ToLower(strings.TrimSpace(trimStr))
	if trimStr == "t" || trimStr == "true" {
		trim = true
	}

	countFromForm, err := strconv.Atoi(countStr)
	if err == nil {
		count = countFromForm
	}

	return
}

func ParseGenParamsToByte(req *http.Request) (defaultDefContent, configDefContent []byte, fields string, count int,
	format string, trim bool, table string, decode bool, input, output string) {

	format = constant.FormatJson
	table = ""
	req.ParseForm()

	defaultDefContent = []byte(ParserPostParams(req, "default", "d", "", true))
	configDefContent = []byte(ParserPostParams(req, "config", "c", "", true))
	trimStr := ParserPostParams(req, "trim", "T", "", false)
	countStr := ParserPostParams(req, "lines", "n", "", false)

	if countStr == "" {
		countStr = "10"
	}

	trimStr = strings.ToLower(strings.TrimSpace(trimStr))
	if trimStr == "t" || trimStr == "true" {
		trim = true
	}

	countFromForm, err := strconv.Atoi(countStr)
	if err == nil {
		count = countFromForm
	}

	return
}
