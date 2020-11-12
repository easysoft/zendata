package server

import (
	"encoding/json"
	"github.com/easysoft/zendata/src/model"
	"io"
	"net/http"
)

func setupCORS(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func outputErr(err error, writer http.ResponseWriter) {
	errRes := errRes(err.Error())
	writeRes(errRes, writer)
}

func writeRes(ret model.ResData, writer http.ResponseWriter) {
	jsonStr, _ := json.Marshal(ret)
	io.WriteString(writer, string(jsonStr))
}

func errRes(msg string) model.ResData {
	return model.ResData{ Code: 0, Msg: msg }
}

func convertDef(data interface{}) (def model.Def) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &def)

	return
}
func convertField(data interface{}) (field model.Field) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &field)

	return
}

func convertSection(data interface{}) (section model.Section) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &section)

	return
}

func convertParams(data interface{}) (mp map[string]string) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &mp)

	return
}
