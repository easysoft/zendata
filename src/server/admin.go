package server

import (
	"encoding/json"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"io"
	"net/http"
)

func AdminHandler(writer http.ResponseWriter, req *http.Request) {
	setupCORS(&writer, req)
	logUtils.PrintToScreen("111")

	ret := new(Ret)
	id := req.FormValue("id")
	//id := req.PostFormValue('id')

	ret.Code = 0
	ret.Param = id
	ret.Msg = "success"
	ret.Data = map[string]interface{}{"key": "value"}
	ret_json,_ := json.Marshal(ret)

	io.WriteString(writer, string(ret_json))
}

type Ret struct{
	Code int
	Param string
	Msg string
	Data interface{}
}
