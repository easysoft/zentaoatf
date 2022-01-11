package serverConfig

import (
	"bytes"
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/res"
	"github.com/spf13/viper"
	"path/filepath"
)

func Init() {
	CONFIG.System.IsRelease = commonUtils.IsRelease()

	CONFIG.System.WorkDir = fileUtils.GetWorkDir()
	CONFIG.System.ExeDir = fileUtils.GetExeDir(CONFIG.System.WorkDir)

	if CONFIG.System.Verbose {
		fmt.Printf("launch %s%s in %s\n",
			CONFIG.System.ExeDir, commConsts.App, CONFIG.System.WorkDir)
	}

	v := viper.New()
	VIPER = v
	VIPER.SetConfigType("yaml")

	configRes := filepath.Join("res", commConsts.AppServer+".yaml")

	yamlDefault, _ := resUtils.ReadRes(configRes)

	if err := VIPER.ReadConfig(bytes.NewBuffer(yamlDefault)); err != nil {
		panic(fmt.Errorf("读取默认配置文件错误: %w ", err))
	}

	if err := VIPER.Unmarshal(&CONFIG); err != nil {
		panic(fmt.Errorf("同步配置文件错误: %w ", err))
	}
	return
}
