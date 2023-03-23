package service

import (
	"encoding/json"
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"

	hostHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/host"
	serverConfig "github.com/easysoft/zentaoatf/internal/server/config"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	commonUtils "github.com/easysoft/zentaoatf/pkg/lib/common"
	httpUtils "github.com/easysoft/zentaoatf/pkg/lib/http"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
)

type HeartbeatService struct {
}

func NewHeartbeatService() *HeartbeatService {
	return &HeartbeatService{}
}

func (s *HeartbeatService) Heartbeat() {
	_, macObj := commonUtils.GetIp()
	req := serverDomain.HostHeartbeatReq{
		MacAddress: macObj.String(),
	}

	url := hostHelper.GenApiUrl("virtual/notifyHost", nil, serverConfig.CONFIG.Host)
	respBytes, err := httpUtils.Post(url, req)
	ok := err == nil

	if ok {
		resp := serverDomain.HostResponse{}
		err = json.Unmarshal(respBytes, &resp)
		if err != nil {
			return
		}

		respDataBytes, _ := json.Marshal(resp.Data)

		respObj := serverDomain.HostHeartbeatResp{}
		err := json.Unmarshal(respDataBytes, &respObj)
		if err == nil && respObj.Token != "" {
			serverConfig.CONFIG.AuthToken = respObj.Token
		}
		if serverConfig.CONFIG.Server == "" && respObj.Server != "" {
			serverConfig.CONFIG.Server = respObj.Server
			serverConfig.CONFIG.Server = httpUtils.AddSepIfNeeded(serverConfig.CONFIG.Server)
		}
	}

	if serverConfig.CONFIG.AuthToken == "" {
		ok = false
	}

	if ok {
		if serverConfig.CONFIG.AuthToken != "" {
			commConsts.ExecFrom = commConsts.FromZentao
		}
		logUtils.Info(i118Utils.I118Prt.Sprintf("success_to_register", url))
	} else {
		logUtils.Info(i118Utils.I118Prt.Sprintf("fail_to_register", url, string(respBytes)))
	}
}
