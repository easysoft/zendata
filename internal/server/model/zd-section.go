package model

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
)

type ZdSection struct {
	BaseModel

	OwnerType string `json:"ownerType"` // field or instances
	OwnerID   uint   `json:"ownerID"`
	Type      string `gorm:"default:interval" json:"type"`
	Value     string `json:"value"`
	Ord       int    `gorm:"default:1" json:"ord"`

	// for range
	Start     string `json:"start"`
	End       string `json:"end"`
	Step      int    `gorm:"default:1" json:"step"`
	Repeat    string `gorm:"default:1" json:"repeat"`
	RepeatTag string `json:"repeatTag"`
	Rand      bool   `gorm:"default:false" json:"rand"`

	// for arr and const
	Text string `gorm:"-" json:"-"`
}

func (*ZdSection) TableName() string {
	return consts.TablePrefix + "field"
}
