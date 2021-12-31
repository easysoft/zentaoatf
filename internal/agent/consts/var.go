package agentConsts

import (
	agentConfig "github.com/aaronchen2k/deeptest/internal/agent/config"
	"github.com/spf13/viper"
)

var (
	CONFIG agentConfig.Config // 配置
	VIPER  *viper.Viper       // viper
)
