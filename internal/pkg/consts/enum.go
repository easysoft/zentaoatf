package consts

type ResultCode int

const (
	ResultCodeSuccess ResultCode = 0
	ResultCodeFail    ResultCode = 1
)

func (e ResultCode) Int() int {
	return int(e)
}

type ResultStatus string

const (
	PASS    ResultStatus = "pass"
	FAIL    ResultStatus = "fail"
	SKIP    ResultStatus = "skip"
	BLOCKED ResultStatus = "blocked"
	UNKNOWN ResultStatus = "unknown"
)

func (e ResultStatus) String() string {
	return string(e)
}
