package serverDomain

type TestAsset struct {
	Title string `json:"title"`
	Key   string `json:"key"`

	Children []TestAsset `json:"children"`
}
