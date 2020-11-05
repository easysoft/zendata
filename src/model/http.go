package model

type ReqData struct {
	Action string `json:"action"`
	Data interface{} `json:"data"`
}

type ResData struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}
