package serverConfig

import (
	"bytes"
	"fmt"
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	"github.com/easysoft/zentaoatf/pkg/lib/common"
	"github.com/easysoft/zentaoatf/pkg/lib/file"
	"github.com/easysoft/zentaoatf/pkg/lib/res"
	"github.com/spf13/viper"
	"path/filepath"
)

func Init() {
	commConsts.IsRelease = commonUtils.IsRelease()

	//commConsts.ZtfDir = fileUtils.GetExeDir(commConsts.WorkDir)
	commConsts.WorkDir = GetServerWorDir()

	if commConsts.Verbose {
		fmt.Printf("\nlaunch %s%s in %s\n", "", commConsts.App, commConsts.WorkDir)
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

func GetServerWorDir() (ret string) {
	home, _ := fileUtils.GetUserHome()
	ret = filepath.Join(home, commConsts.App)

	ret = fileUtils.AddFilePathSepIfNeeded(ret)
	fileUtils.MkDirIfNeeded(ret)

	return
}
