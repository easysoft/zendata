package model

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
)

type ZdExcel struct {
	BaseModel

	Title string `json:"title"`
	Sheet string `json:"sheet"`

	ChangeTime string `json:"changeTime"`
	Yaml       string `json:"yaml"`
	Path       string `json:"path" yaml:"-"`
	Folder     string `json:"folder" yaml:"-"`
	FileName   string `json:"fileName" yaml:"-"`
	ReferName  string `json:"referName" yaml:"-"`
}

func (*ZdExcel) TableName() string {
	return consts.TablePrefix + "excel"
}
