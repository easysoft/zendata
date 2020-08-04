package service

import (
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	"github.com/easysoft/zendata/src/utils/vari"
	"net/http"
	"net/url"
	"strconv"
)

func ParseRequestParams(req *http.Request) (defaultFile, yamlFile, fields string, count int,
		human string, format, table string, decode bool, input, output string) {
	query := req.URL.Query()

	defaultFile = GetRequestParams(query,"default", "d")
	yamlFile = GetRequestParams(query,"config", "c")
	countStr := GetRequestParams(query,"lines", "n")
	if countStr == "" {
		countStr = "10"
	}
	count, _ = strconv.Atoi(countStr)
	fields = GetRequestParams(query,"field", "F")

	format = constant.FormatJson
	table = ""

	human = GetRequestParams(query,"human", "H")

	if req.Method == http.MethodPost {
		// save to files
		req.ParseForm()
		defaultDefContent := req.FormValue("default")
		configDefContent := req.FormValue("config")

		if defaultDefContent != "" {
			defaultFile = vari.ExeDir + "._default.yaml"
			fileUtils.WriteFile(defaultFile, defaultDefContent)
		}
		if configDefContent != "" {
			yamlFile = vari.ExeDir + "._config.yaml"
			fileUtils.WriteFile(yamlFile, configDefContent)
		}
	}

	return
}

func GetRequestParams(values url.Values, name, short string) (val string) {
	for key, item := range values {
		if key == name || key == short {
			val = item[len(item) - 1]
		}
	}
	return val
}