package serverDomain

import (
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	"github.com/kataras/iris/v12"
)

type TestAsset struct {
	Title string `json:"title"`
	Key   string `json:"key"`
	Path  string `json:"path"`

	Type        commConsts.TreeNodeType `json:"type"`
	ModuleId    int                     `json:"moduleId"`
	CaseId      int                     `json:"caseId"`
	Lang        string                  `json:"lang"`
	ScriptCount int                     `json:"scriptCount"`
	Slots       iris.Map                `json:"slots"`

	WorkspaceId   int                 `json:"workspaceId"`
	WorkspaceType commConsts.TestTool `json:"workspaceType"`
	Children      []*TestAsset        `json:"children"`

	Checkable bool `json:"checkable"`
	IsLeaf    bool `json:"isLeaf"`
}
