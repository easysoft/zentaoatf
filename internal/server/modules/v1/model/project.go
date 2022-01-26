package model

import commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"

type Project struct {
	BaseModel

	Type commConsts.TestType `json:"type"`
	Path string              `json:"path"`
	Name string              `json:"name"`
	Desc string              `json:"desc" gorm:"column:descr"`

	IsDefault bool `json:"isDefault"`
}

func (Project) TableName() string {
	return "biz_project"
}
