package model

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
)

type ZdRefer struct {
	BaseModel

	OwnerType string `json:"ownerType"` // field or instances
	OwnerID   uint   `json:"ownerID"`
	Type      string `json:"type"`

	Value string `json:"value"`
	File  string `json:"file"`
	Sheet string `json:"sheet"`

	ColName   string `json:"colName"`
	ColIndex  int    `json:"colIndex"`
	Condition string `json:"condition"`
	Count     int    `json:"count"`
	CountTag  string `json:"countTag"`
	Step      int    `json:"step"`
	Rand      bool   `json:"rand"`
	HasTitle  bool   `json:"hasTitle"`
}

func (*ZdRefer) TableName() string {
	return consts.TablePrefix + "refer"
}
