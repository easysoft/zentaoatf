package constant

import (
	"fmt"
	"os"
)

var (
	ConfigVer  = 1
	ConfigFile = "conf.yaml"

	UrlZentaoSettings = "zentaoSettings"
	UrlImportProject  = "importProject"
	UrlSubmitResult   = "submitResults"
	UrlReportBug      = "reportBug"

	ExtNameSuite = "cs"
	ExtNameJson  = "json"
	ExtNameTxt   = "txt"

	LanguageDefault = "en"
	LanguageEN      = "en"
	LanguageZH      = "zh"

	EnRes = fmt.Sprintf("res%smessages_en.json", string(os.PathSeparator))
	ZhRes = fmt.Sprintf("res%smessages_zh.json", string(os.PathSeparator))

	ScriptDir = fmt.Sprintf("scripts%s", string(os.PathSeparator))
	LogDir    = fmt.Sprintf("logs%s", string(os.PathSeparator))

	LeftWidth = 36
	MinWidth  = 130
	MinHeight = 36

	CmdViewHeight = 10

	CuiRunOutputView = "panelFileContent"

	RequestTypePathInfo = "PATH_INFO"

	Usage = `
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

	Example = `
 $>atf.exe run scripts-demo/tc-01.bat                执行本项目自带的测试用，非windows系统使用tc-01.sh脚本。

 $>atf.exe co                                        交互式导出禅道测试用例，将提示用户输入导出类型和编号。
 $>atf.exe co -product 1 -language python            导出编号为1的产品测试用例，使用python语言，缩写-p -l。
 $>atf.exe co -p 1 -m 16 -l python                   导出产品编号为1、模块编号为16的测试用例。
 $>atf.exe co -s 1 -l python                         导出编号为1的套件所含测试用例。
 $>atf.exe co -t 1 -l python                         导出编号为1的测试单所含用例。

 $>atf.exe run dir1 dir2 tc01.py                     执行目录dir1和dir2目录下，以及tc01.py文件的用例。
 $>atf.exe run c:\scripts all.cs                     执行本机all.cs测试套件，脚本在c:\scripts中。
 $>atf.exe run c:\scripts -suite 1                   执行禅道中编号为1的套件, 缩写-s。
 $>atf.exe run c:\scripts -task 1                    执行禅道系统中编号为1的任务, 缩写-t
 $>atf.exe run c:\scripts c:\19-08-27\result.txt     执行指定结果文件中失败的用例，缩写-r。

 $>atf.exe ci tc01.py                                将脚本里面修改的用例信息，同步到禅道系统。
 $>atf.exe cr logs/2019-08-28T164819                 提交测试结果到禅道系统。
 $>atf.exe cb logs/2019-08-28T164819 -case 1         将编号为1的用例结果提交为缺陷，缩写-c。

 $>atf.exe list dir1 .                               列出目录dir1，以及当前目录下的所有脚本文件。
 $>atf.exe view tc01.py tc02.py                      查看指定路径的测试脚本。
`
)
