package model

import (
	constant "github.com/easysoft/zendata/src/utils/const"
)

var (
	CommonPrefix = "zd_"
	Models = []interface{}{ &ZdDef{}, &ZdField{}, &ZdSection{}, &ZdRefer{},
		&ZdRanges{}, &ZdRangesItem{}, &ZdText{}, &ZdConfig{}, &ZdInstances{}, &ZdInstancesItem{}, &ZdExcel{} }
)

type ZdDef struct {
	Model
	Title  string `gorm:"column:title" json:"title"`
	Path   string `gorm:"column:path" json:"path" yaml:"-"`
	Type   string `gorm:"column:type" json:"type"`
	Desc   string `gorm:"column:desc" json:"desc"`
	Yaml   string `gorm:"yaml" json:"yaml"`
	Folder string `gorm:"-" json:"folder" yaml:"-"`
}
func (*ZdDef) TableName() string {
	return constant.TablePrefix + "def"
}

type ZdField struct {
	Model
	DefID uint `gorm:"column:defID" json:"defID"`
	ParentID uint `gorm:"column:parentID" json:"parentID"`
	Field string `gorm:"column:field" json:"field"`
	Note string `gorm:"column:note" json:"note"`

	Range string `gorm:"column:range" json:"range"`
	Exp  string `gorm:"column:exp" json:"exp"`
	Prefix string `gorm:"column:prefix" json:"prefix"`
	Postfix string `gorm:"column:postfix" json:"postfix"`
	Loop string `gorm:"column:loop" json:"loop"`
	Loopfix string `gorm:"column:loopfix" json:"loopfix"`
	Format string `gorm:"column:format" json:"format"`
	Type string `gorm:"column:type" json:"type"`
	Mode string `gorm:"column:mode" json:"mode"`
	Length int `gorm:"column:length" json:"length"`
	LeftPad string `gorm:"column:leftPad" json:"leftPad"`
	RightPad string `gorm:"column:rightPad" json:"rightPad"`
	Rand bool `gorm:"column:rand" json:"rand"`

	ConfigID	uint `gorm:"column:configID" json:"configID"`
	Config	string `gorm:"column:config" json:"config"`
	Use	string `gorm:"column:use" json:"use"`
	UseID	uint `gorm:"column:useID" json:"useID"`
	From	string `gorm:"column:fromCol" json:"fromCol"`
	Select	string `gorm:"column:selectCol" json:"selectCol"`
	Where	string `gorm:"column:whereCol" json:"whereCol"`
	Limit	int `gorm:"column:limitCol" json:"limitCol"`

	Ord int `gorm:"column:ord;default:1" json:"ord"`
	Children []*ZdField `gorm:"-" json:"children"`
	Froms []*ZdField    `gorm:"-" json:"froms"`

	// for range edit
	IsRange bool         `gorm:"column:isRange;default:true" json:"isRange"`
	Sections []ZdSection `gorm:"ForeignKey:fieldID" json:"sections"`

	// for refer edit
	Refer ZdRefer `gorm:"ForeignKey:fieldID" json:"refer"`
}
func (*ZdField) TableName() string {
	return constant.TablePrefix + "field"
}

type ZdRefer struct {
	Model
	FieldID uint   `gorm:"column:fieldID" json:"fieldID"`
	Type    string `gorm:"column:type" json:"type"`
	File   string `gorm:"column:file" json:"file"`
	ColName     string    `gorm:"column:colName" json:"colName"`
	ColIndex     int    `gorm:"column:colIndex" json:"colIndex"`
	Count int    `gorm:"column:count" json:"count"`
	HasTitle bool `gorm:"column:hasTitle" json:"hasTitle"`
}
func (*ZdRefer) TableName() string {
	return constant.TablePrefix + "refer"
}

type ZdSection struct {
	Model
	FieldID uint   `gorm:"column:fieldID" json:"fieldID"`
	Type    string `gorm:"column:type;default:scope" json:"type"`
	Value     string `gorm:"column:value" json:"value"`
	Ord     int    `gorm:"column:ord;default:1" json:"ord"`

	// for range
	Start string `gorm:"column:start" json:"start"`
	End string `gorm:"column:end" json:"end"`
	Step string `gorm:"column:step;default:1" json:"step"`
	Repeat string `gorm:"column:repeat;default:1" json:"repeat"`
	Rand bool `gorm:"column:rand;default:false" json:"rand"`

	// for arr and const
	Text string `gorm:"-" json:"-"`
}
func (*ZdSection) TableName() string {
	return constant.TablePrefix + "section"
}

type ZdRanges struct {
	Model
	Title  string `gorm:"column:title" json:"title"`
	Name  string `gorm:"column:name" json:"name"`
	Desc   string `gorm:"column:desc" json:"desc"`
	Path   string `gorm:"column:path" json:"path" yaml:"-"`

	Yaml   string `gorm:"yaml" json:"yaml"`
	Folder string `gorm:"-" json:"folder" yaml:"-"`

	Field string `gorm:"column:field" json:"field"`
	Note string `gorm:"column:note" json:"note"`

	Prefix string `gorm:"column:prefix" json:"prefix"`
	Postfix string `gorm:"column:postfix" json:"postfix"`
	Format string `gorm:"column:format" json:"format"`

	Ranges []ZdRangesItem `gorm:"ForeignKey:rangesID" json:"ranges"`
}
func (*ZdRanges) TableName() string {
	return constant.TablePrefix + "ranges"
}

type ZdRangesItem struct {
	Model
	Name string `gorm:"column:name" json:"name"`
	Value string `gorm:"column:value" json:"value"`
	RangesID uint `gorm:"column:rangesID" json:"rangesID"`
	Ord int `gorm:"column:ord" json:"ord"`

	// for tree node
	ParentID uint `gorm:"-" json:"parentID"`
	Children []*ZdRangesItem `gorm:"-" json:"children"`
}
func (*ZdRangesItem) TableName() string {
	return constant.TablePrefix + "rangesItem"
}

type ZdText struct {
	Model
	Title  string `gorm:"column:title" json:"title"`
	Name  string `gorm:"column:name" json:"name"`
	Content   string `gorm:"column:content" json:"content"`
	Path   string `gorm:"column:path" json:"path" yaml:"-"`

	Folder string `gorm:"-" json:"folder" yaml:"-"`
}
func (*ZdText) TableName() string {
	return constant.TablePrefix + "text"
}

type ZdExcel struct {
	Model
	Title  string `gorm:"column:title" json:"title"`
	Name  string `gorm:"column:name" json:"name"`
	Text   string `gorm:"column:desc" json:"desc"`
	Path   string `gorm:"column:path" json:"path" yaml:"-"`

	Folder string `gorm:"-" json:"folder" yaml:"-"`
}
func (*ZdExcel) TableName() string {
	return constant.TablePrefix + "excel"
}

type ZdConfig struct {
	Model
	Title  string `gorm:"column:title" json:"title"`
	Name  string `gorm:"column:name" json:"name"`
	Desc   string `gorm:"column:desc" json:"desc"`
	Path   string `gorm:"column:path" json:"path" yaml:"-"`

	Yaml   string `gorm:"yaml" json:"yaml"`
	Folder string `gorm:"-" json:"folder" yaml:"-"`
	Field string `gorm:"column:field" json:"field"`
	Note string `gorm:"column:note" json:"note"`

	Prefix string `gorm:"column:prefix" json:"prefix"`
	Postfix string `gorm:"column:postfix" json:"postfix"`
	Loop string `gorm:"column:loop" json:"loop"`
	Loopfix string `gorm:"column:loopfix" json:"loopfix"`
	Format string `gorm:"column:format" json:"format"`
}
func (*ZdConfig) TableName() string {
	return constant.TablePrefix + "config"
}

type ZdInstances struct {
	Model
	Title  string `gorm:"column:title" json:"title"`
	Name  string `gorm:"column:name" json:"name"`
	Desc   string `gorm:"column:desc" json:"desc"`
	Path   string `gorm:"column:path" json:"path" yaml:"-"`

	Yaml   string `gorm:"yaml" json:"yaml"`
	Folder string `gorm:"-" json:"folder" yaml:"-"`

	Field string `gorm:"column:field" json:"field"`

	Instances []ZdInstancesItem `gorm:"ForeignKey:instancesID" json:"instances"`
}
func (*ZdInstances) TableName() string {
	return constant.TablePrefix + "instances"
}

type ZdInstancesItem struct {
	Model
	DefID uint `gorm:"column:defID" json:"defID"`
	ParentID uint `gorm:"column:parentID" json:"parentID"`
	Field string `gorm:"column:field" json:"field"`
	Note string `gorm:"column:note" json:"note"`

	Range string `gorm:"column:range" json:"range"`
	Exp  string `gorm:"column:exp" json:"exp"`
	Prefix string `gorm:"column:prefix" json:"prefix"`
	Postfix string `gorm:"column:postfix" json:"postfix"`
	Loop string `gorm:"column:loop" json:"loop"`
	Loopfix string `gorm:"column:loopfix" json:"loopfix"`
	Format string `gorm:"column:format" json:"format"`
	Type string `gorm:"column:type" json:"type"`
	Mode string `gorm:"column:mode" json:"mode"`
	Length int `gorm:"column:length" json:"length"`
	LeftPad string `gorm:"column:leftPad" json:"leftPad"`
	RightPad string `gorm:"column:rightPad" json:"rightPad"`
	Rand bool `gorm:"column:rand" json:"rand"`

	ConfigID	uint `gorm:"column:configID" json:"configID"`
	Config	string `gorm:"column:config" json:"config"`
	Use	string `gorm:"column:use" json:"use"`
	UseID	uint `gorm:"column:useID" json:"useID"`
	From	string `gorm:"column:fromCol" json:"fromCol"`
	Select	string `gorm:"column:selectCol" json:"selectCol"`
	Where	string `gorm:"column:whereCol" json:"whereCol"`
	Limit	int `gorm:"column:limitCol" json:"limitCol"`

	Ord int `gorm:"column:ord;default:1" json:"ord"`
	Children []*ZdInstancesItem `gorm:"-" json:"children"`
	Froms []*ZdInstancesItem    `gorm:"-" json:"froms"`

	// for range edit
	IsRange bool         `gorm:"column:isRange;default:true" json:"isRange"`
	Sections []ZdSection `gorm:"ForeignKey:fieldID" json:"sections"`

	// for refer edit
	Refer ZdRefer `gorm:"ForeignKey:fieldID" json:"refer"`
}
func (*ZdInstancesItem) TableName() string {
	return constant.TablePrefix + "instancesItem"
}
