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
	seq string `json:"seq"`
	name string `json:"name"`
	path string `json:"path"`
	desc string `json:"desc"`
}

func (*Data) TableName() string {
	return constant.TablePrefix + "data"
}
