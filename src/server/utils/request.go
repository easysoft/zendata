package serverUtils

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

func SetupCORS(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func OutputErr(err error, writer http.ResponseWriter) {
	errRes := ErrRes(err.Error())
	WriteRes(errRes, writer)
}

func WriteRes(ret model.ResData, writer http.ResponseWriter) {
	jsonStr, _ := json.Marshal(ret)
	io.WriteString(writer, string(jsonStr))
}

func ErrRes(msg string) model.ResData {
	return model.ResData{Code: 0, Msg: msg}
}

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
			val = item[len(item)-1]
		}
	}
	return val
}

func ParserPostParams(req *http.Request, paramName1, paramName2 string, dft string, isFile bool) (ret string) {
	if paramName2 != "" && req.FormValue(paramName2) != "" {
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
func ParserPostParams2(req *http.Request, paramName1, paramName2 string, dft string, isFile bool) (ret []byte) {
	ret1 := ""
	if paramName2 != "" && req.FormValue(paramName2) != "" {
		ret1 = req.FormValue(paramName2)
	} else if paramName1 != "" && req.FormValue(paramName1) != "" { // high priority than paramName2
		ret1 = req.FormValue(paramName1)
	}
	if isFile && ret1 == "" {
		postFile, _, _ := req.FormFile(paramName2)
		if postFile != nil {
			defer postFile.Close()

			buf := bytes.NewBuffer(nil)
			io.Copy(buf, postFile)
			ret = buf.Bytes()
		}

		if ret == nil {
			postFile, _, _ = req.FormFile(paramName1)
			if postFile != nil {
				defer postFile.Close()

				buf := bytes.NewBuffer(nil)
				io.Copy(buf, postFile)
				ret = buf.Bytes()
			}
		}
	}

	if ret == nil {
		ret = []byte(dft)
	}

	return
}
