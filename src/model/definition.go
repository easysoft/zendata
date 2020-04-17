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
	Prefix   string  `yaml:"prefix"`
	Postfix  string  `yaml:"postfix"`
	Loop  int  `yaml:"loop"`
	Loopfix  string  `yaml:"loopfix"`
	Expect  string  `yaml:"expect"`
	Fields   []Field `yaml:"fields,flow"`

	Precision int
}
