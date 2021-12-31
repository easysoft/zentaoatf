package model

type TestSuite struct {
	BaseModel

	Version int    `json:"version"`
	Name    string `json:"name"`
	Code    string `json:"code"`
	Desc    string `json:"desc"`

	ProjectId uint `json:"projectId"`
}

func (TestSuite) TableName() string {
	return "biz_test_suite"
}
