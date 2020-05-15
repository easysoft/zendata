package model

type ClsBase struct {
	Title   string `yaml:"title"`
	Desc    string `yaml:"desc"`
	Author  string `yaml:"author"`
	Version string `yaml:"version"`
}

type ClsRange struct {
	ClsBase
	FieldBase
	Field string
	Ranges map[string]string  `yaml:"ranges"`
}
type ClsInst struct {
	ClsBase
	Field string
	Instances []struct {
		FieldBase
		Instance string  `yaml:"Instance"`
	} `yaml:"fields,flow"`
}
type DefData struct {
	ClsBase
	Fields  []DefField `yaml:"fields,flow"`
}
type DefField struct {
	FieldBase
	Field     string  `yaml:"field"`
	Range    string  `yaml:"range"`

	Fields   []DefField `yaml:"fields,flow"`
}

type FieldBase struct {
	Note     string  `yaml:"note"`

	From	string  `yaml:"from"`
	Select	string  `yaml:"select"`
	Where	string  `yaml:"where"`
	Use	string  `yaml:"use"`

	Prefix   string  `yaml:"prefix"`
	Postfix  string  `yaml:"postfix"`
	Loop  int  `yaml:"loop"`
	Loopfix  string  `yaml:"loopfix"`
	Format  string  `yaml:"format"`
	IsNumb  bool  `yaml:"isNumb"`
	Expect  string  `yaml:"expect"`

	Precision int
}

type FieldValue struct {
	FieldBase
	Field     string  `yaml:"field"`

	Values   []interface{}
	Children []FieldValue
}