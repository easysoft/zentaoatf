package model

type TestExec struct {
	BaseModel

	Name  string `json:"name"`
	Cases string `json:"cases"`
}

func (TestExec) TableName() string {
	return "biz_test_exec"
}
