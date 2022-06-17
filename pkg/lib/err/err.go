package errUtils

func New(text string) error {
	return &errorStr{text}
}

type errorStr struct {
	str string
}

func (e *errorStr) Error() string {
	return e.str
}
