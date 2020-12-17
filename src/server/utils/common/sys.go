package serverUtils

import (
	"bytes"
	"errors"
	serverModel "github.com/easysoft/zentaoatf/src/server/model"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

func GetSysInfo() (info serverModel.SysInfo) {
	info.AgentDir = vari.AgentDir

	info.SysArch = runtime.GOARCH
	info.SysCores = runtime.GOMAXPROCS(0)

	info.OsType = runtime.GOOS
	info.OsName, _ = os.Hostname()

	envs := os.Environ()
	for _, env := range envs {
		if strings.Index(env, "LC_CTYPE=") > -1 { // LC_CTYPE=zh_CN.UTF-8
			info.Lang = strings.Split(env, "=")[1]
		}
	}

	return
}

func GetUserHome() (string, error) {
	user, err := user.Current()
	if nil == err {
		return user.HomeDir, nil
	}

	// cross compile support

	if "windows" == runtime.GOOS {
		return homeWindows()
	}

	// Unix-like system, so just assume Unix
	return homeUnix()
}

func homeUnix() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

func homeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}
