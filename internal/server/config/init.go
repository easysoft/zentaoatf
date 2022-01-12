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
	commConsts.IsRelease = commonUtils.IsRelease()

	commConsts.WorkDir = fileUtils.GetWorkDir()
	commConsts.ExeDir = fileUtils.GetExeDir(commConsts.WorkDir)

	if commConsts.Verbose {
		fmt.Printf("launch %s%s in %s\n",
			commConsts.ExeDir, commConsts.App, commConsts.WorkDir)
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
