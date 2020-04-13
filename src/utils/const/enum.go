package constant

type LangType int

const (
	GO LangType = iota
	LUA
	PERL
	PHP
	PYTHON
	RUBY
	SHELL
	TCL
	AUTOIT
)

func (c LangType) String() string {
	switch c {
	case GO:
		return "go"
	case LUA:
		return "lua"
	case PERL:
		return "perl"
	case PHP:
		return "php"
	case PYTHON:
		return "python"
	case RUBY:
		return "ruby"
	case SHELL:
		return "shell"
	case TCL:
		return "tcl"
	case AUTOIT:
	return "autoit"
	}
	return "unknown"
}

type ResultStatus int

const (
	PASS ResultStatus = iota
	FAIL
	SKIP
	BLOCKED
)

func (c ResultStatus) String() string {
	switch c {
	case PASS:
		return "pass"
	case FAIL:
		return "fail"
	case SKIP:
		return "skip"
	case BLOCKED:
		return "blocked"
	}

	return "UNKNOWN"
}

type RunMode int

const (
	RunModeDir RunMode = iota
	RunModeBatch
	RunModeSuite
	RunModeScript
)

func (c RunMode) String() string {
	switch c {
	case RunModeDir:
		return "dir"
	case RunModeBatch:
		return "batch"
	case RunModeSuite: // can be show with cui by select a suite file
		return "suite"
	case RunModeScript: // can be show with cui by select a script file
		return "script"
	}
	return "unknown"
}
