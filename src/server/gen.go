package server

import (
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	"github.com/easysoft/zendata/src/utils/vari"
	"net/http"
	"strconv"
)

func ParseGenParams(req *http.Request) (defaultFile, configFile, fields string, count int,
		format, table string, decode bool, input, output string) {
	query := req.URL.Query()

	defaultFile = ParserGetParams(query,"default", "d")
	configFile = ParserGetParams(query,"config", "c")
	countStr := ParserGetParams(query,"lines", "n")
	if countStr == "" {
		countStr = "10"
	}

	fields = ParserGetParams(query,"field", "F")

	format = constant.FormatJson
	table = ""

	if req.Method == http.MethodPost {
		req.ParseForm()

		countStr = ParserPostParams(req, "lines", "n", countStr, false)
		defaultDefContent := ParserPostParams(req, "default", "d", "", true)
		configDefContent := ParserPostParams(req, "config", "c", "", true)

		if defaultDefContent != "" {
			defaultFile = vari.WorkDir + "._default.yaml"
			fileUtils.WriteFile(defaultFile, defaultDefContent)
		}
		if configDefContent != "" {
			configFile = vari.WorkDir + "._config.yaml"
			fileUtils.WriteFile(configFile, configDefContent)
		}
	}

	countFromPForm, err := strconv.Atoi(countStr)
	if err == nil {
		count = countFromPForm
	}

	return
}
