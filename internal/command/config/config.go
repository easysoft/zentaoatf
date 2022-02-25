package commandConfig

import (
	"bytes"
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	configUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/config"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	commonUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/display"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	langUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/lang"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	resUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/res"
	serverConfig "github.com/aaronchen2k/deeptest/internal/server/config"
	"github.com/fatih/color"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"reflect"
)

func InitConfig() {
	commConsts.IsRelease = commonUtils.IsRelease()

	commConsts.WorkDir = fileUtils.GetWorkDir()
	commConsts.ConfigPath = commConsts.WorkDir + commConsts.ConfigFile

	if commConsts.Verbose {
		fmt.Printf("\nlaunch %s%s in %s\n", "", commConsts.App, commConsts.WorkDir)
	}
	v := viper.New()
	serverConfig.VIPER = v
	serverConfig.VIPER.SetConfigType("yaml")

	configRes := filepath.Join("res", commConsts.AppServer+".yaml")
	yamlDefault, _ := resUtils.ReadRes(configRes)

	if err := serverConfig.VIPER.ReadConfig(bytes.NewBuffer(yamlDefault)); err != nil {
		panic(fmt.Errorf("读取默认配置文件错误: %w ", err))
	}

	if err := serverConfig.VIPER.Unmarshal(&serverConfig.CONFIG); err != nil {
		panic(fmt.Errorf("同步配置文件错误: %w ", err))
	}

	return
}

func Init() {
	InitConfig()
	serverConfig.InitLog()

	CheckConfigPermission()

	// screen size
	InitScreenSize()

	// internationalization
	i118Utils.Init(commConsts.Language, commConsts.AppServer)

	langUtils.GetExtToNameMap()

	commConsts.ComeFrom = "cmd"
	return
}

func CheckConfigPermission() {
	err := fileUtils.MkDirIfNeeded(commConsts.WorkDir + "conf")
	if err != nil {
		msg := i118Utils.Sprintf("perm_deny", commConsts.WorkDir)
		logUtils.ExecConsolef(color.FgRed, msg)
		os.Exit(0)
	}
}
func InitScreenSize() {
	w, h := display.GetScreenSize()
	consts.ScreenWidth = w
	consts.ScreenHeight = h
}

//func CheckRequestConfig() {
//	conf := configUtils.LoadByProjectPath(commConsts.WorkDir)
//	if conf.Url == "" || conf.Username == "" || conf.Password == "" {
//		stdinUtils.InputForRequest()
//	}
//}

func PrintCurrConfig() {
	logUtils.ExecConsole(color.FgCyan, "\n"+i118Utils.Sprintf("current_config"))
	conf := configUtils.LoadByProjectPath(commConsts.WorkDir)
	val := reflect.ValueOf(conf)
	typeOfS := val.Type()
	for i := 0; i < reflect.ValueOf(conf).NumField(); i++ {
		if !commonUtils.IsWin() && i > 4 {
			break
		}

		val := val.Field(i)
		name := typeOfS.Field(i).Name

		fmt.Printf("  %s: %v \n", name, val.Interface())
	}
}
