package utils

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/fatih/color"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"sync"
)

var Prefer model.Preference

func InitPreference() {
	// preference from yaml
	Prefer = getInst()

	// screen size
	InitScreenSize()

	// internationalization
	InitI118(Prefer.Language)

	if strings.Index(os.Args[0], "atf") > -1 && (len(os.Args) > 1 && os.Args[1] != "set") {
		PrintPreference()
	}
}

func SetPreference(param string, val string, dumb bool) {
	buf, _ := ioutil.ReadFile(PreferenceFile)
	yaml.Unmarshal(buf, &Prefer)

	if param == "lang" {
		Prefer.Language = val
		if !dumb {
			color.Cyan(I118Prt.Sprintf("set_preference", I118Prt.Sprintf("lang"), I118Prt.Sprintf(Prefer.Language)))
		}
	} else if param == "workDir" {
		val = ConvertWorkDir(val)

		Prefer.WorkDir = val
		updateWorkDirHistory()
		if !dumb {
			color.Cyan(I118Prt.Sprintf("set_preference", I118Prt.Sprintf("workDir"), Prefer.WorkDir))
		}
	}
	data, _ := yaml.Marshal(&Prefer)
	ioutil.WriteFile(PreferenceFile, data, 0666)
}

func getInst() model.Preference {
	var once sync.Once
	once.Do(func() {
		Prefer = model.Preference{}
		if FileExist(PreferenceFile) {
			buf, _ := ioutil.ReadFile(PreferenceFile)
			yaml.Unmarshal(buf, &Prefer)
		} else { // init
			Prefer.Language = "en"
			Prefer.WorkDir = ConvertWorkDir(".")

			history := model.WorkHistory{Id: uuid.NewV4().String(), ProjectPath: Prefer.WorkDir}
			Prefer.WorkHistories = []model.WorkHistory{history}

			data, _ := yaml.Marshal(&Prefer)
			ioutil.WriteFile(PreferenceFile, data, 0666)
		}
	})
	return Prefer
}

func PrintPreference() {
	color.Cyan(I118Prt.Sprintf("current_preference", ""))

	val := reflect.ValueOf(Prefer)
	typeOfS := val.Type()
	for i := 0; i < reflect.ValueOf(Prefer).NumField(); i++ {
		val := val.Field(i)
		name := typeOfS.Field(i).Name

		if !RunFromCui && (name == "Width" || name == "Height" || name == "WorkHistories") {
			continue
		}
		fmt.Printf("  %s: %v \n", typeOfS.Field(i).Name, val.Interface())
	}
}

func PrintPreferenceToView() {
	cmdView, _ := Cui.View("cmd")
	fmt.Fprintln(cmdView, color.CyanString(I118Prt.Sprintf("current_preference", "")))

	val := reflect.ValueOf(Prefer)
	typeOfS := val.Type()
	for i := 0; i < reflect.ValueOf(Prefer).NumField(); i++ {
		val := val.Field(i)
		fmt.Fprintln(cmdView, fmt.Sprintf("  %s: %v", typeOfS.Field(i).Name, val.Interface()))
	}
}

func updateWorkDirHistory() {
	histories := Prefer.WorkHistories

	// 已经是第一个，不做操作
	if histories[0].ProjectPath == Prefer.WorkDir {
		return
	}

	// 移除元素
	idx := -1
	for i, item := range histories {
		if item.ProjectPath == Prefer.WorkDir {
			idx = i
		}
	}
	if idx > -1 {
		histories = append(histories[:idx], histories[idx+1:]...)
	}

	// 头部插入元素
	config := ReadCurrConfig()

	history := model.WorkHistory{Id: uuid.NewV4().String(), ProjectName: config.ProjectName, ProjectPath: Prefer.WorkDir,
		EntityType: config.EntityType, EntityVal: config.EntityVal}

	histories = append([]model.WorkHistory{history}, histories...)

	// 保存最后10个
	if len(histories) > 10 {
		histories = histories[:10]
	}

	Prefer.WorkHistories = histories
}
