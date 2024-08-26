package commDomain

import (
	"github.com/kataras/iris/v12"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
)

type WsResp struct {
	Msg       string                   `json:"msg"`
	IsRunning string                   `json:"isRunning,omitempty"`
	Category  commConsts.WsMsgCategory `json:"category"`

	Info iris.Map `json:"info,omitempty"`
}

type MqMsg struct {
	Namespace string `json:"namespace"`
	Room      string `json:"room"`
	Event     string `json:"event"`
	Content   string `json:"content"`
}
