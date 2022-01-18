package commConsts

type ZentaoRequestType string

const (
	PathInfo ZentaoRequestType = "PATH_INFO"
	Get      ZentaoRequestType = "GET"
	Empty    ZentaoRequestType = ""
)

func (e ZentaoRequestType) String() string {
	return string(e)
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

type ExecCmd string

const (
	ExecCase   ExecCmd = "execCase"
	ExecModule ExecCmd = "execModule"
	ExecSuite  ExecCmd = "execSuite"
	ExecTask   ExecCmd = "execTask"
)

func (e ExecCmd) String() string {
	return string(e)
}
