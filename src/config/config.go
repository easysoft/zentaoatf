package config

import (
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
	"sync"
)

type ScreenSizeStruct struct {
	width  int
	height int
}

type Config struct {
	Language string
}

var config Config

func InitConfig() {
	config = GetInst()

	// language
	p := GetI118(config.Language)
	color.Blue(p.Sprintf("current_config", ""))

	val := reflect.ValueOf(config)
	typeOfS := val.Type()
	for i := 0; i < reflect.ValueOf(config).NumField(); i++ {
		fmt.Printf("  %s: %v \n", typeOfS.Field(i).Name, val.Field(i).Interface())
	}
}

func Set(param string, val string) {
	buf, _ := ioutil.ReadFile(ConfFile)
	yaml.Unmarshal(buf, &config)

	if param == "lang" {
		config.Language = val

		data, _ := yaml.Marshal(&config)
		ioutil.WriteFile(ConfFile, data, 0666)

		config := GetInst()
		p := GetI118(GetInst().Language)
		color.Blue(p.Sprintf("set_config", p.Sprintf("lang"), p.Sprintf(config.Language)))
	}

}

func GetInst() Config {
	var once sync.Once
	once.Do(func() {
		config = Config{}
		buf, _ := ioutil.ReadFile(ConfFile)
		yaml.Unmarshal(buf, &config)
	})
	return config
}
