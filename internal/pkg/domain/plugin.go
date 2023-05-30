package commDomain

type PluginInstallReq struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type PluginStartReq struct {
}

type ZapScanReq struct {
	Session string `json:"session"`
}
