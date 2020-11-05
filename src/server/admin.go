package server

import (
	"encoding/json"
	"io"
	"net/http"
)

func AdminHandler(writer http.ResponseWriter, req *http.Request) {
	setupCORS(&writer, req)

	data := map[string]interface{}{}
	ParserJsonParams(req, &data)

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
