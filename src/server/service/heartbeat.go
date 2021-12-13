package service

import (
	serverModel "github.com/easysoft/zentaoatf/src/server/domain"
	"github.com/easysoft/zentaoatf/src/server/utils/common"
	serverConst "github.com/easysoft/zentaoatf/src/server/utils/const"
	"github.com/easysoft/zentaoatf/src/service/client"
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	zentaoUtils "github.com/easysoft/zentaoatf/src/utils/zentao"
)

var (
	sysInfo serverModel.SysInfo
)

type HeartBeatService struct {
}

func NewHeartBeatService() *HeartBeatService {
	return &HeartBeatService{}
}

func (s *HeartBeatService) HeartBeat(isBusy bool) {
	if sysInfo.OsName == "" {
		sysInfo = serverUtils.GetSysInfo()
	}

	// send request
	zentaoService.GetConfig(vari.Config.Url)

	url := vari.Config.Url + zentaoUtils.GenApiUri("agent", "heartbeat", "")
	data := map[string]interface{}{"type": vari.Platform, "sys": sysInfo}

	status := serverConst.VmActive
	if isBusy {
		status = serverConst.VmBusy
	}
	data["status"] = status

	_, ok := client.PostObject(url, data, false)
	if ok {
		logUtils.PrintTo(i118Utils.Sprintf("success_heart_beat"))
	} else {
		logUtils.PrintTo(i118Utils.Sprintf("fail_heart_beat"))
	}

	return
}
