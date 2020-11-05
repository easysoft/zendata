package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	"github.com/easysoft/zendata/src/utils/vari"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func ParseRequestParams(req *http.Request) (defaultFile, configFile, fields string, count int,
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

func ParserGetParams(values url.Values, name, short string) (val string) {
	for key, item := range values {
		if key == name || key == short {
			val = item[len(item) - 1]
		}
	}
	return val
}

func ParserPostParams(req *http.Request, paramName1, paramName2 string, dft string, isFile bool) (ret string) {
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

func ParserJsonParams(req *http.Request, obj *map[string]interface{}) {
	con, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(con))

	err := req.ParseForm()
	if err != nil {
		log.Fatal("parse form error ",err)
	}

	json.NewDecoder(req.Body).Decode(&obj)
	for key,value := range *obj{
		log.Println("key:",key," => value :",value)
	}
}
