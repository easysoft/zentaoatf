package serverConfig

const (
	WebCheckInterval = 60 * 60

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
