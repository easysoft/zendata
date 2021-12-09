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
	query := req.URL.Query()

	trimStr := ParserGetParams(query, "trim", "T")
	countStr := ParserGetParams(query, "lines", "n")
	if countStr == "" {
		countStr = "10"
	}

	fields = ParserGetParams(query, "field", "F")

	format = constant.FormatJson
	table = ""

	req.ParseForm()

	defaultDefContent = ParserPostParamsToByte(req, "default", "d", "", true)
	configDefContent = ParserPostParamsToByte(req, "config", "c", "", true)
	trimStr = string(ParserPostParamsToByte(req, "trim", "T", trimStr, false))
	countStr = string(ParserPostParamsToByte(req, "lines", "n", countStr, false))

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
