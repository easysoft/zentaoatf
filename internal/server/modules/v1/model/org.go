package model

import ()

type Org struct {
	BaseModel

	Name string `json:"name"`
	Desc string `json:"desc" gorm:"column:descr"`
}

func (Org) TableName() string {
	return "biz_org"
}
