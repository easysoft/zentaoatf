package misc

type LangType int

const (
	PHP LangType = iota
	PYTHON
	GO
	UNKNOWN
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
