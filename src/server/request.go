package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/easysoft/zendata/src/model"
	"io"
	"log"
	"net/http"
	"net/url"
)

func ParserJsonReq(bytes []byte, obj *model.ReqData) (err error) {
	err = json.Unmarshal(bytes, &obj)
	if err != nil {
		log.Println(fmt.Sprintf("parse json error %s", err))
		return
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
