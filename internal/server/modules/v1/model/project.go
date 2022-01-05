package model

type Project struct {
	BaseModel

	Path string `json:"path"`
	Name string `json:"name"`
	Desc string `json:"desc" gorm:"column:descr"`

	IsDefault bool `json:"isDefault"`
}

func (Project) TableName() string {
	return "biz_project"
}
