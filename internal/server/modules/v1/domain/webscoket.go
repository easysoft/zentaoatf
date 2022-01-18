package serverDomain

import commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"

type WsMsg struct {
	Act         commConsts.ExecCmd `json:"act"`
	Cases       []string           `json:"cases"`
	ProjectPath string             `json:"projectPath"`
}
