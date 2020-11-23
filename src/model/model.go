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
	Type   string `gorm:"column:type" json:"type"`
	Desc   string `gorm:"column:desc" json:"desc"`

	Yaml   string `gorm:"yaml" json:"yaml"`
	Path   string `gorm:"column:path" json:"path" yaml:"-"`
	Folder string `gorm:"column:folder" json:"folder" yaml:"-"`
	FileName  string `gorm:"column:fileName" json:"fileName" yaml:"-"`
	ReferName  string `gorm:"column:referName" json:"referName" yaml:"-"`

	Fields []ZdField `gorm:"-" json:"fields"`
}
func (*ZdDef) TableName() string {
	return constant.TablePrefix + "def"
}

type ZdField struct {
	Model
	Field string `gorm:"column:field" json:"field"`
	Note string `gorm:"column:note" json:"note"`

	Range string `gorm:"column:range" json:"range"`
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
	Config	string `gorm:"column:config" json:"config"`
	Use	string `gorm:"column:use" json:"use"`
	From	string `gorm:"column:fromCol" json:"fromCol"`
	Select	string `gorm:"column:selectCol" json:"selectCol"`
	Where	string `gorm:"column:whereCol" json:"whereCol"`
	Limit	int `gorm:"column:limitCol" json:"limitCol"`

	Exp  string `gorm:"column:exp" json:"exp"`
	DefID uint `gorm:"column:defID" json:"defID"`
	ParentID uint `gorm:"column:parentID" json:"parentID"`
	UseID	uint `gorm:"column:useID" json:"useID"`
	ConfigID	uint `gorm:"column:configID" json:"configID"`

	Ord int `gorm:"column:ord;default:1" json:"ord"`
	Fields []*ZdField `gorm:"-" json:"fields"`
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

type ZdSection struct {
	Model
	OwnerType string   `gorm:"column:ownerType" json:"ownerType"` // field or instances
	OwnerID uint   `gorm:"column:ownerID" json:"ownerID"`
	Type    string `gorm:"column:type;default:scope" json:"type"`
	Value   string `gorm:"column:value" json:"value"`
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

type ZdRefer struct {
	Model
	OwnerType string   `gorm:"column:ownerType" json:"ownerType"` // field or instances
	OwnerID   uint   `gorm:"column:ownerID" json:"ownerID"`
	Type      string `gorm:"column:type" json:"type"`
	File      string `gorm:"column:file" json:"file"`
	ColName   string `gorm:"column:colName" json:"colName"`
	ColIndex  int    `gorm:"column:colIndex" json:"colIndex"`
	Count     int    `gorm:"column:count" json:"count"`
	HasTitle  bool   `gorm:"column:hasTitle" json:"hasTitle"`
}
func (*ZdRefer) TableName() string {
	return constant.TablePrefix + "refer"
}

type ZdRanges struct {
	Model
	Title    string `gorm:"column:title" json:"title"`
	Desc     string `gorm:"column:desc" json:"desc"`
	Prefix   string `gorm:"column:prefix" json:"prefix"`
	Postfix  string `gorm:"column:postfix" json:"postfix"`
	Format   string `gorm:"column:format" json:"format"`

	Yaml   string `gorm:"yaml" json:"yaml"`
	Path   string `gorm:"column:path" json:"path" yaml:"-"`
	Folder string `gorm:"folder" json:"folder" yaml:"-"`
	FileName  string `gorm:"column:fileName" json:"fileName" yaml:"-"`
	ReferName  string `gorm:"column:referName" json:"referName" yaml:"-"`

	Ranges []ZdRangesItem `gorm:"ForeignKey:rangesID" json:"ranges" yaml:"-"`
	RangeMap map[string]string  `gorm:"-" yaml:"ranges"`
}
func (*ZdRanges) TableName() string {
	return constant.TablePrefix + "ranges"
}

type ZdRangesItem struct {
	Model
	RangesID uint `gorm:"column:rangesID" json:"rangesID"`
	Field string `gorm:"column:name" json:"field"`
	Value string `gorm:"column:value" json:"value"`
	Ord int `gorm:"column:ord" json:"ord"`

	// for tree node
	ParentID uint `gorm:"-" json:"parentID"`
	Fields []*ZdRangesItem `gorm:"-" json:"fields"`
}
func (*ZdRangesItem) TableName() string {
	return constant.TablePrefix + "ranges_item"
}

type ZdInstances struct {
	Model     `yaml:"-"`
	Title     string            `gorm:"column:title" json:"title" yaml:"title,omitempty"`
	Desc      string            `gorm:"column:desc" json:"desc" yaml:"desc,omitempty"`

	Yaml   string `gorm:"yaml" json:"yaml" yaml:"-"`
	Path   string `gorm:"column:path" json:"path" yaml:"-"`
	Folder string `gorm:"folder" json:"folder" yaml:"-"`

	FileName  string `gorm:"column:fileName" json:"fileName" yaml:"-"`
	ReferName  string `gorm:"column:referName" json:"referName" yaml:"-"`

	Instances []ZdInstancesItem `gorm:"ForeignKey:instancesID" json:"instances" yaml:"instances"`
}
func (*ZdInstances) TableName() string {
	return constant.TablePrefix + "instances"
}

type ZdInstancesItem struct {
	Model            `yaml:"-"`

	Instance string `gorm:"column:instance" json:"instance" yaml:"instance,omitempty"`
	Note string `gorm:"column:note" json:"note" yaml:"note,omitempty"`

	Field string `gorm:"column:field" json:"field" yaml:"field,omitempty"`
	Range string `gorm:"column:range" json:"range" yaml:"range,omitempty"`
	Prefix string `gorm:"column:prefix" json:"prefix" yaml:"prefix,omitempty"`
	Postfix string `gorm:"column:postfix" json:"postfix" yaml:"postfix,omitempty"`
	Loop string `gorm:"column:loop" json:"loop" yaml:"loop,omitempty"`
	Loopfix string `gorm:"column:loopfix" json:"loopfix" yaml:"loopfix,omitempty"`
	Format string `gorm:"column:format" json:"format" yaml:"format,omitempty"`
	Type string `gorm:"column:type" json:"type" yaml:"type,omitempty"`
	Mode string `gorm:"column:mode" json:"mode" yaml:"mode,omitempty"`
	Length int `gorm:"column:length" json:"length" yaml:"length,omitempty"`
	LeftPad string `gorm:"column:leftPad" json:"leftPad" yaml:"leftPad,omitempty"`
	RightPad string `gorm:"column:rightPad" json:"rightPad" yaml:"rightPad,omitempty"`
	Rand bool `gorm:"column:rand" json:"rand" yaml:"rand,omitempty"`

	Config	string `gorm:"column:config" json:"config" yaml:"config,omitempty"`
	Use	string `gorm:"column:use" json:"use" yaml:"use,omitempty"`
	From	string `gorm:"column:fromCol" json:"fromCol" yaml:"from,omitempty"`
	Select	string `gorm:"column:selectCol" json:"selectCol" yaml:"select,omitempty"`
	Where	string `gorm:"column:whereCol" json:"whereCol" yaml:"where,omitempty"`
	Limit	int `gorm:"column:limitCol" json:"limitCol" yaml:"limit,omitempty"`

	Exp  string `gorm:"column:exp" json:"exp" yaml:"exp,omitempty"`
	InstancesID      uint `gorm:"column:instancesID" json:"instancesID" yaml:"-"`
	ParentID         uint `gorm:"column:parentID" json:"parentID" yaml:"-"`
	ConfigID         uint `gorm:"column:configID" json:"configID" yaml:"-"`
	UseID            uint `gorm:"column:useID" json:"useID" yaml:"-"`

	Ord    int                `gorm:"column:ord;default:1" json:"ord" yaml:"-"`
	Fields []*ZdInstancesItem `gorm:"-" json:"fields" yaml:"fields,omitempty"`
	Froms  []*ZdInstancesItem `gorm:"-" json:"froms" yaml:"froms,omitempty"`

	// for range edit
	IsRange bool         `gorm:"column:isRange;default:true" json:"isRange" yaml:"-"`
	Sections []ZdSection `gorm:"ForeignKey:fieldID" json:"sections" yaml:"-"`

	// for refer edit
	Refer ZdRefer `gorm:"ForeignKey:fieldID" json:"refer" yaml:"-"`
}
func (*ZdInstancesItem) TableName() string {
	return constant.TablePrefix + "instances_item"
}

type ZdConfig struct {
	Model
	Title    string `gorm:"column:title" json:"title"`
	Desc     string `gorm:"column:desc" json:"desc"`

	Prefix string `gorm:"column:prefix" json:"prefix"`
	Postfix string `gorm:"column:postfix" json:"postfix"`
	Loop string `gorm:"column:loop" json:"loop"`
	Loopfix string `gorm:"column:loopfix" json:"loopfix"`
	Format string `gorm:"column:format" json:"format"`

	Yaml   string `gorm:"yaml" json:"yaml"`
	Path   string `gorm:"column:path" json:"path" yaml:"-"`
	Folder string `gorm:"folder" json:"folder" yaml:"-"`
	FileName  string `gorm:"column:fileName" json:"fileName" yaml:"-"`
	ReferName  string `gorm:"column:referName" json:"referName" yaml:"-"`
}
func (*ZdConfig) TableName() string {
	return constant.TablePrefix + "config"
}

type ZdText struct {
	Model
	Title     string `gorm:"column:title" json:"title"`

	Content   string `gorm:"column:content" json:"content"`
	Path   string `gorm:"column:path" json:"path" yaml:"-"`
	Folder string `gorm:"folder" json:"folder" yaml:"-"`
	FileName  string `gorm:"column:fileName" json:"fileName" yaml:"-"`
	ReferName  string `gorm:"column:referName" json:"referName" yaml:"-"`
}
func (*ZdText) TableName() string {
	return constant.TablePrefix + "text"
}

type ZdExcel struct {
	Model
	Title     string `gorm:"column:title" json:"title"`
	Sheet     string `gorm:"column:sheet" json:"sheet"`

	Yaml   string `gorm:"yaml" json:"yaml"`
	Path   string `gorm:"column:path" json:"path" yaml:"-"`
	Folder string `gorm:"folder" json:"folder" yaml:"-"`
	FileName  string `gorm:"column:fileName" json:"fileName" yaml:"-"`
	ReferName  string `gorm:"column:referName" json:"referName" yaml:"-"`
}
func (*ZdExcel) TableName() string {
	return constant.TablePrefix + "excel"
}
