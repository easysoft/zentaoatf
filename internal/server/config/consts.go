package serverConfig

const (
	HeartbeatInterval = 15
	JobCheckInterval  = 15

	WsDefaultNameSpace = "default"
	WsDefaultRoom      = "default"
	WsEvent            = "OnChat"

	ApiPath = "/api/v1"
	WsPath  = ApiPath + "/ws"

	ZentaoCasePrefix   = "zentao-case-"
	ZentaoModulePrefix = "zentao-module-"
)

var (
	SortMap = map[string]string{
		"ascend":  "ASC",
		"descend": "DESC",
		"":        "ASC",
	}
)
