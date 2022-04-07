package serverDomain

type TestScript struct {
	Id      int    `json:"id"`
	Version int    `json:"version"`
	Name    string `json:"name"`
	Code    string `json:"code"`
	Desc    string `json:"desc"`
	Lang    string `json:"lang"`

	Path        string `json:"path"`
	WorkspaceId int    `json:"workspaceId"`
}

type FilterItem struct {
	Label string `json:"label"`
	Value int    `json:"value"`
}
