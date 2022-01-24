package repo

import (
	"gorm.io/gorm"
)

type TestExecRepo struct {
	DB          *gorm.DB     `inject:""`
	ProjectRepo *ProjectRepo `inject:""`
}

func NewTestExecRepo() *TestExecRepo {
	return &TestExecRepo{}
}
