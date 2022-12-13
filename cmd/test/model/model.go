package model

import "time"

type TableInfo struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default string
	Extra   string
}

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
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`

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
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`

	Name        string `json:"name"`
	Code        string `json:"code"`
	ZipCode     string `json:"zipCode"`
	State       string `json:"state"`
	StateShort  string `json:"stateShort"`
	StateShort2 string `json:"stateShort2"`
}

func (DataCity) TableName() string {
	return "biz_data_city"
}

type DataIdiom struct {
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`

	Derivation   string `json:"derivation"`
	Example      string `json:"derivation"`
	Explanation  string `json:"derivation"`
	Pinyin       string `json:"derivation"`
	Word         string `json:"derivation"`
	Abbreviation string `json:"derivation"`
}

func (DataIdiom) TableName() string {
	return "biz_data_idiom"
}

type DataXiehouyu struct {
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`

	Riddle string `json:"riddle"`
	Answer string `json:"answer"`
	Tag    string `json:"tag"`
}

func (DataXiehouyu) TableName() string {
	return "biz_data_xiehouyu"
}

type DataDict struct {
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`

	Word        string `json:"word"`
	OldWord     string `json:"oldWord"`
	Strokes     string `json:"strokes"`
	Pinyin      string `json:"pinyin"`
	Radicals    string `json:"radicals"`
	Explanation string `json:"explanation"`
}

func (DataDict) TableName() string {
	return "biz_data_dict"
}

type DataChronology struct {
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`

	Content string `json:"content"`
}

func (DataChronology) TableName() string {
	return "biz_data_chronology"
}

type DataCompany struct {
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`

	Content string `json:"content"`
}

func (DataCompany) TableName() string {
	return "biz_data_company"
}

type DataFiveElements struct {
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`

	Content string `json:"content"`
}

func (DataFiveElements) TableName() string {
	return "biz_data_five_elements"
}

type DataHeavenlyStems struct {
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`

	Content string `json:"content"`
}

func (DataHeavenlyStems) TableName() string {
	return "biz_data_heavenly_stems"
}

type DataOccupation struct {
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`

	Content string `json:"content"`
}

func (DataOccupation) TableName() string {
	return "biz_data_occupation"
}

type DataPlanet struct {
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`

	Content string `json:"content"`
}

func (DataPlanet) TableName() string {
	return "biz_data_planet"
}

type DataSeason struct {
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`

	Content string `json:"content"`
}

func (DataSeason) TableName() string {
	return "biz_data_season"
}

type DataEarthlyBranches struct {
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`

	Content string `json:"content"`
}

func (DataEarthlyBranches) TableName() string {
	return "biz_data_earthly_branches"
}

type DataCompanyAbbreviation struct {
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`

	Content string `json:"content"`
}

func (DataCompanyAbbreviation) TableName() string {
	return "biz_data_company_abbreviation"
}

type DataChineseFamily struct {
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`

	Name   string `json:"name"`
	Pinyin string `json:"pinyin"`
	Double bool   `json:"double"`
}

func (DataChineseFamily) TableName() string {
	return "biz_data_chinese_family"
}

type DataChineseGiven struct {
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`

	Name   string `json:"name"`
	Pinyin string `json:"pinyin"`
	Sex    string `json:"sex"`
}

func (DataChineseGiven) TableName() string {
	return "biz_data_chinese_given"
}

type DataEnglishFamily struct {
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`

	Name  string `json:"name"`
	Index string `json:"index"`
}

func (DataEnglishFamily) TableName() string {
	return "biz_data_english_family"
}

type DataEnglishGiven struct {
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`

	Name  string `json:"name"`
	Index string `json:"index"`
	Sex   string `json:"sex"`
}

func (DataEnglishGiven) TableName() string {
	return "biz_data_english_given"
}

type DataWordsInternetArgot struct {
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`
}
type DataWordsPreposition struct {
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`
}
type DataWordsPronoun struct {
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`
}
type DataWordsAdverb struct {
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`
}
type DataWordsVerb struct {
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`
}
type DataWordsAuxiliary struct {
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`
}
type DataWordsNoun struct {
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`
}
type DataWordsAdjectivePredicate struct {
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`
}
type DataWordsAdjective struct {
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`
}
type DataWordsNumeral struct {
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`
}
type DataWordsConjunction struct {
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`
}
type DataWordsMeasure struct {
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`
}

type DataWordTagGroup struct {
	Id        uint           `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time     `json:"createdAt,omitempty"`
	UpdatedAt *time.Time     `json:"updatedAt,omitempty"`
	Deleted   bool           `json:"-" gorm:"default:false"`
	Disabled  bool           `json:"disabled,omitempty" gorm:"default:false"`
	Name      string         `gorm:"uniqueIndex" json:"name"`
	Tags      []*DataWordTag `gorm:"many2many:biz_data_word_tag_group_biz_data_word_tag" json:"tags"`
}

func (DataWordTagGroup) TableName() string {
	return "biz_data_word_tag_group"
}

type DataWordTag struct {
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`

	Name  string      `gorm:"uniqueIndex" json:"name"`
	Words []*DataWord `gorm:"many2many:biz_data_word_biz_data_word_tag" json:"words"`

	Groups []*DataWordTagGroup `gorm:"many2many:biz_data_word_tag_group_biz_data_word_tag" json:"tags"`
}

func (DataWordTag) TableName() string {
	return "biz_data_word_tag"
}

type DataWord struct {
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`

	Word       string `json:"word"`
	TagGroupId uint   `json:"tagGroupId"`

	Tags []*DataWordTag `gorm:"many2many:biz_data_word_biz_data_word_tag" json:"tags"`
}

func (DataWord) TableName() string {
	return "biz_data_word"
}

type DataColor struct {
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`

	English string `json:"english"`
	Chinese string `json:"chinese"`
	Hex     string `json:"hex"`
	Rgb     string `json:"rgb"`
}

func (DataColor) TableName() string {
	return "biz_data_color"
}

type DataComm struct {
	Id        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`
}

type DataFood struct {
	DataComm
	Name string `json:"name"`
}

func (DataFood) TableName() string {
	return "biz_data_food"
}

type DataAnimalPlant struct {
	DataComm
	Name string `json:"name"`
}

func (DataAnimalPlant) TableName() string {
	return "biz_data_animal_plant"
}

type DataFruit struct {
	DataComm
	Name string `json:"name"`
}

func (DataFruit) TableName() string {
	return "biz_data_fruit"
}

type DataConstellation struct {
	DataComm
	Name string `json:"name"`
}

func (DataConstellation) TableName() string {
	return "biz_data_constellation"
}

type DataZodiac struct {
	DataComm
	Name string `json:"name"`
}

func (DataZodiac) TableName() string {
	return "biz_data_zodiac"
}

type EightDiagram struct {
	DataComm
	Name string `json:"name"`
}

func (EightDiagram) TableName() string {
	return "biz_data_eight_diagram"
}

type Dynasty struct {
	DataComm
	Name string `json:"name"`
}

func (Dynasty) TableName() string {
	return "biz_data_dynasty"
}

type CarBrand struct {
	DataComm
	Name string `json:"name"`
}

func (CarBrand) TableName() string {
	return "biz_data_car_brand"
}

type CarComponent struct {
	DataComm
	Name string `json:"name"`
}

func (CarComponent) TableName() string {
	return "biz_data_car_component"
}

type PcOs struct {
	DataComm
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
	Version   string `json:"version"`
	Website   string `json:"website"`
}

func (PcOs) TableName() string {
	return "biz_data_pc_os"
}

type PcFileExt struct {
	DataComm
	Name string `json:"name"`
	Desc string `json:"desc"`
}

func (PcFileExt) TableName() string {
	return "biz_data_pc_file_ext"
}

type PhoneModel struct {
	DataComm

	Brand     string `json:"brand"`
	BrandName string `json:"brandName"`

	Model     string `json:"model"`
	ModelName string `json:"modelName"`

	Area string `json:"area"`
}

func (PhoneModel) TableName() string {
	return "biz_data_phone_model"
}

type Place struct {
	DataComm
	Name string `json:"name"`
}

func (Place) TableName() string {
	return "biz_data_place"
}

type SongData struct {
	DataComm
	Name   string `json:"name"`
	Singer string
	Lyric  []string
}

type Song struct {
	DataComm
	Name   string `json:"name"`
	Singer string
	Lyric  string
}

func (Song) TableName() string {
	return "biz_data_song"
}

type Bank struct {
	DataComm
	Name string `json:"name"`
}

func (Bank) TableName() string {
	return "biz_data_bank"
}

type Advert struct {
	DataComm
	Name string `json:"name"`
}

func (Advert) TableName() string {
	return "biz_data_advert"
}

type IdiomSimple struct {
	DataComm
	Name string `json:"name"`
}

func (IdiomSimple) TableName() string {
	return "biz_data_idiom_simple"
}

type PoetryAncient struct {
	DataComm
	Content    string `json:"content"`
	CategoryId uint   `json:"categoryId"`
}

func (PoetryAncient) TableName() string {
	return "biz_data_poetry_ancient"
}

type PoetryCategory struct {
	DataComm
	Name   string `json:"name" gorm:"unique_index, not null"`
	EnName string `json:"cnName" gorm:"unique_index, not null"`
}

func (PoetryCategory) TableName() string {
	return "biz_data_poetry_category"
}

type Medicine struct {
	DataComm
	Name string `json:"name"`
}

func (Medicine) TableName() string {
	return "biz_data_medicine"
}

type Joke struct {
	DataComm
	Content string `json:"content"`
}

func (Joke) TableName() string {
	return "biz_data_joke"
}
