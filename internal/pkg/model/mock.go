package model

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
)

type ZdMock struct {
	Model
	Name  string `json:"name"`
	Path  string `json:"path"`
	Desc  string `json:"desc"`
	DefId uint   `json:"defId"`

	SpecContent string `json:"specContent"`
	MockContent string `json:"mockContent"`
	DataContent string `json:"dataContent"`
	DataPath    string `json:"dataPath"`
}

func (*ZdMock) TableName() string {
	return consts.TablePrefix + "mock"
}

type ZdMockSampleSrc struct {
	Model
	MockId uint   `json:"mockId"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}

func (*ZdMockSampleSrc) TableName() string {
	return consts.TablePrefix + "mock_sample_src"
}
