package model

type Interpreter struct {
	BaseModel

	Lang string `json:"lang"`
	Path string `json:"path"`
}

func (Interpreter) TableName() string {
	return "biz_interpreter"
}
