package model

type TestCase struct {
	BaseModel

	Version int `json:"version"`
	Name    string  `json:"name"`
	Desc    string  `json:"desc"`
}

func (TestCase) TableName() string {
	return "biz_test_case"
}
