package serverDomain

import commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"

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

type CreateScriptReq struct {
	Mode        commConsts.NodeCreateMode `json:"mode"`
	Type        commConsts.NodeCreateType `json:"type"`
	Target      string                    `json:"target"`
	Name        string                    `json:"name"`
	WorkspaceId int                       `json:"workspaceId"`
	ProductId   int                       `json:"productId"`
}

type MoveScriptReq struct {
	DragKey      string             `json:"dragKey"`
	DropKey      string             `json:"dropKey"`
	DropPosition commConsts.DropPos `json:"dropPosition"`
	WorkspaceId  int                `json:"workspaceId"`
}
