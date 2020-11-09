package model

import (
	constant "github.com/easysoft/zendata/src/utils/const"
)

type Def struct {
	Name string `gorm:"column:name" json:"name"`
	Path string `gorm:"column:path" json:"path"`
	Desc string `gorm:"column:desc" json:"desc"`
	Folder string `gorm:"-" json:"folder"`
	Model
}

func (*Def) TableName() string {
	return constant.TablePrefix + "def"
}
