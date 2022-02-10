package controller

import (
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
)

type TestSuiteCtrl struct {
	TestSuiteService *service.TestSuiteService `inject:""`
	BaseCtrl
}

func NewTestSuiteCtrl() *TestSuiteCtrl {
	return &TestSuiteCtrl{}
}
