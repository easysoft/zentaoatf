package service

import (
	"github.com/easysoft/zentaoatf/src/server/domain"
)

type ExecService struct {
}

func NewExecService() *ExecService {
	return &ExecService{}
}

func (s *ExecService) Exec(build domain.Build) (reply domain.OptResult) {
	// TODO: run testing

	return
}
