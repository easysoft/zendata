package serverUtils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/easysoft/zendata/internal/pkg/model"
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

func WriteRes(ret model.RespData, writer http.ResponseWriter) {
	jsonStr, _ := json.Marshal(ret)
	io.WriteString(writer, string(jsonStr))
}

func ErrRes(msg string) model.RespData {
	return model.RespData{Code: 0, Msg: msg}
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

func ParserPostParams(req *http.Request, paramName, paramNameShort string, dft string, isFile bool) (ret string) {
	if paramNameShort != "" && req.FormValue(paramNameShort) != "" {
		ret = req.FormValue(paramNameShort)
	} else if paramName != "" && req.FormValue(paramName) != "" { // high priority than paramName2
		ret = req.FormValue(paramName)
	}

	if isFile && ret == "" {
		postFile, _, _ := req.FormFile(paramNameShort)
		if postFile != nil {
			defer postFile.Close()

			buf := bytes.NewBuffer(nil)
			io.Copy(buf, postFile)
			ret = buf.String()
		}

		if ret == "" {
			postFile, _, _ = req.FormFile(paramName)
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
