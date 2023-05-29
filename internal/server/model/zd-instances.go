package model

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
)

type ZdInstances struct {
	BaseModel

	Title string `json:"title" yaml:"title,omitempty"`
	Desc  string `json:"desc" yaml:"desc,omitempty"`

	Yaml   string `json:"yaml" yaml:"-"`
	Path   string `json:"path" yaml:"-"`
	Folder string `json:"folder" yaml:"-"`

	FileName  string `json:"fileName" yaml:"-"`
	ReferName string `json:"referName" yaml:"-"`

	From      string            `gorm:"-" json:"from"`
	Instances []ZdInstancesItem `gorm:"ForeignKey:instancesID" json:"instances" yaml:"instances"`
}

func (*ZdInstances) TableName() string {
	return consts.TablePrefix + "instances"
}
