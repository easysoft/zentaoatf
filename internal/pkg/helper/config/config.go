package configHelper

import (
	"path/filepath"
	"reflect"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	"github.com/easysoft/zentaoatf/internal/server/core/dao"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	commonUtils "github.com/easysoft/zentaoatf/pkg/lib/common"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	stringUtils "github.com/easysoft/zentaoatf/pkg/lib/string"
	"gopkg.in/ini.v1"
)

func LoadBySite(site model.Site) (config commDomain.WorkspaceConf) {
	config = commDomain.WorkspaceConf{
		Url:      site.Url,
		Username: site.Username,
		Password: site.Password,
	}

	config.Url = commonUtils.AddSlashForUrl(config.Url)

	return config
}

func LoadByConfigPath(configPath string) (config commDomain.WorkspaceConf) {
	ini.MapTo(&config, configPath)
	config.Url = commonUtils.AddSlashForUrl(config.Url)
	return config
}
func LoadByWorkspacePath(workspacePath string) (config commDomain.WorkspaceConf) {
	if workspacePath == "" {
		//从db获取interpreter的路径
		GetInterpreterConfig(&config)
		return
	}

	pth := filepath.Join(workspacePath, commConsts.ConfigDir, commConsts.ConfigFile)

	return LoadByConfigPath(pth)
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

func GetFieldVal(config commDomain.WorkspaceConf, key string) (val string) {
	key = stringUtils.UcFirst(key)

	immutable := reflect.ValueOf(config)
	value := immutable.FieldByName(key)

	if value.IsValid() {
		val = value.String()
	}

	return
}

func SetFieldVal(config *commDomain.WorkspaceConf, key string, val string) string {
	key = stringUtils.UcFirst(key)

	mutable := reflect.ValueOf(config).Elem()
	mutable.FieldByName(key).SetString(val)

	return val
}

func GetInterpreterConfig(config *commDomain.WorkspaceConf) (err error) {
	interps := []model.Interpreter{}
	dao.GetDB().Model(&model.Interpreter{}).Where("NOT deleted")
	err = dao.GetDB().Find(&interps).Error
	mp := map[string]string{}

	for _, item := range interps {
		mp[item.Lang] = item.Path
	}

	if config.Language == "" {
		config.Language = commConsts.LanguageZh
	}
	config.Javascript = mp["javascript"]
	config.Lua = mp["lua"]
	config.Perl = mp["perl"]
	config.Php = mp["php"]
	config.Python = mp["python"]
	config.Ruby = mp["ruby"]
	config.Tcl = mp["tcl"]
	config.Autoit = mp["autoit"]
	return
}

func UpdateAllInterpreterConfig() {
	var workspaces []model.Workspace
	dao.GetDB().Model(&model.Workspace{}).
		Where("NOT deleted").
		Find(&workspaces)

	for _, item := range workspaces {
		if item.Type != commConsts.ZTF {
			continue
		}

		UpdateInterpreterConfig(item)
	}
}

func UpdateInterpreterConfig(workspace model.Workspace) (err error) {
	interps := []model.Interpreter{}
	dao.GetDB().Model(&model.Interpreter{}).Where("NOT deleted")
	err = dao.GetDB().Find(&interps).Error
	mp := map[string]string{}

	for _, item := range interps {
		mp[item.Lang] = item.Path
	}

	conf := ReadFromFile(workspace.Path)
	if conf.Language == "" {
		conf.Language = commConsts.LanguageZh
	}

	conf.Javascript = mp["javascript"]
	conf.Lua = mp["lua"]
	conf.Perl = mp["perl"]
	conf.Php = mp["php"]
	conf.Python = mp["python"]
	conf.Go = mp["go"]
	conf.Ruby = mp["ruby"]
	conf.Tcl = mp["tcl"]
	conf.Autoit = mp["autoit"]

	SaveToFile(conf, workspace.Path)

	return
}
