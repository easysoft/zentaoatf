package configUtils

import (
	"github.com/easysoft/zentaoatf/src/model"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/display"
	"github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func SaveConfig(dir string, url string, entityType string, entityVal string,
	productId int, projectId int, langType string, independentFile bool, name string,
	account string, password string) error {
	config := model.Config{Url: url, EntityType: entityType,
		ProductId: productId, ProjectId: projectId, LangType: langType, IndependentFile: independentFile, ProjectName: name,
		Account: account, Password: password}

	config.EntityType = entityType
	config.EntityVal = entityVal

	if dir == "" {
		dir = vari.Prefer.WorkDir
	}

	data, _ := yaml.Marshal(&config)
	ioutil.WriteFile(dir+constant.ConfigFile, data, 0666)

	return nil
}

func ReadProjectConfig(projectPath string) model.Config {
	return ReadConfig(projectPath)
}

func ReadCurrConfig() model.Config {
	return ReadConfig(vari.Prefer.WorkDir)
}

func ReadConfig(dir string) model.Config {
	configPath := dir + constant.ConfigFile
	var config model.Config

	if !fileUtils.FileExist(configPath) {
		saveEmptyConfig(dir)
	}
	buf, _ := ioutil.ReadFile(configPath)
	yaml.Unmarshal(buf, &config)

	return config
}

func saveEmptyConfig(dir string) error {
	SaveConfig(dir, "", "", "", 0, 0, "", false, "",
		"", "")

	return nil
}

func InitScreenSize() {
	w, h := display.GetScreenSize()
	vari.Prefer.Width = w
	vari.Prefer.Height = h
}
