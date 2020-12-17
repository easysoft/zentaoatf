package serverModel

type SysInfo struct {
	SysArch  string `json:"sysArch"`
	SysCores int    `json:"sysCores"`

	OsType    string      `json:"osType"`
	OsName    interface{} `json:"osName"`
	OsVersion interface{} `json:"osVersion"`

	Local string `json:"local"`
	Lang  string `json:"lang"`
}
