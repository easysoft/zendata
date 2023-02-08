package model

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
)

type ZdMock struct {
	Model
	Name string `json:"name"`
	Desc string `json:"desc"`

	SpecContent string `json:"specContent"`
	MockContent string `json:"mockContent"`
	DataContent string `json:"dataContent"`
}

func (*ZdMock) TableName() string {
	return consts.TablePrefix + "mock"
}
