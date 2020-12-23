package domain

type SysInfo struct {
	AgentDir string `json:"agentDir,omitempty"`

	SysArch  string `json:"sysArch,omitempty"`
	SysCores int    `json:"sysCores,omitempty"`

	OsType    string `json:"osType,omitempty"`
	OsName    string `json:"osName,omitempty"`
	OsVersion string `json:"osVersion,omitempty"`

	Local string `json:"local,omitempty"`
	Lang  string `json:"lang,omitempty"`

	IP  string `json:"ip,omitempty"`
	MAC string `json:"mac,omitempty"`
}
