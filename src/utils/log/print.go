package logUtils

import (
	"encoding/json"
	"fmt"
	stringUtils "github.com/easysoft/zentaoatf/src/utils/string"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/fatih/color"
	"strings"
)

func PrintUsage() {
	usage :=
		`
 help                 查看使用帮助。
 set                  全局设置语言、禅道站点连接参数。
 co      checkout     导出禅道系统中的用例，已存在的将更新标题和步骤描述。可指定产品、套件、测试单编号。
 up      update       从禅道系统更新已存在的用例。可指定产品、套件、测试单编号。
 run                  执行测试用例。可指定要目录、套件、任务、结果或脚本的路径，多个参数之间用空格隔开。
 ci      commit       将执行结果提交到禅道系统中。可指定测试日志目录，会弹出命令行图形界面。
 bug                  将执行结果中的失败用例，作为缺陷提交到餐到系统。可指定测试日志目录和用例编号，弹出命令行图形界面。
 ls      list         查看测试用例列表。可指定目录或文件列表。
 view                 查看测试用例详情。可指定目录或文件列表。
`

	fmt.Println(color.CyanString("\nUsage: "))
	fmt.Fprintf(color.Output, "%s\n", usage)
}

func PrintToSide(msg string) {
	if !vari.RunFromCui {
		fmt.Println(msg)
		return
	}
	slideView, _ := vari.Cui.View("side")
	slideView.Clear()

	fmt.Fprintln(slideView, msg)
}

func PrintToMainNoScroll(msg string) {
	if !vari.RunFromCui {
		fmt.Println(msg)
		return
	}
	mainView, _ := vari.Cui.View("main")
	mainView.Clear()

	fmt.Fprintln(mainView, msg)
}

func PrintToCmd(msg string) {
	if !vari.RunFromCui {
		fmt.Println(msg)
		return
	}
	cmdView, _ := vari.Cui.View("cmd")
	_, _ = fmt.Fprintln(cmdView, msg)
}
func PrintStructToCmd(obj interface{}) {
	str := stringUtils.StructToStr(obj)
	PrintToCmd(str)
}

func ClearSide() {
	slideView, _ := vari.Cui.View("side")
	slideView.Clear()
}

func PrintUnicode(str []byte) {
	var a interface{}

	temp := strings.Replace(string(str), "\\\\", "\\", -1)

	err := json.Unmarshal([]byte(temp), &a)

	var msg string
	if err == nil {
		msg = fmt.Sprint(a)
	} else {
		msg = temp
	}

	if !vari.RunFromCui {
		fmt.Println(msg)
	} else {
		PrintToCmd(msg)
	}
}
