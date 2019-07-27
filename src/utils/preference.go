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

type Preference struct {
	Language string
	WorkDir  string

	Width  int
	Height int
}

var Conf Preference

func InitPreference() {
	// preference from yaml
	Conf = getInst()

	// screen size
	InitScreenSize()

	// internationalization
	InitI118(Conf.Language)

	if strings.Index(os.Args[0], "atf") > -1 && (len(os.Args) > 1 && os.Args[1] != "set") {
		PrintPreference()
	}
}

func SetPreference(param string, val string, dumb bool) {
	buf, _ := ioutil.ReadFile(PreferenceFile)
	yaml.Unmarshal(buf, &Conf)

	if param == "lang" {
		Conf.Language = val
		if !dumb {
			color.Blue(I118Prt.Sprintf("set_preference", I118Prt.Sprintf("lang"), I118Prt.Sprintf(Conf.Language)))
		}
	} else if param == "workDir" {
		val = convertWorkDir(val)

		Conf.WorkDir = val
		if !dumb {
			color.Blue(I118Prt.Sprintf("set_preference", I118Prt.Sprintf("workDir"), Conf.WorkDir))
		}
	}
	data, _ := yaml.Marshal(&Conf)
	ioutil.WriteFile(PreferenceFile, data, 0666)
}

func getInst() Preference {
	var once sync.Once
	once.Do(func() {
		Conf = Preference{}
		if FileExist(PreferenceFile) {
			buf, _ := ioutil.ReadFile(PreferenceFile)
			yaml.Unmarshal(buf, &Conf)
		} else { // init
			Conf.Language = "en"
			Conf.WorkDir = convertWorkDir("./")

			data, _ := yaml.Marshal(&Conf)
			ioutil.WriteFile(PreferenceFile, data, 0666)
		}
	})
	return Conf
}

func PrintPreference() {
	color.Blue(I118Prt.Sprintf("current_preference", ""))

	val := reflect.ValueOf(Conf)
	typeOfS := val.Type()
	for i := 0; i < reflect.ValueOf(Conf).NumField(); i++ {
		val := val.Field(i)
		fmt.Printf("  %s: %v \n", typeOfS.Field(i).Name, val.Interface())
	}
}

func PrintPreferenceToView(v *gocui.View) {
	fmt.Fprintln(v, color.BlueString(I118Prt.Sprintf("current_preference", "")))

	val := reflect.ValueOf(Conf)
	typeOfS := val.Type()
	for i := 0; i < reflect.ValueOf(Conf).NumField(); i++ {
		val := val.Field(i)
		fmt.Fprintln(v, fmt.Sprintf("  %s: %v", typeOfS.Field(i).Name, val.Interface()))
	}
}

func convertWorkDir(path string) string {
	if path == "./" || path == "." {
		path, _ = filepath.Abs(`.`)
		path = path + string(os.PathSeparator)
	} else {
		if strings.LastIndex(path, "/") != len(path)-1 {
			path = path + string(os.PathSeparator)
		}
	}

	return path
}
