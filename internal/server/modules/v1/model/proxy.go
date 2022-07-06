package model

type Proxy struct {
	BaseModel

	Name string `json:"name"`
	Path string `json:"path"`
}

func (Proxy) TableName() string {
	return "biz_proxy"
}
