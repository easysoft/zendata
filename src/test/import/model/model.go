package model

import "time"

type DataCategory1 struct {
	BaseModel

	Name string `json:"name"`
	Desc string `json:"desc" gorm:"column:descr"`
}

func (DataCategory1) TableName() string {
	return "biz_data_category1"
}

type DataCategory2 struct {
	BaseModel

	Name      string `json:"name"`
	Desc      string `json:"desc" gorm:"column:descr"`
	DataTable string `json:"desc"`

	ParentId uint `json:"parentId"`
}

func (DataCategory2) TableName() string {
	return "biz_data_category2"
}

type BaseModel struct {
	ID        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`
}

type DataCountry struct {
	BaseModel

	ContinentId  int    `json:"continentId"`
	Continent    string `json:"continent"`
	AreaCode     string `json:"areaCode"`
	EnglishShort string `json:"englishShort"`
	EnglishFull  string `json:"englishFull"`
	ChineseShort string `json:"chineseShort"`
	ChineseFull  string `json:"chineseFull"`
}

func (DataCountry) TableName() string {
	return "biz_data_country"
}

type DataCity struct {
	BaseModel

	Name    string `json:"name"`
	Code    string `json:"code"`
	ZipCode string `json:"zipCode"`
	State   string `json:"state"`
}

func (DataCity) TableName() string {
	return "biz_data_city"
}

type DataColor struct {
	BaseModel

	English string `json:"english"`
	Chinese string `json:"chinese"`
	Hex     string `json:"hex"`
	Rgb     string `json:"rgb"`
}

func (DataColor) TableName() string {
	return "biz_data_color"
}

type DataChineseFamily struct {
	BaseModel

	Name   string `json:"name"`
	Pinyin string `json:"pinyin"`
	Double bool   `json:"double"`
}

func (DataChineseFamily) TableName() string {
	return "biz_data_chinese_family"
}

type DataChineseGiven struct {
	BaseModel

	Name   string `json:"name"`
	Pinyin string `json:"pinyin"`
	Sex    string `json:"sex"`
}

func (DataChineseGiven) TableName() string {
	return "biz_data_chinese_given"
}

type DataEnglishFamily struct {
	BaseModel

	Name  string `json:"name"`
	Index string `json:"index"`
}

func (DataEnglishFamily) TableName() string {
	return "biz_data_english_family"
}

type DataEnglishGiven struct {
	BaseModel

	Name  string `json:"name"`
	Index string `json:"index"`
	Sex   string `json:"sex"`
}

func (DataEnglishGiven) TableName() string {
	return "biz_data_english_given"
}

type DataWordsInternetArgot struct {
	BaseModel
}
type DataWordsPreposition struct {
	BaseModel
}
type DataWordsPronoun struct {
	BaseModel
}
type DataWordsAdverb struct {
	BaseModel
}
type DataWordsVerb struct {
	BaseModel
}
type DataWordsAuxiliary struct {
	BaseModel
}
type DataWordsNoun struct {
	BaseModel
}
type DataWordsAdjectivePredicate struct {
	BaseModel
}
type DataWordsAdjective struct {
	BaseModel
}
type DataWordsNumeral struct {
	BaseModel
}
type DataWordsConjunction struct {
	BaseModel
}
type DataWordsMeasure struct {
	BaseModel
}
