package model

type Server struct {
	BaseModel

	Path    string `json:"path"`
	Default bool   `json:"default"`
}

func (Server) TableName() string {
	return "biz_server"
}
