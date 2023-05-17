package execHelper

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	websocketHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/websocket"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
)

func PrintErrMsg(key string, err error, wsMsg *websocket.Message) (
	stdOutput string, errOutput string) {

	if commConsts.ExecFrom == commConsts.FromClient {
		websocketHelper.SendOutputMsg(err.Error(), "", iris.Map{"key": key}, wsMsg)
	}
	logUtils.ExecConsolef(color.FgRed, err.Error())
	logUtils.ExecFilef(err.Error())

	return "", err.Error()
}
