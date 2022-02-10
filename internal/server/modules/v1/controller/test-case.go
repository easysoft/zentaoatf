package controller

import (
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
)

type TestCaseCtrl struct {
	TestCaseService *service.TestCaseService `inject:""`
	BaseCtrl
}

func NewTestCaseCtrl() *TestCaseCtrl {
	return &TestCaseCtrl{}
}
