package model

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"time"
)

type ClsInfo struct {
	Title   string `yaml:"title"`
	Desc    string `yaml:"desc"`
	Author  string `yaml:"author,omitempty"`
	Version string `yaml:"version,omitempty"`
}
type ClsBase struct {
	ClsInfo `yaml:",inline"`
	From    string `yaml:"from,omitempty"`
	Type    string `yaml:"type,omitempty"`
}

// config res
type ResConfig struct {
	ClsInfo     `yaml:",inline"`
	FieldSimple `yaml:",inline"`
}

// range res
type ResRanges struct {
	ClsInfo `yaml:",inline"`
	FileDir string            `yaml:"fileDir,omitempty"`
	Ranges  map[string]string `yaml:"ranges"`
}

// instance res
type ResInstances struct {
	ClsBase   `yaml:",inline"`
	FileDir   string             `yaml:"fileDir,omitempty"`
	Instances []ResInstancesItem `yaml:"instances,flow"`
}
type ResInstancesItem struct {
	FieldBase `yaml:",inline"`
	Instance  string     `yaml:"instance"`
	Fields    []DefField `yaml:"fields,flow"`
	Froms     []DefField `yaml:"froms,flow"`
}

// common item
type DefData struct {
	ClsBase `yaml:",inline"`
	Fields  []DefField `yaml:"fields"`

	Content string `yaml:"content,omitempty"` // for article only
}
type DefField struct {
	FieldBase `yaml:",inline"`
	Fields    []DefField `yaml:"fields,omitempty"`
	Length    int        `yaml:"length,omitempty"`
	LeftPad   string     `yaml:"leftpad,omitempty"`
	RightPad  string     `yaml:"rightpad,omitempty"`
	Path      string     `yaml:"path,omitempty"`
	FileDir   string     `yaml:"fileDir,omitempty"`
	Join      bool       `yaml:"join,omitempty"`

	Froms []DefField `yaml:"froms,flow,omitempty"`

	Values                []interface{} `yaml:"-"`
	ValuesWithPlaceholder []string      `yaml:"-"`
}

type DefArticle struct {
	Author  string `yaml:"author"`
	From    string `yaml:"from"`
	Title   string `yaml:"title"`
	Type    string `yaml:"type"`
	Version string `yaml:"version"`
	Content string `yaml:"content"`
}

// base struct
type FieldBase struct {
	FieldSimple `yaml:",inline"`

	Config string `yaml:"config,omitempty"`
	Select string `yaml:"select,omitempty"`
	Where  string `yaml:"where,omitempty"`
	Limit  int    `yaml:"limit,omitempty"`

	IsNumb bool   `yaml:"isNumb,omitempty"`
	Expect string `yaml:"expect,omitempty"`

	Precision int `yaml:"precision,omitempty"`
}
type DefSimple struct {
	ClsBase `yaml:",inline"`
	Fields  []FieldSimple `yaml:"fields"`
}
type FieldSimple struct {
	Field   string `yaml:"field,omitempty"`
	Note    string `yaml:"note,omitempty"`
	Range   string `yaml:"range,omitempty"`
	Value   string `yaml:"value,omitempty"`
	Prefix  string `yaml:"prefix,omitempty"`
	Postfix string `yaml:"postfix,omitempty"`
	Divider string `yaml:"divider,omitempty"`
	Loop    string `yaml:"loop,omitempty"`
	Loopfix string `yaml:"loopfix,omitempty"`
	Format  string `yaml:"format,omitempty"`
	Rand    bool   `yaml:"rand,omitempty"`
	Type    string `yaml:"type,omitempty"`
	Mode    string `yaml:"mode,omitempty"`
	Items   int    `yaml:"records,omitempty"`
	//ParentItems int    `yaml:"-"`
	//ParentJoin  bool   `yaml:"-"`
	From string `yaml:"from,omitempty"`
	Use  string `yaml:"use,omitempty"`

	LoopStart          int  `yaml:"-"`
	LoopEnd            int  `yaml:"-"`
	LoopIndex          int  `yaml:"-"`
	IsRand             bool `yaml:"-"`
	ReferToAnotherYaml bool `yaml:"-"`

	PrefixRange  *Range `yaml:"-"`
	PostfixRange *Range `yaml:"-"`
}

type Range struct {
	Values []interface{}
	IsRand bool
}

type FieldWithValues struct {
	FieldBase             `yaml:",inline"`
	Field                 string `yaml:"field"`
	Values                []interface{}
	ValuesWithPlaceholder []string
}

type DefInfo struct {
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`

	Fields    interface{} `yaml:"fields,omitempty"`    // is yaml
	Range     string      `yaml:"range,omitempty"`     // is config
	Ranges    interface{} `yaml:"ranges,omitempty"`    // is ranges
	Instances interface{} `yaml:"instances,omitempty"` // is instances
}

func (def *DefSimple) Init(tableName, author, desc, version string) {
	def.Title = "table " + tableName
	def.Author = author
	def.Desc = desc
	def.Version = version
}
func (fld *FieldSimple) Init(field string) {
	fld.Field = field
}

type DefExport struct {
	ClsBase `yaml:",inline"`
	XFields []DefFieldExport `yaml:"xfields,flow"` // control orders
}
type DefFieldExport struct {
	Field   string `yaml:"field"`
	Prefix  string `yaml:"prefix,omitempty"`
	Postfix string `yaml:"postfix,omitempty"`
	Divider string `yaml:"divider,omitempty"`
	Select  string `yaml:"select,omitempty"`
	Where   string `yaml:"where,omitempty"`
	Rand    bool   `yaml:"rand"`
	Limit   int    `yaml:"limit,omitempty"`
}
type Article struct {
	Title   string         `yaml:"title"`
	Desc    string         `yaml:"desc"`
	Author  string         `yaml:"author"`
	Type    string         `yaml:"type"`
	XFields []ArticleField `yaml:"xfields,flow"` // control orders
}
type ArticleField struct {
	Field   string `yaml:"field"`
	Range   string `yaml:"range,omitempty"`
	Prefix  string `yaml:"prefix,omitempty"`
	Postfix string `yaml:"postfix,omitempty"`
}
type ArticleSent struct {
	Type    string
	Val     string
	IsParag bool
	IsSent  bool
}

type ResFile struct {
	Path      string    `json:"path"`
	Title     string    `json:"title"`
	Desc      string    `json:"desc"`
	ResType   string    `json:"resType"`
	UpdatedAt time.Time `json:"updatedAt"`

	FileName  string `json:"fileName"`
	ReferName string `json:"referName"`
}
type ResField struct {
	ID    uint   `json:"id"`
	Index int    `json:"index"`
	Name  string `json:"name"`
}

type Dir struct {
	Name string `json:"name"`
	Path string `json:-`
}

type MockData struct {
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Author  string `json:"author"`
	Version string `json:"version"`

	Paths map[string]map[string]map[string]map[string]*EndPoint `json:"paths"`
}

type EndPoint struct {
	Method    HttpMethod               `json:"method"`
	Summary   string                   `json:"summary"`
	Config    string                   `json:"config"`
	MediaType string                   `json:"mediaType"`
	Type      consts.OpenApiSchemaType `json:"type"`
	Lines     int                      `json:"lines"`
	Fields    string                   `json:"fields"`
}

// for preview
type MockPaths struct {
	PathKey   string        `json:"pathKey"`
	PathValue []MockMethods `json:"pathValue"`
}
type MockMethods struct {
	MethodKey   HttpMethod  `json:"methodKey"`
	MethodValue []MockCodes `json:"methodValue"`
}
type MockCodes struct {
	CodeKey   string       `json:"codeKey"`
	CodeValue []MockMedias `json:"codeValue"`
}
type MockMedias struct {
	MediaKey   string `json:"mediaKey"`
	RespSample string `json:"respSample"`
}

type HttpMethod string

const (
	Get    HttpMethod = "get"
	Post   HttpMethod = "post"
	Put    HttpMethod = "put"
	Delete HttpMethod = "delete"

	Patch   HttpMethod = "path"
	Head    HttpMethod = "head"
	Connect HttpMethod = "connect"
	Options HttpMethod = "options"
	Trace   HttpMethod = "trace"
)

func (e HttpMethod) String() string {
	return string(e)
}
