package config

import (
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
)

type Config struct {
	Language string
}

var config Config

func GetInst() Config {
	var once sync.Once
	once.Do(func() {
		config = Config{}
		buf, _ := ioutil.ReadFile(ConfigFile)
		yaml.Unmarshal(buf, &config)
	})
	return config
}

func InitConfig() {
	config := Config{}
	buf, _ := ioutil.ReadFile("config.yaml")
	yaml.Unmarshal(buf, &config)

	fmt.Println(color.BlueString("current config %+v", config))

	p := GetI118(config.Language)
	fmt.Println(p.Sprintf("HELLO_1", "Peter"))
}

func Set(param string, val string) {
	buf, _ := ioutil.ReadFile("config.yaml")
	yaml.Unmarshal(buf, &config)

	if param == "lang" {
		config.Language = val

		data, _ := yaml.Marshal(&config)
		ioutil.WriteFile(ConfigFile, data, 0666)
	}

}
