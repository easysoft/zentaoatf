package model

type Server struct {
	BaseModel

	Name      string `json:"name"`
	Path      string `json:"path"`
	IsDefault bool   `json:"is_default"`
}

func (Server) TableName() string {
	return "biz_server"
}
