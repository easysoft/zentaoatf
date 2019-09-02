package configUtils

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	"github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/display"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
	"sync"
)

func InitConfig() {
	// preference from yaml
	vari.Config = getInst()

	// screen size
	InitScreenSize()

	// internationalization
	i118Utils.InitI118(vari.Config.Language)

}

func InitScreenSize() {
	w, h := display.GetScreenSize()
	vari.ScreenWidth = w
	vari.ScreenHeight = h
}

func getInst() model.Config {
	var once sync.Once
	once.Do(func() {
		vari.Config = model.Config{}
		if fileUtils.FileExist(constant.ConfigFile) {
			buf, _ := ioutil.ReadFile(constant.ConfigFile)
			yaml.Unmarshal(buf, &vari.Config)

			if vari.Config.Version != constant.ConfigVer { // init
				if vari.Config.Language != "en" && vari.Config.Language != "zh" {
					vari.Config.Language = "en"
				}

				SaveConfig(vari.Config.Language, vari.Config.Url, vari.Config.Account, vari.Config.Password)
			}
		} else { // init
			vari.Config = saveEmptyConfig()
		}
	})
	return vari.Config
}

func SaveConfig(language string, url string, account string, password string) error {
	config := ReadCurrConfig()

	config.Version = constant.ConfigVer

	if language != "" {
		config.Language = language
	}
	if url != "" {
		config.Url = url
	}
	if account != "" {
		config.Account = account
	}
	if password != "" {
		config.Password = password
	}

	data, _ := yaml.Marshal(&config)
	ioutil.WriteFile(constant.ConfigFile, data, 0666)

	vari.Config = ReadCurrConfig()
	return nil
}

func PrintCurrConfig() {
	logUtils.PrintToStdOut(i118Utils.I118Prt.Sprintf("current_config"), color.FgCyan)

	val := reflect.ValueOf(vari.Config)
	typeOfS := val.Type()
	for i := 0; i < reflect.ValueOf(vari.Config).NumField(); i++ {
		val := val.Field(i)
		name := typeOfS.Field(i).Name

		if !vari.RunFromCui && (name == "Width" || name == "Height" || name == "WorkHistories") {
			continue
		}
		fmt.Printf("  %s: %v \n", typeOfS.Field(i).Name, val.Interface())
	}
	fmt.Println("")
}

func ReadCurrConfig() model.Config {
	configPath := constant.ConfigFile
	var config model.Config

	if !fileUtils.FileExist(configPath) {
		saveEmptyConfig()
	}
	buf, _ := ioutil.ReadFile(configPath)
	yaml.Unmarshal(buf, &config)

	config.Url = commonUtils.UpdateUrl(config.Url)

	return config
}

func saveEmptyConfig() model.Config {
	config := model.Config{Version: constant.ConfigVer, Language: "en", Url: "", Account: "", Password: ""}

	data, _ := yaml.Marshal(&config)
	ioutil.WriteFile(constant.ConfigFile, data, 0666)

	return config
}
