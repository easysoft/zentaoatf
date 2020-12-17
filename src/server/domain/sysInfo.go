package domain

type SysInfo struct {
	AgentDir string `json:"agentDir"`

	SysArch  string `json:"sysArch"`
	SysCores int    `json:"sysCores"`

	OsType    string `json:"osType"`
	OsName    string `json:"osName"`
	OsVersion string `json:"osVersion"`

	Local string `json:"local"`
	Lang  string `json:"lang"`

	IP  string `json:"ip"`
	MAC string `json:"mac"`
}
