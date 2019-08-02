package misc

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
	}
	return "unknown"
}

type ResultStatus int

const (
	PASS ResultStatus = iota
	FAIL
	SKIP
)

func (c ResultStatus) String() string {
	switch c {
	case PASS:
		return "PASS"
	case FAIL:
		return "FAIL"
	case SKIP:
		return "SKIP"
	}
	return "UNKNOWN"
}

type RunMode int

const (
	DIR RunMode = iota
	BATCH
	SUITE
	SCRIPT
)

func (c RunMode) String() string {
	switch c {
	case DIR:
		return "dir"
	case BATCH:
		return "batch"
	case SUITE:
		return "suite"
	case SCRIPT:
		return "script"
	}
	return "unknown"
}
