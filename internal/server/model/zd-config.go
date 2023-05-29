package model

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
)

type ZdConfig struct {
	BaseModel

	Title string `json:"title"`
	Desc  string `json:"desc"`

	Range   string `json:"range"`
	Prefix  string `json:"prefix"`
	Postfix string `json:"postfix"`
	Loop    string `json:"loop"`
	Loopfix string `json:"loopfix"`
	Format  string `json:"format"`

	Yaml      string `json:"yaml"`
	Path      string `json:"path" yaml:"-"`
	Folder    string `json:"folder" yaml:"-"`
	FileName  string `json:"fileName" yaml:"-"`
	ReferName string `json:"referName" yaml:"-"`

	// for range edit
	IsRange  bool        `gorm:"default:true" json:"isRange" yaml:"-"`
	Sections []ZdSection `gorm:"ForeignKey:ownerID" json:"sections" yaml:"-"`
}

func (*ZdConfig) TableName() string {
	return consts.TablePrefix + "config"
}
