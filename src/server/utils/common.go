package serverUtils

import (
	"encoding/json"
	"github.com/easysoft/zendata/src/model"
)

func ConvertDef(data interface{}) (def model.ZdDef) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &def)

	return
}
func ConvertField(data interface{}) (field model.ZdField) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &field)

	return
}

func ConvertSection(data interface{}) (section model.ZdSection) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &section)

	return
}

func ConvertRefer(data interface{}) (refer model.ZdRefer) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &refer)

	return
}

func ConvertResFile(data interface{}) (refer model.ResFile) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &refer)

	return
}
func ConvertRanges(data interface{}) (ranges model.ZdRanges) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &ranges)

	return
}
func ConvertRangesItem(data interface{})  (item model.ZdRangesItem) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &item)

	return
}

func ConvertExcel(data interface{}) (ranges model.ZdExcel) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &ranges)

	return
}
func ConvertInstances(data interface{}) (ranges model.ZdInstances) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &ranges)

	return
}
func ConvertConfig(data interface{}) (ranges model.ZdConfig) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &ranges)

	return
}
func ConvertText(data interface{}) (ranges model.ZdText) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &ranges)

	return
}

func ConvertParams(data interface{}) (mp map[string]string) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &mp)

	return
}
