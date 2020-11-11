package server

import (
	"encoding/json"
	"github.com/easysoft/zendata/src/model"
	defServer "github.com/easysoft/zendata/src/server/def"
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
		ret.Data, err = defServer.List()
	} else if reqData.Action == "getDef" {
		var def model.Def
		def, err = defServer.Get(reqData.Id)

		ret.Data = def
	} else if reqData.Action == "saveDef" {
		def := convertDef(reqData.Data)

		if def.ID == 0 {
			err = defServer.Create(&def)
		} else {
			err = defServer.Update(&def)
		}

		ret.Data = def
	} else if reqData.Action == "removeDef" {
		err = defServer.Remove(reqData.Id)
	} else if reqData.Action == "getDefFieldTree" {
		ret.Data, err = defServer.GetDefFieldTree(reqData.Id)


	} else if reqData.Action == "getDefField" {
		ret.Data, err = defServer.GetDefField(reqData.Id)
	} else if reqData.Action == "createDefField" {
		var field *model.Field
		field, err = defServer.CreateDefField(0, uint(reqData.Id), "新字段", reqData.Mode)
		ret.Data, err = defServer.GetDefFieldTree(int(field.DefID))
		ret.Field = field
	} else if reqData.Action == "saveDefField" {
		field := convertField(reqData.Data)
		err = defServer.SaveDefField(&field)
	}

	if err != nil {
		ret.Code = 0
	}

	bytes, _ = json.Marshal(ret)
	io.WriteString(writer, string(bytes))
}
