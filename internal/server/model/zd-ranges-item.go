package model

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
)

type ZdRangesItem struct {
	BaseModel

	RangesID uint   `json:"rangesID"`
	Field    string `json:"field"`
	Ord      int    `json:"ord"`

	Value    string      `json:"value"`
	Sections []ZdSection `gorm:"ForeignKey:ownerID" json:"sections"`

	// for tree node
	ParentID uint            `gorm:"-" json:"parentID"`
	Fields   []*ZdRangesItem `gorm:"-" json:"fields"`
}

func (*ZdRangesItem) TableName() string {
	return consts.TablePrefix + "ranges_item"
}
