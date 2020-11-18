package model

type ClsBase struct {
	Title   string `yaml:"title"`
	Desc    string `yaml:"desc"`
	Author  string `yaml:"author,omitempty"`
	Version string `yaml:"version,omitempty"`

	From string        `yaml:"from,omitempty"`
	Type  string  `yaml:"type,omitempty"`
}

// instance res
type ResInsts struct {
	ClsBase   `yaml:",inline"`
	Field string        `yaml:"field"`
	Instances []ResInst `yaml:"instances,flow"`

}
type ResInst struct {
	FieldBase   `yaml:",inline"`
	Instance string  `yaml:"instance"`
	Fields  []DefField `yaml:"fields,flow"`

	Froms []DefField `yaml:"froms,flow"`
}

// range res
type ResRanges struct {
	ClsBase   `yaml:",inline"`
	Field string        `yaml:"field"`
	Ranges map[string]string  `yaml:"ranges"`
}

// common item
type DefData struct {
	ClsBase   `yaml:",inline"`
	Fields  []DefField `yaml:"fields,flow"`
}
type DefField struct {
	FieldBase `yaml:",inline"`
	Fields    []DefField `yaml:"fields,flow,omitempty"`
	Length    int        `yaml:"length,omitempty"`
	LeftPad   string     `yaml:"leftpad,omitempty"`
	RightPad  string     `yaml:"rightpad,omitempty"`
	Path      string     `yaml:"path,omitempty"`

	Froms []DefField `yaml:"froms,flow,omitempty"`
}

// base struct
type FieldBase struct {
	FieldSimple   `yaml:",inline"`

	Config	string  `yaml:"config,omitempty"`
	From	string  `yaml:"from,omitempty"`
	Select	string  `yaml:"select,omitempty"`
	Where	string  `yaml:"where,omitempty"`
	Limit	int  `yaml:"limit,omitempty"`
	Use	string  `yaml:"use,omitempty"`

	IsNumb  bool  `yaml:"isNumb,omitempty"`
	Expect  string  `yaml:"expect,omitempty"`

	Precision int  `yaml:"precision,omitempty"`
}
type DefSimple struct {
	ClsBase   `yaml:",inline"`
	Fields  []FieldSimple `yaml:"fields"`
}
type FieldSimple struct {
	Field  string  `yaml:"field"`
	Note     string  `yaml:"note,omitempty"`
	Range    string  `yaml:"range"`
	Value string  `yaml:"value,omitempty"`
	Prefix   string  `yaml:"prefix,omitempty"`
	Postfix  string  `yaml:"postfix,omitempty"`
	Loop  string  `yaml:"loop,omitempty"`
	Loopfix  string  `yaml:"loopfix,omitempty"`
	Format  string  `yaml:"format,omitempty"`
	Rand  bool  `yaml:"rand,omitempty"`
	Type  string  `yaml:"type,omitempty"`
	Mode  string  `yaml:"mode,omitempty"`

	LoopStart          int  `yaml:"-"`
	LoopEnd            int  `yaml:"-"`
	LoopIndex          int  `yaml:"-"`
	IsRand             bool `yaml:"-"`
	ReferToAnotherYaml bool `yaml:"-"`
}

type FieldWithValues struct {
	FieldBase   `yaml:",inline"`
	Field     string  `yaml:"field"`
	Values   []interface{}
	ValuesWithPlaceholder []string
}

type DefInfo struct {
	Title string   `yaml:"title"`
	Desc string  `yaml:"desc"`

	Fields interface{}  `yaml:"fields,omitempty"` // is yaml
	Range string  `yaml:"range,omitempty"` // is config
	Ranges interface{}  `yaml:"ranges,omitempty"` // is ranges
	Instances interface{}  `yaml:"instances,omitempty"` // is instances
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
	ClsBase   `yaml:",inline"`
	XFields  []DefFieldExport `yaml:"xfields,flow"` // control orders
}
type DefFieldExport struct {
	Field string  `yaml:"field"`
	Prefix string  `yaml:"prefix,omitempty"`
	Postfix  string  `yaml:"postfix,omitempty"`

	Select	string  `yaml:"select,omitempty"`
	Where	string  `yaml:"where,omitempty"`
	Rand  bool  `yaml:"rand"`
	Limit	int  `yaml:"limit,omitempty"`
}
type Article struct {
	Title   string `yaml:"title"`
	Desc    string `yaml:"desc"`
	Author  string `yaml:"author"`
	Type  string  `yaml:"type"`
	XFields  []ArticleField `yaml:"xfields,flow"` // control orders
}
type ArticleField struct {
	Field string  `yaml:"field"`
	Range  string  `yaml:"range,omitempty"`
	Prefix string  `yaml:"prefix,omitempty"`
	Postfix  string  `yaml:"postfix,omitempty"`
}
type ArticleSent struct {
	Type string
	Val string
	IsParag bool
	IsSent bool
}

type ResFile struct {
	Path string `json:"path"`
	Name    string `json:"name"`
	Title string `json:"title"`
	Desc   string `json:"desc"`
	ResType string `json:"resType"`
}
type ResField struct {
	Index int `json:"index"`
	Name    string `json:"name"`
}
