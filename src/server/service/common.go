package service

import (
	commonUtils "github.com/easysoft/zentaoatf/src/server/utils/common"
	"github.com/easysoft/zentaoatf/src/service/client"
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	zentaoUtils "github.com/easysoft/zentaoatf/src/utils/zentao"
)

type CommonService struct {
}

func NewCommonService() *CommonService {
	return &CommonService{}
}

func (s *CommonService) HeartBeat() {
	sysInfo := commonUtils.GetSysInfo()

	// send request
	zentaoService.GetConfig(vari.Config.Url)

	url := vari.Config.Url + zentaoUtils.GenApiUri("agent", "heartBeat", "")
	data := map[string]interface{}{"sys": sysInfo}

	_, ok := client.PostObject(url, data, false)
	if ok {
		logUtils.PrintTo("heart beat success")
	} else {
		logUtils.PrintTo("heart beat fail")
	}

	return
}
