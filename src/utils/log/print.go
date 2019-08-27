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
 help                  查看使用帮助。
 set                   全局设置语言、禅道站点连接参数。
 co      checkout      导出禅道系统中的用例，已存在的将更新标题和步骤描述。可指定产品、套件、测试单编号。
 up      update        从禅道系统更新已存在的用例。可指定产品、模块、套件、测试单编号。
 run                   执行测试用例。可指定目录、套件、脚本、测试结果文件的路径，也支持套件和任务的编号，多个参数之间用空格隔开。
 ci      commit        将执行结果提交到禅道系统中。可指定测试日志目录，会弹出命令行图形界面。
 bug                   将执行结果中的失败用例，作为缺陷提交到餐到系统。可指定测试日志目录和用例编号，弹出命令行图形界面。
 ls      list          查看测试用例列表。可指定目录和文件的列表，之间用空格隔开。
 view                  查看测试用例详情。可指定目录和文件的列表，之间用空格隔开。
`

	example :=
		`
 $>atf.exe run -f scripts-demo/tc-01.bat             执行本项目自带的测试用，非windows系统使用tc-01.sh脚本

 $>atf.exe co                                        交互式导出测试用例，将提示用户输入导出类型和编号。
 $>atf.exe co -p 1 -l python                         从禅道导出编号为1的产品测试用例，使用python语言。
 $>atf.exe co -p 1 -m 16 -l python                   从禅道导出产品编号为1、模块编号为16的测试用例。
 $>atf.exe co -s 1 -l python                         从禅道导出编号为1的套件所含测试用例。
 $>atf.exe co -t 1 -l python                         从禅道导出编号为1的测试单所含用例。

 $>atf.exe run dir1 dir2 tc01.py                     执行目录dir1和dir2目录下，以及tc01.py文件的用例
 $>atf.exe run c:\scripts -s 1                       执行禅道系统中编号为1的套件所含用例，脚本在目录c:\scripts中
 $>atf.exe run c:\scripts -s all.cs                  执行本机all.cs套件下的所有用例
 $>atf.exe run c:\scripts -t 1                       执行禅道系统中编号为1的任务所含用例
 $>atf.exe run c:\scripts -r c:\19-08-27\result.txt  执行指定结果文件中失败的用例，脚本在目录c:\scripts中
 $>atf.exe list dir1 .                               列出目录dir1，以及当前目录下的所有脚本文件
 $>atf.exe view tc01.py tc02.py                      查看指定路径的测试脚本     

`

	fmt.Println(color.CyanString("\nUsage: "))
	fmt.Fprintf(color.Output, "%s\n", usage)

	fmt.Println(color.CyanString("Example: "))
	fmt.Fprintf(color.Output, "%s\n", example)
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
