package commConsts

import "github.com/aaronchen2k/deeptest/internal/comm/domain"

var (
	ProjectConfig domain.ProjectConf

	Language  = "zh"
	Verbose   bool
	IsRelease bool
	ExeDir    string
	WorkDir   string

	RequestType   ZentaoRequestType
	ZenTaoVersion string
	SessionId     string
	SessionVar    string
	RequestFix    string
)
