package agentViper

import (
	"bytes"
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/agent/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	myZap "github.com/aaronchen2k/deeptest/internal/pkg/core/zap"
	"github.com/fsnotify/fsnotify"
	"github.com/snowlyg/helper/dir"
	"github.com/spf13/viper"
)

// Init 初始化系统配置
// - 第一次初始化系统配置，会自动生成配置文件 config.yaml 和 casbin 的规则文件 rbac_model.conf
// - 热监控系统配置项，如果发生变化会重写配置文件内的配置项
func Init() {
	config := consts.ConfigFileAgent
	fmt.Printf("配置文件路径为%s\n", config)

	v := viper.New()
	agentConsts.VIPER = v
	agentConsts.VIPER.SetConfigType("yaml")

	if !dir.IsExist(config) { //没有配置文件，写入默认配置
		var yamlDefault = []byte(`
system:
 level: debug # debug,release,test
 addr: "127.0.0.1:8085"
 time-format: "2006-01-02 15:04:05"
zap:
 level: info
 format: console
 prefix: '[OP-ONLINE]'
 director: log
 link-name: latest_log
 show-line: true
 encode-level: LowercaseColorLevelEncoder
 stacktrace-key: stacktrace
 log-in-console: true`)

		if err := agentConsts.VIPER.ReadConfig(bytes.NewBuffer(yamlDefault)); err != nil {
			panic(fmt.Errorf("读取默认配置文件错误: %w ", err))
		}

		if err := agentConsts.VIPER.Unmarshal(&agentConsts.CONFIG); err != nil {
			panic(fmt.Errorf("同步配置文件错误: %w ", err))
		}

		if err := agentConsts.VIPER.WriteConfigAs(config); err != nil {
			panic(fmt.Errorf("写入配置文件错误: %w ", err))
		}
		return
	}

	// 存在配置文件，读取配置文件内容
	agentConsts.VIPER.SetConfigFile(config)
	err := agentConsts.VIPER.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("读取配置错误: %w ", err))
	}

	// 监控配置文件变化
	agentConsts.VIPER.WatchConfig()
	agentConsts.VIPER.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置发生变化:", e.Name)
		if err := agentConsts.VIPER.Unmarshal(&agentConsts.CONFIG); err != nil {
			fmt.Println(err)
		}
	})
	if err := v.Unmarshal(&agentConsts.CONFIG); err != nil {
		fmt.Println(err)
	}

	myZap.ZapInst = agentConsts.CONFIG.Zap
}
