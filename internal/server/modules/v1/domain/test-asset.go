package serverDomain

import (
	"github.com/kataras/iris/v12"
)

type TestAsset struct {
	Title string `json:"title"`
	Key   string `json:"key"`
	Path  string `json:"path"`

	IsDir       bool     `json:"isDir"`
	ScriptCount int      `json:"scriptCount"`
	Slots       iris.Map `json:"slots"`

	Children []*TestAsset `json:"children"`
}
