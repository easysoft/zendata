package model

type ClsBase struct {
	Title   string `yaml:"title"`
	Desc    string `yaml:"desc"`
	Author  string `yaml:"author"`
	Version string `yaml:"version"`

	From string        `yaml:"from"`
	Type  string  `yaml:"type"`
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
	FieldBase   `yaml:",inline"`
	Fields   []DefField `yaml:"fields,flow"`
	Width int `yaml:"width"`
	LeftPad string `yaml:"leftPad"`
	RightPad string `yaml:"rightPad"`
	Path string

	Froms []DefField `yaml:"froms,flow"`
}

// base struct
type FieldBase struct {
	FieldSimple   `yaml:",inline"`

	Config	string  `yaml:"config"`
	From	string  `yaml:"from"`
	Select	string  `yaml:"select"`
	Where	string  `yaml:"where"`
	Limit	int  `yaml:"limit"`
	Use	string  `yaml:"use"`

	IsNumb  bool  `yaml:"isNumb"`
	Expect  string  `yaml:"expect"`

	Precision int
}
type DefSimple struct {
	ClsBase   `yaml:",inline"`
	Fields  []FieldSimple `yaml:"fields"`
}
type FieldSimple struct {
	Field  string  `yaml:"field"`
	Note     string  `yaml:"note"`
	Range    string  `yaml:"range"`
	Prefix   string  `yaml:"prefix"`
	Postfix  string  `yaml:"postfix"`
	Loop  string  `yaml:"loop"`
	Loopfix  string  `yaml:"loopfix"`
	Format  string  `yaml:"format"`
	Rand  bool  `yaml:"rand"`

	LoopStart  int  `yaml:"-"`
	LoopEnd  int  `yaml:"-"`
	LoopIndex int  `yaml:"-"`
	IsRand bool  `yaml:"-"`
	IsReferYaml bool  `yaml:"-"`
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