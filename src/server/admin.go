package server

import (
	"encoding/json"
	"github.com/easysoft/zendata/src/model"
	"io"
	"io/ioutil"
	"net/http"
)

func AdminHandler(writer http.ResponseWriter, req *http.Request) {
	setupCORS(&writer, req)

	bytes, err := ioutil.ReadAll(req.Body)
	if len(bytes) == 0 {
		return
	}

	reqData := model.ReqData{}
	err = ParserJsonReq(bytes, &reqData)
	if err != nil {
		outputErr(err, writer)
		return
	}

	ret := model.ResData{ Code: 1, Msg: "success"}
	if reqData.Action == "listDef" {
		ret.Data = ListData()
	}

	jsonStr, _ := json.Marshal(ret)

	io.WriteString(writer, string(jsonStr))
}
