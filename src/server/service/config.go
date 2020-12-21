package service

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/server/domain"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
)

var ()

type ConfigService struct {
}

func NewConfigService() *ConfigService {
	return &ConfigService{}
}

func (s *ConfigService) Update(req domain.ReqData) {
	conf := model.Config{}

	reqStr, _ := json.Marshal(req.Data)
	err := json.Unmarshal(reqStr, &conf)
	if err != nil {
		logUtils.PrintTo(fmt.Sprintf("error: %v", err))
		return
	}

	if conf.Version != 0 {
		vari.Config.Version = conf.Version
	}
	if conf.Language != "" {
		vari.Config.Language = conf.Language
	}
	if conf.Url != "" {
		vari.Config.Url = conf.Url
	}
	if conf.Account != "" {
		vari.Config.Account = conf.Account
	}
	if conf.Password != "" {
		vari.Config.Password = conf.Password
	}
	if conf.Javascript != "" {
		vari.Config.Javascript = conf.Javascript
	}
	if conf.Lua != "" {
		vari.Config.Lua = conf.Lua
	}
	if conf.Perl != "" {
		vari.Config.Perl = conf.Perl
	}
	if conf.Php != "" {
		vari.Config.Php = conf.Php
	}
	if conf.Python != "" {
		vari.Config.Python = conf.Python
	}
	if conf.Ruby != "" {
		vari.Config.Ruby = conf.Ruby
	}
	if conf.Tcl != "" {
		vari.Config.Tcl = conf.Tcl
	}
	if conf.Lua != "" {
		vari.Config.Autoit = conf.Autoit
	}

	configUtils.SaveConfig(vari.Config)
}
