package consts

type ResultCode int

const (
	ResultCodeSuccess ResultCode = 0
	ResultCodeFail    ResultCode = 1
)

func (e ResultCode) Int() int {
	return int(e)
}
