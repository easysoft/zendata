package model

import (
	"time"
)

type ReqData struct {
	Action string `json:"action"`
	Id int  `json:"id"`
	DomainId int  `json:"domainId"`
	Mode string `json:"mode"`
	Data interface{} `json:"data"`

	Src int `json:"src"`
	Dist int `json:"dist"`
}

type ResData struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
	Model interface{} `json:"model"`
	Res   interface{} `json:"res"`
}

type Model struct {
	ID        uint      `gorm:"column:id;primary_key" json:"id"`
	CreatedAt time.Time `gorm:"column:createTime" json:"createTime" yaml:"-"`
	UpdatedAt time.Time `gorm:"column:updateTime" json:"updateTime" yaml:"-"`

	Disabled bool `gorm:"column:disabled;default:false" json:"disabled" yaml:"-"`
	Deleted  bool `gorm:"column:deleted;default:false" json:"deleted" yaml:"-"`
}
