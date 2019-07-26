package utils

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/jroimartin/gocui"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
)

type Config struct {
	Language string
	WorkDir  string

	Width  int
	Height int
}

var Conf Config

func InitConfig() {
	// config from yaml
	Conf = getInst()

	// screen size
	InitScreenSize()

	// internationalization
	InitI118(Conf.Language)

	if strings.Index(os.Args[0], "atf") > -1 && Conf.Language != "" && os.Args[1] != "set" {
		PrintConfig()
	}
}

func Set(param string, val string) {
	buf, _ := ioutil.ReadFile(ConfFile)
	yaml.Unmarshal(buf, &Conf)

	if param == "lang" {
		Conf.Language = val
		color.Blue(I118Prt.Sprintf("set_config", I118Prt.Sprintf("lang"), I118Prt.Sprintf(Conf.Language)))
	} else if param == "workDir" {
		val = getWorkDir(val)

		Conf.WorkDir = val
		color.Blue(I118Prt.Sprintf("set_config", I118Prt.Sprintf("workDir"), Conf.WorkDir))
	}
	data, _ := yaml.Marshal(&Conf)
	ioutil.WriteFile(ConfFile, data, 0666)
}

func getInst() Config {
	var once sync.Once
	once.Do(func() {
		Conf = Config{}
		if FileExist(ConfFile) {
			buf, _ := ioutil.ReadFile(ConfFile)
			yaml.Unmarshal(buf, &Conf)
		} else { // init
			Conf.Language = "en"
			Conf.WorkDir = getWorkDir("./")

			data, _ := yaml.Marshal(&Conf)
			ioutil.WriteFile(ConfFile, data, 0666)
		}
	})
	return Conf
}

func PrintConfig() {
	color.Blue(I118Prt.Sprintf("current_config", ""))

	val := reflect.ValueOf(Conf)
	typeOfS := val.Type()
	for i := 0; i < reflect.ValueOf(Conf).NumField(); i++ {
		val := val.Field(i)
		fmt.Printf("  %s: %v \n", typeOfS.Field(i).Name, val.Interface())
	}
}

func PrintConfigToView(v *gocui.View) {
	fmt.Fprintln(v, color.BlueString(I118Prt.Sprintf("current_config", "")))

	val := reflect.ValueOf(Conf)
	typeOfS := val.Type()
	for i := 0; i < reflect.ValueOf(Conf).NumField(); i++ {
		val := val.Field(i)
		fmt.Fprintln(v, fmt.Sprintf("  %s: %v", typeOfS.Field(i).Name, val.Interface()))
	}
}

func getWorkDir(path string) string {
	if path == "./" {
		path, _ = filepath.Abs(`.`)
		if !IsRelease() { // remove 'bin' on dev mode
			path = path + string(os.PathSeparator)
		}
	}

	return path
}
