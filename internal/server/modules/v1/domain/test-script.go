package serverDomain

type TestScript struct {
	Version int    `json:"version"`
	Name    string `json:"name"`
	Code    string `json:"code"`
	Desc    string `json:"desc"`
	Lang    string `json:"lang"`

	ProjectId uint `json:"projectId"`
}
