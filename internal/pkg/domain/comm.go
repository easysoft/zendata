package domain

// Response
type Response struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type BizError struct {
	Code int64
	Msg  string
}
