package controller

import (
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
)

type TestSetCtrl struct {
	TestSetService *service.TestSetService `inject:""`
	BaseCtrl
}

func NewTestSetCtrl() *TestSetCtrl {
	return &TestSetCtrl{}
}
