package model

type Proxy struct {
	BaseModel

	Path string `json:"path"`
}

func (Proxy) TableName() string {
	return "biz_proxy"
}
