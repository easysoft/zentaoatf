package commandConfig

import (
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	configUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/config"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	commonUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/display"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	stdinUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/stdin"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
	"github.com/fatih/color"
	"os"
	"reflect"
)

type ConfigCtrl struct {
	ProjectRepo *repo.ProjectRepo `inject:""`
}

func CheckConfigPermission() {
	//err := syscall.Access(vari.ExeDir, syscall.O_RDWR)

	err := fileUtils.MkDirIfNeeded(commConsts.ExeDir + "conf")
	if err != nil {
		msg := i118Utils.Sprintf("perm_deny", commConsts.ExeDir)
		logUtils.ExecConsolef(color.FgRed, msg)
		os.Exit(0)
	}
}

func InitScreenSize() {
	w, h := display.GetScreenSize()
	consts.ScreenWidth = w
	consts.ScreenHeight = h
}

func CheckRequestConfig() {
	conf := configUtils.LoadByProjectPath(commConsts.WorkDir)
	if conf.Url == "" || conf.Username == "" || conf.Password == "" {
		stdinUtils.InputForRequest()
	}
}

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
