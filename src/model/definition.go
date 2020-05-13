package model

type Definition struct {
	Title string `yaml:"title"`
	Desc string `yaml:"desc"`
	Author string `yaml:"author"`
	Version string`yaml:"version"`

	Fields  []Field `yaml:"fields,flow"`
}

type Field struct {
	Name     string  `yaml:"name"`
	Note     string  `yaml:"note"`
	Type string  `yaml:"type"`
	Range    string  `yaml:"range"`
	Filter    string  `yaml:"filter"`
	Prefix   string  `yaml:"prefix"`
	Postfix  string  `yaml:"postfix"`
	Loop  int  `yaml:"loop"`
	Loopfix  string  `yaml:"loopfix"`
	Format  string  `yaml:"format"`
	IsNumb  bool  `yaml:"isNumb"`
	Expect  string  `yaml:"expect"`
	Fields   []Field `yaml:"fields,flow"`

	Precision int
}

type FieldValue struct {
	Name     string
	Type     string
	Precision int
	Level int

	Values   []interface{}
	Children []FieldValue
}