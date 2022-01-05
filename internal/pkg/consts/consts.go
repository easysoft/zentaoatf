package consts

import "os"

const (
	PthSep = string(os.PathSeparator)

	App       = "ztf"
	AppServer = "server"
	AppAgent  = "agent"

	RequestTypePathInfo = "PATH_INFO"
)
