package configUtils

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	commonUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	stringUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/string"
	"gopkg.in/ini.v1"
	"path"
	"path/filepath"
	"reflect"
)

func LoadByProjectPath(projectPath string) (config commDomain.ProjectConf) {
	pth := filepath.Join(projectPath, commConsts.ConfigDir, commConsts.ConfigFile)
	ini.MapTo(&config, pth)

	config.Url = commonUtils.AddSlashForUrl(config.Url)

	return config
}

func SaveConfig(config commDomain.ProjectConf, projectPath string) (err error) {
	SaveToFile(config, projectPath)

	return
}

func ReadFromFile(projectPath string) (config commDomain.ProjectConf) {
	pth := filepath.Join(projectPath, commConsts.ConfigDir, commConsts.ConfigFile)
	ini.MapTo(&config, pth)

	config.Url = commonUtils.AddSlashForUrl(config.Url)

	return config
}

func SaveToFile(config commDomain.ProjectConf, projectPath string) (err error) {
	pth := filepath.Join(projectPath, commConsts.ConfigDir, commConsts.ConfigFile)
	fileUtils.MkDirIfNeeded(path.Dir(pth))

	config.Version = commConsts.ConfigVersion

	cfg := ini.Empty()
	cfg.ReflectFrom(&config)

	cfg.SaveTo(pth)
	logUtils.Infof("Successfully update config file %s.", pth)

	return nil
}

func GetFieldVal(config commDomain.ProjectConf, key string) string {
	key = stringUtils.UcFirst(key)

	immutable := reflect.ValueOf(config)
	val := immutable.FieldByName(key).String()

	return val
}
