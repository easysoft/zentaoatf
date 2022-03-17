package model

import commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"

type Workspace struct {
	BaseModel

	Path      string              `json:"path"`
	Name      string              `json:"name"`
	Desc      string              `json:"desc" gorm:"column:descr"`
	Type      commConsts.TestTool `json:"type" gorm:"default:ztf"`
	Cmd       string              `json:"cmd"`
	ProductId uint                `json:"productId"`

	IsDefault bool `json:"isDefault"`
}

func (Workspace) TableName() string {
	return "biz_workspace"
}
