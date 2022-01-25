package consts

type ResultCode int

const (
	ResultCodeSuccess ResultCode = 0
	ResultCodeFail    ResultCode = 1
)

func (e ResultCode) Int() int {
	return int(e)
}

type OsType string

const (
	OsWindows OsType = "windows"
	OsLinux   OsType = "linux"
	OsMac     OsType = "mac"
)

func (e OsType) String() string {
	return string(e)
}
func (OsType) Get(osName string) OsType {
	return OsType(osName)
}

type TestType string

const (
	TestFunc TestType = "func"
	TestUnit TestType = "unit"
	TestAuto TestType = "auto"
)

func (e TestType) String() string {
	return string(e)
}
func (TestType) Get(testType string) TestType {
	return TestType(testType)
}
