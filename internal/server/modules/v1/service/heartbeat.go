package service

import (
	"encoding/json"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	serverConfig "github.com/easysoft/zentaoatf/internal/server/config"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	dateUtils "github.com/easysoft/zentaoatf/pkg/lib/date"
	httpUtils "github.com/easysoft/zentaoatf/pkg/lib/http"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	"time"
)

type HeartbeatService struct {
}

func NewHeartbeatService() *HeartbeatService {
	return &HeartbeatService{}
}

func (s *HeartbeatService) Heartbeat() {
	req := serverDomain.ZentaoHeartbeatReq{
		Ip:   serverConfig.CONFIG.Ip,
		Port: serverConfig.CONFIG.Port,
	}

	if serverConfig.CONFIG.AuthToken == "" || serverConfig.CONFIG.ExpiredDate.Unix() < time.Now().Unix() { // re-apply token using secret
		req.Secret = serverConfig.CONFIG.Secret
	}

	url := zentaoHelper.GenApiUrl("ztf/heartbeat", nil, serverConfig.CONFIG.Server)
	respBytes, err := httpUtils.Post(url, req)
	ok := err == nil

	if ok {
		respObj := serverDomain.ZentaoHeartbeatResp{}
		err := json.Unmarshal(respBytes, &respObj)
		if err == nil && respObj.Token != "" {
			respObj.ExpiredDate, _ = dateUtils.UnitToDate(respObj.ExpiredTimeUnix)
			serverConfig.CONFIG.AuthToken = respObj.Token
			serverConfig.CONFIG.ExpiredDate = respObj.ExpiredDate
		}
	}

	if serverConfig.CONFIG.AuthToken == "" {
		ok = false
	}

	if ok {
		logUtils.Info(i118Utils.I118Prt.Sprintf("success_to_register", url))
	} else {
		logUtils.Info(i118Utils.I118Prt.Sprintf("fail_to_register", url, respBytes))
	}
}
