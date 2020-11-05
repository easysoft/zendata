package server

import (
	"encoding/json"
	"github.com/easysoft/zendata/src/model"
	"io"
	"net/http"
)

func AdminHandler(writer http.ResponseWriter, req *http.Request) {
	setupCORS(&writer, req)

	reqData := map[string]interface{}{}
	ParserJsonParams(req, &reqData)

	ret := model.ResData{ Code: 1, Msg: "success", Data: reqData }
	jsonStr, _ := json.Marshal(ret)

	io.WriteString(writer, string(jsonStr))
}
