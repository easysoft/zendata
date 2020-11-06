package model

import constant "github.com/easysoft/zendata/src/utils/const"

var CommonPrefix = "zd_"

type ReqData struct {
	Action string `json:"action"`
	Data interface{} `json:"data"`
}

type ResData struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}

type Data struct {
	Seq string `gorm:"column:seq" json:"seq"`
	Name string `gorm:"column:name" json:"name"`
	Path string `gorm:"column:path" json:"path"`
	Desc string `gorm:"column:desc" json:"desc"`
}

func (*Data) TableName() string {
	return constant.TablePrefix + "data"
}
