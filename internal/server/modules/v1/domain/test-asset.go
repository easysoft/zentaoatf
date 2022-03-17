package serverDomain

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/kataras/iris/v12"
)

type TestAsset struct {
	Title string `json:"title"`
	Key   string `json:"key"`
	Path  string `json:"path"`

	Type        commConsts.TreeNodeType `json:"type"`
	ScriptCount int                     `json:"scriptCount"`
	Slots       iris.Map                `json:"slots"`

	Children []*TestAsset `json:"children"`
}
