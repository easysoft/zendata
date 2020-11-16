package serverUtils

import (
	"encoding/json"
	"github.com/easysoft/zendata/src/model"
)

func ConvertDef(data interface{}) (def model.Def) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &def)

	return
}
func ConvertField(data interface{}) (field model.Field) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &field)

	return
}

func ConvertSection(data interface{}) (section model.Section) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &section)

	return
}

func ConvertRefer(data interface{}) (refer model.Refer) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &refer)

	return
}

func ConvertResFile(data interface{}) (refer model.ResFile) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &refer)

	return
}

func ConvertParams(data interface{}) (mp map[string]string) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &mp)

	return
}
