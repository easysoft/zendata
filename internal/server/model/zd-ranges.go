package model

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
)

type ZdRanges struct {
	BaseModel

	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Prefix  string `json:"prefix"`
	Postfix string `json:"postfix"`
	Format  string `json:"format"`

	Yaml      string `json:"yaml"`
	Path      string `json:"path" yaml:"-"`
	Folder    string `json:"folder" yaml:"-"`
	FileName  string `json:"fileName" yaml:"-"`
	ReferName string `json:"referName" yaml:"-"`

	Ranges   []ZdRangesItem    `gorm:"ForeignKey:rangesID" json:"ranges" yaml:"-"`
	RangeMap map[string]string `gorm:"-" yaml:"ranges"`
}

func (*ZdRanges) TableName() string {
	return consts.TablePrefix + "ranges"
}
