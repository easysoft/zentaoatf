package model

type Site struct {
	BaseModel

	Name     string `json:"name"`
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`

	IsDefault bool `json:"isDefault"`
}

func (Site) TableName() string {
	return "biz_site"
}
