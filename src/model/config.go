package model

type Config struct {
	Version  int
	Language string

	Url      string
	Account  string
	Password string

	Javascript string
	Lua        string
	Perl       string
	Php        string
	Python     string
	Ruby       string
	Tcl        string
	Autoit     string
}
