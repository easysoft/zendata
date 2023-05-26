package model

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
)

type ZdDef struct {
	BaseModel

	Title string `json:"title"`
	Type  string `json:"type"`
	Desc  string `json:"desc"`

	Yaml     string `json:"yaml"`
	Path     string `json:"path" yaml:"-"`
	Folder   string `json:"folder" yaml:"-"`
	FileName string `json:"fileName" yaml:"-"`
	IsMock   bool   `json:"isMock"`

	ReferName string    `json:"referName" yaml:"-"`
	From      string    `json:"from"`
	Fields    []ZdField `json:"fields"`
}

func (ZdDef) TableName() string {
	return consts.TablePrefix + "def"
}
