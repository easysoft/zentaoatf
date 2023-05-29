package commandConfig

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	configHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/config"
	langHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/lang"
	serverConfig "github.com/easysoft/zentaoatf/internal/server/config"
	"github.com/easysoft/zentaoatf/pkg/consts"
	commonUtils "github.com/easysoft/zentaoatf/pkg/lib/common"
	"github.com/easysoft/zentaoatf/pkg/lib/display"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	resUtils "github.com/easysoft/zentaoatf/pkg/lib/res"
	"github.com/fatih/color"
	"github.com/spf13/viper"
)

func InitConfig() {
	commConsts.IsRelease = commonUtils.IsRelease()

	commConsts.WorkDir = fileUtils.GetWorkDir()
	commConsts.ZtfDir = fileUtils.GetZTFDir()

	if !commConsts.IsRelease {
		log.Println("WorkDir=" + commConsts.WorkDir)
		log.Println("ZtfDir=" + commConsts.ZtfDir)
	}

	commConsts.ConfigPath = filepath.Join(commConsts.WorkDir, commConsts.ConfigDir, commConsts.ConfigFile)
	if commConsts.IsRelease {
		commConsts.ConfigPath = filepath.Join(commConsts.ZtfDir, commConsts.ConfigDir, commConsts.ConfigFile)
	}

	config := configHelper.LoadByConfigPath(commConsts.ConfigPath)
	if config.Language != "" {
		commConsts.Language = config.Language
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
	InitLog()

	CheckConfigPermission()

	InitScreenSize()

	i118Utils.Init(commConsts.Language, commConsts.AppServer)

	langHelper.GetExtToNameMap()
	langHelper.GetEditorExtToLangMap()

	commConsts.ExecFrom = commConsts.FromCmd
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

func PrintCurrConfig() {
	logUtils.ExecConsole(color.FgCyan, "\n"+i118Utils.Sprintf("current_config"))
	conf := configHelper.LoadByWorkspacePath(commConsts.ZtfDir)
	val := reflect.ValueOf(conf)
	typeOfS := val.Type()
	for i := 0; i < reflect.ValueOf(conf).NumField(); i++ {
		if !commonUtils.IsWin() && i >= 4 {
			break
		}

		val := val.Field(i)
		name := typeOfS.Field(i).Name

		fmt.Printf("  %s: %v \n", name, val.Interface())
	}
}
