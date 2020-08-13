package service

import (
	"bytes"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	"github.com/easysoft/zendata/src/utils/vari"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

func ParseRequestParams(req *http.Request) (defaultFile, configFile, fields string, count int,
		format, table string, decode bool, input, output string) {
	query := req.URL.Query()

	defaultFile = GetRequestParams(query,"default", "d")
	configFile = GetRequestParams(query,"config", "c")
	countStr := GetRequestParams(query,"lines", "n")
	if countStr == "" {
		countStr = "10"
	}

	fields = GetRequestParams(query,"field", "F")

	format = constant.FormatJson
	table = ""

	if req.Method == http.MethodPost {
		req.ParseForm()

		countStr = GetPostParams(req, "lines", "n", countStr, false)
		defaultDefContent := GetPostParams(req, "default", "d", "", true)
		configDefContent := GetPostParams(req, "config", "c", "", true)

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

func GetRequestParams(values url.Values, name, short string) (val string) {
	for key, item := range values {
		if key == name || key == short {
			val = item[len(item) - 1]
		}
	}
	return val
}

func GetPostParams(req *http.Request, paramName1, paramName2 string, dft string, isFile bool) (ret string) {
	if paramName2 != "" &&  req.FormValue(paramName2) != "" {
		ret = req.FormValue(paramName2)
	} else if paramName1 != "" && req.FormValue(paramName1) != "" { // high priority than paramName2
		ret = req.FormValue(paramName1)
	}

	if isFile && ret == "" {
		postFile, _, _ := req.FormFile(paramName2)
		if postFile != nil {
			defer postFile.Close()

			buf := bytes.NewBuffer(nil)
			io.Copy(buf, postFile)
			ret = buf.String()
		}

		if ret == "" {
			postFile, _, _ = req.FormFile(paramName1)
			if postFile != nil {
				defer postFile.Close()

				buf := bytes.NewBuffer(nil)
				io.Copy(buf, postFile)
				ret = buf.String()
			}
		}
	}

	if ret == "" {
		ret = dft
	}

	return
}