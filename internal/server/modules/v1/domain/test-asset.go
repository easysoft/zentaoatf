package serverDomain

type TestAsset struct {
	Title string `json:"title"`
	Key   string `json:"key"`
	Path  string `json:"path"`

	IsDir       bool `json:"isDir"`
	ScriptCount int  `json:"scriptCount"`

	Children []*TestAsset `json:"children"`
}
