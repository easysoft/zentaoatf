package utils

import (
	"github.com/easysoft/zentaoatf/src/model"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func SaveConfig(url string, entityType string, entityVal string, langType string, singleFile bool, name string) error {
	config := model.Config{Url: url, EntityType: entityType, LangType: langType, SingleFile: singleFile, ProjectName: name}

	config.EntityType = entityType
	config.EntityVal = entityVal

	data, _ := yaml.Marshal(&config)
	ioutil.WriteFile(Prefer.WorkDir+ConfigFile, data, 0666)

	return nil
}

func SaveEmptyConfig() error {
	SaveConfig("", "", "", "", false, "")

	return nil
}

func ReadProjectConfig(projectPath string) model.Config {
	return ReadConfig(projectPath)
}

func ReadCurrConfig() model.Config {
	return ReadConfig(Prefer.WorkDir)
}

func ReadConfig(dir string) model.Config {
	configPath := dir + ConfigFile
	var config model.Config
	if !FileExist(configPath) {
		SaveEmptyConfig()
	}
	buf, _ := ioutil.ReadFile(configPath)
	yaml.Unmarshal(buf, &config)

	return config
}
