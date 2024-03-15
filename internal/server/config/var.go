package serverConfig

import (
	"github.com/spf13/viper"
)

var (
	CONFIG     Config              // 配置
	VIPER      *viper.Viper        // viper
	PermRoutes []map[string]string // 权限路由
)
