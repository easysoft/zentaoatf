package model

type TestSet struct {
	BaseModel

	Version int    `json:"version"`
	Name    string `json:"name"`
	Code    string `json:"code"`
	Desc    string `json:"desc"`

	ProjectId uint `json:"projectId"`
}

func (TestSet) TableName() string {
	return "biz_test_set"
}
