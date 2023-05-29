package model

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
)

type ZdInstancesItem struct {
	BaseModel

	Instance string `json:"instance" yaml:"instance,omitempty"`
	Note     string `json:"note" yaml:"note,omitempty"`

	Field    string `json:"field" yaml:"field,omitempty"`
	Range    string `json:"range" yaml:"range,omitempty"`
	Prefix   string `json:"prefix" yaml:"prefix,omitempty"`
	Postfix  string `json:"postfix" yaml:"postfix,omitempty"`
	Loop     string `json:"loop" yaml:"loop,omitempty"`
	Loopfix  string `json:"loopfix" yaml:"loopfix,omitempty"`
	Format   string `json:"format" yaml:"format,omitempty"`
	Type     string `json:"type" yaml:"type,omitempty"`
	Mode     string `json:"mode" yaml:"mode,omitempty"`
	Items    int    `json:"records,omitempty" yaml:"records,omitempty"`
	Length   int    `json:"length" yaml:"length,omitempty"`
	LeftPad  string `json:"leftPad" yaml:"leftPad,omitempty"`
	RightPad string `json:"rightPad" yaml:"rightPad,omitempty"`
	Rand     bool   `json:"rand" yaml:"rand,omitempty"`

	Config string `json:"config" yaml:"config,omitempty"`
	Use    string `json:"use" yaml:"use,omitempty"`
	From   string `json:"fromCol" yaml:"from,omitempty"`
	Select string `json:"selectCol" yaml:"select,omitempty"`
	Where  string `json:"whereCol" yaml:"where,omitempty"`
	Limit  int    `json:"limitCol" yaml:"limit,omitempty"`

	Exp         string `json:"exp" yaml:"exp,omitempty"`
	InstancesID uint   `json:"instancesID" yaml:"-"`
	ParentID    uint   `json:"parentID" yaml:"-"`
	ConfigID    uint   `json:"configID" yaml:"-"`
	UseID       uint   `json:"useID" yaml:"-"`

	Ord    int                `gorm:"default:1" json:"ord" yaml:"-"`
	Fields []*ZdInstancesItem `gorm:"-" json:"fields" yaml:"fields,omitempty"`
	Froms  []*ZdInstancesItem `gorm:"-" json:"froms" yaml:"froms,omitempty"`

	// for range edit
	IsRange  bool        `gorm:"default:true" json:"isRange" yaml:"-"`
	Sections []ZdSection `gorm:"ForeignKey:ownerID" json:"sections" yaml:"-"`

	// for refer edit
	Refer ZdRefer `gorm:"ForeignKey:ownerID" json:"refer" yaml:"-"`
}

func (*ZdInstancesItem) TableName() string {
	return consts.TablePrefix + "instances_item"
}
