package configUtils

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/utils/common"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/file"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/fatih/color"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"sync"
)

func InitPreference() {
	// preference from yaml
	vari.Prefer = getInst()

	// screen size
	InitScreenSize()

	// internationalization
	i118Utils.InitI118(vari.Prefer.Language)

	if strings.Index(os.Args[0], "atf") > -1 && (len(os.Args) > 1 && os.Args[1] != "set") {
		PrintCurrPreference()
	}
}

func SetLanguage(lang string, dumb bool) {
	buf, _ := ioutil.ReadFile(constant.PreferenceFile)
	yaml.Unmarshal(buf, &vari.Prefer)

	vari.Prefer.Language = lang

	data, _ := yaml.Marshal(&vari.Prefer)
	ioutil.WriteFile(constant.PreferenceFile, data, 0666)

	// re-init language resource
	i118Utils.InitI118(vari.Prefer.Language)

	if !dumb {
		logUtils.PrintToCmd(color.CyanString(i118Utils.I118Prt.Sprintf("set_preference", i118Utils.I118Prt.Sprintf("lang"),
			i118Utils.I118Prt.Sprintf(vari.Prefer.Language))))
	}
}

func SetWorkDir(dir string, dumb bool) {
	fileUtils.MkDirIfNeeded(dir)

	buf, _ := ioutil.ReadFile(constant.PreferenceFile)
	yaml.Unmarshal(buf, &vari.Prefer)

	dir = commonUtils.ConvertWorkDir(dir)

	vari.Prefer.WorkDir = dir
	UpdateWorkDirHistoryForSwitch()

	data, _ := yaml.Marshal(&vari.Prefer)
	ioutil.WriteFile(constant.PreferenceFile, data, 0666)

	if !dumb {
		logUtils.PrintToCmd(color.CyanString(i118Utils.I118Prt.Sprintf("set_preference", i118Utils.I118Prt.Sprintf("workDir"), vari.Prefer.WorkDir)))
	}
}

func getInst() model.Preference {
	var once sync.Once
	once.Do(func() {
		vari.Prefer = model.Preference{}
		if fileUtils.FileExist(constant.PreferenceFile) {
			buf, _ := ioutil.ReadFile(constant.PreferenceFile)
			yaml.Unmarshal(buf, &vari.Prefer)
		} else { // init
			vari.Prefer.Language = "en"
			vari.Prefer.WorkDir = commonUtils.ConvertWorkDir(".")

			history := model.WorkHistory{Id: uuid.NewV4().String(), ProjectPath: vari.Prefer.WorkDir}
			vari.Prefer.WorkHistories = []model.WorkHistory{history}

			data, _ := yaml.Marshal(&vari.Prefer)
			ioutil.WriteFile(constant.PreferenceFile, data, 0666)
		}
	})
	return vari.Prefer
}

func PrintCurrPreference() {
	color.Cyan(i118Utils.I118Prt.Sprintf("current_preference", ""))

	val := reflect.ValueOf(vari.Prefer)
	typeOfS := val.Type()
	for i := 0; i < reflect.ValueOf(vari.Prefer).NumField(); i++ {
		val := val.Field(i)
		name := typeOfS.Field(i).Name

		if !vari.RunFromCui && (name == "Width" || name == "Height" || name == "WorkHistories") {
			continue
		}
		fmt.Printf("  %s: %v \n", typeOfS.Field(i).Name, val.Interface())
	}
}

func PrintPreferenceToView() {
	cmdView, _ := vari.Cui.View("cmd")
	fmt.Fprintln(cmdView, color.CyanString(i118Utils.I118Prt.Sprintf("current_preference", "")))

	val := reflect.ValueOf(vari.Prefer)
	typeOfS := val.Type()
	for i := 0; i < reflect.ValueOf(vari.Prefer).NumField(); i++ {
		val := val.Field(i)
		fmt.Fprintln(cmdView, fmt.Sprintf("  %s: %v", typeOfS.Field(i).Name, val.Interface()))
	}
}

func UpdateWorkDirHistoryForGenerate() { // update the first one
	conf := ReadCurrConfig()

	vari.Prefer.WorkHistories[0].ProjectName = conf.ProjectName
	vari.Prefer.WorkHistories[0].EntityType = conf.EntityType
	vari.Prefer.WorkHistories[0].EntityVal = conf.EntityVal

	data, _ := yaml.Marshal(&vari.Prefer)
	ioutil.WriteFile(constant.PreferenceFile, data, 0666)
}

func UpdateWorkDirHistoryForSwitch() {
	histories := vari.Prefer.WorkHistories

	// 已经是第一个，不做操作
	if histories[0].ProjectPath == vari.Prefer.WorkDir {
		return
	}

	// 移除元素
	idx := -1
	for i, item := range histories {
		if item.ProjectPath == vari.Prefer.WorkDir {
			idx = i
		}
	}
	if idx > -1 {
		histories = append(histories[:idx], histories[idx+1:]...)
	}

	// 头部插入元素
	conf := ReadCurrConfig()

	history := model.WorkHistory{Id: uuid.NewV4().String(), ProjectName: conf.ProjectName, ProjectPath: vari.Prefer.WorkDir,
		EntityType: conf.EntityType, EntityVal: conf.EntityVal}

	histories = append([]model.WorkHistory{history}, histories...)

	// 只保存最后10个
	if len(histories) > 10 {
		histories = histories[:10]
	}

	vari.Prefer.WorkHistories = histories
}
