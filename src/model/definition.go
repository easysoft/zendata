package model

type ClsBase struct {
	Title   string `yaml:"title"`
	Desc    string `yaml:"desc"`
	Author  string `yaml:"author"`
	Version string `yaml:"version"`
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
}

// range res
type ResRanges struct {
	ClsBase   `yaml:",inline"`
	FieldBase   `yaml:",inline"`
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

	Path string
}

// base struct
type FieldBase struct {
	FieldSimple   `yaml:",inline"`

	Config	string  `yaml:"config"`
	From	string  `yaml:"from"`
	Select	string  `yaml:"select"`
	Where	string  `yaml:"where"`
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
	Loop  int  `yaml:"loop"`
	Loopfix  string  `yaml:"loopfix"`
	Format  string  `yaml:"format"`
}


type FieldValue struct {
	FieldBase   `yaml:",inline"`
	Field     string  `yaml:"field"`
	Values   []interface{}
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