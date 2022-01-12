package serverConfig

import (
	"github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/aaronchen2k/deeptest/internal/comm/domain"
	commonUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"gopkg.in/ini.v1"
	"path"
	"path/filepath"
)

func ReadConfig(projectPath string) (config domain.ProjectConfig) {
	pth := filepath.Join(projectPath, commConsts.ConfigDir, commConsts.ConfigFile)
	ini.MapTo(&config, pth)

	config.Url = commonUtils.AddSlashForUrl(config.Url)

	return config
}

func SaveConfig(config domain.ProjectConfig, projectPath string) (err error) {
	pth := filepath.Join(projectPath, commConsts.ConfigDir, commConsts.ConfigFile)
	fileUtils.MkDirIfNeeded(path.Dir(pth))

	config.Version = CONFIG.System.Version

	cfg := ini.Empty()
	cfg.ReflectFrom(&config)

	cfg.SaveTo(pth)
	logUtils.Infof("Successfully update config file %s.", pth)

	return nil
}
