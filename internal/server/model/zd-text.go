package model

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
)

type ZdText struct {
	BaseModel

	Title string `json:"title"`

	Content   string `json:"content"`
	Path      string `json:"path" yaml:"-"`
	Folder    string `json:"folder" yaml:"-"`
	FileName  string `json:"fileName" yaml:"-"`
	ReferName string `json:"referName" yaml:"-"`
}

func (*ZdText) TableName() string {
	return consts.TablePrefix + "text"
}
