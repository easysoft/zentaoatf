package service

import (
	serverModel "github.com/easysoft/zentaoatf/src/server/domain"
	commonUtils "github.com/easysoft/zentaoatf/src/server/utils/common"
	serverConst "github.com/easysoft/zentaoatf/src/server/utils/const"
	"github.com/easysoft/zentaoatf/src/service/client"
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
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
		sysInfo = commonUtils.GetSysInfo()
	}

	// send request
	zentaoService.GetConfig(vari.Config.Url)

	url := vari.Config.Url + zentaoUtils.GenApiUri("agent", "heartBeat", "")
	data := map[string]interface{}{"type": vari.Platform, "sys": sysInfo}

	status := serverConst.VmActive
	if isBusy {
		status = serverConst.VmBusy
	}
	data["status"] = status

	_, ok := client.PostObject(url, data, false)
	if ok {
		logUtils.PrintTo("heart beat success")
	} else {
		logUtils.PrintTo("heart beat fail")
	}

	return
}
