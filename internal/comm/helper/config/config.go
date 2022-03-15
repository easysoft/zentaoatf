package configHelper

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	commonUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	stringUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/string"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"gopkg.in/ini.v1"
	"path/filepath"
	"reflect"
)

func LoadByWorkspacePath(workspacePath string) (config commDomain.WorkspaceConf) {
	pth := filepath.Join(workspacePath, commConsts.ConfigDir, commConsts.ConfigFile)
	ini.MapTo(&config, pth)

	config.Url = commonUtils.AddSlashForUrl(config.Url)

	return config
}

func UpdateSite(site model.Site, workspacePath string) (err error) {
	config := LoadByWorkspacePath(workspacePath)

	config.Url = site.Url
	config.Username = site.Username
	config.Password = site.Password

	SaveToFile(config, workspacePath)

	return
}

func SaveConfig(config commDomain.WorkspaceConf, workspacePath string) (err error) {
	SaveToFile(config, workspacePath)

	return
}

func ReadFromFile(workspacePath string) (config commDomain.WorkspaceConf) {
	pth := filepath.Join(workspacePath, commConsts.ConfigDir, commConsts.ConfigFile)
	ini.MapTo(&config, pth)

	config.Url = commonUtils.AddSlashForUrl(config.Url)

	return config
}

func SaveToFile(config commDomain.WorkspaceConf, workspacePath string) (err error) {
	pth := filepath.Join(workspacePath, commConsts.ConfigDir, commConsts.ConfigFile)
	fileUtils.MkDirIfNeeded(filepath.Dir(pth))

	config.Version = commConsts.ConfigVersion

	cfg := ini.Empty()
	cfg.ReflectFrom(&config)

	cfg.SaveTo(pth)
	logUtils.Infof("Success to update config file %s.", pth)

	return nil
}

func GetFieldVal(config commDomain.WorkspaceConf, key string) string {
	key = stringUtils.UcFirst(key)

	immutable := reflect.ValueOf(config)
	val := immutable.FieldByName(key).String()

	return val
}

func SetFieldVal(config *commDomain.WorkspaceConf, key string, val string) string {
	key = stringUtils.UcFirst(key)

	mutable := reflect.ValueOf(config).Elem()
	mutable.FieldByName(key).SetString(val)

	return val
}
