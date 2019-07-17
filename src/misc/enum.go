package misc

type LangType int

const (
	PHP LangType = iota
	PYTHON
	GO
)

func (c LangType) String() string {
	switch c {
	case PHP:
		return "php"
	case PYTHON:
		return "python"
	case GO:
		return "go"
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

const (
	SuiteExt string = "suite"
)
