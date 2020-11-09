package model

import (
	constant "github.com/easysoft/zendata/src/utils/const"
)

type Def struct {
	Model
	Name string `gorm:"column:name" json:"name"`
	Path string `gorm:"column:path" json:"path"`
	Desc string `gorm:"column:desc" json:"desc"`
	Folder string `gorm:"-" json:"folder"`
}
func (*Def) TableName() string {
	return constant.TablePrefix + "def"
}

type Field struct {
	Model
	DefID string `gorm:"column:defID" json:"defID"`
	ParentID uint `gorm:"column:parentID" json:"parentID"`
	Field string `gorm:"column:field" json:"field"`
	Note string `gorm:"column:note" json:"note"`
	Prefix string `gorm:"column:prefix" json:"prefix"`
	Postfix string `gorm:"column:postfix" json:"postfix"`
	Ord int `gorm:"column:ord" json:"ord"`

	Children []*Field `gorm:"-" json:"children"`

	Key string `gorm:"-" json:"key"`
	Value string   `gorm:"-" json:"value"`
	Title string   `gorm:"-" json:"title"`
}
func (*Field) TableName() string {
	return constant.TablePrefix + "field"
}
