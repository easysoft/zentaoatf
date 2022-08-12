package main

/**

$>ztf set   根据系统提示，设置语言、禅道地址、账号等，Windows下会提示输入语言解释程序。

cid=0
pid=0

中文set >> Success
英文set >> Success

*/
import (
	"fmt"
	"regexp"
	"runtime"
	"time"

	expect "github.com/google/goexpect"
)

var (
	languageRE    = regexp.MustCompile("Enter the language you want to use|请输入你期望使用的语言")
	configRE      = regexp.MustCompile("Do you want to config Zentao site|现在开始配置禅道系统的同步参数")
	urlRE         = regexp.MustCompile("Zentao site url|请输入禅道站点网址")
	accountRE     = regexp.MustCompile("Zentao account|请输入登录账号")
	passwordRE    = regexp.MustCompile("Zentao password|请输入账号密码")
	interpreterRE = regexp.MustCompile("Do you want to config script interpreter|现在配置脚本的解释程序")
	successRE     = regexp.MustCompile("Success")
	langs         = map[string]string{
		"php":        "D:\\Program Files\\phpstudy_pro\\Extensions\\php\\php7.4.3nts\\php.exe",
		"javascript": "",
		"lua":        "",
		"perl":       "",
		"python":     "",
		"ruby":       "",
		"tcl":        "",
		"go":         "",
	}
	newline = "\n"
)

func testSet(language string) {
	cmd := "ztf set"
	child, _, err := expect.Spawn(cmd, -1)
	if err != nil {
		fmt.Println(err)
	}
	defer child.Close()
	if _, _, err = child.Expect(languageRE, time.Second); err != nil {
		fmt.Printf("%s: %s%s", languageRE, err, newline)
		return
	}

	if err = child.Send(language + newline); err != nil {
		fmt.Println(err)
		return
	}
	if _, _, err := child.Expect(configRE, time.Second); err != nil {
		fmt.Printf("%s: %s%s", configRE, err, newline)
		return
	}
	if err = child.Send("y" + newline); err != nil {
		fmt.Println(err)
		return
	}
	if _, _, err = child.Expect(urlRE, time.Second); err != nil {
		fmt.Printf("%s: %s%s", urlRE, err, newline)
		return
	}
	if err = child.Send("http://pms.test/" + newline); err != nil {
		fmt.Println(err)
		return
	}
	if _, _, err = child.Expect(accountRE, time.Second); err != nil {
		fmt.Printf("%s: %s%s", accountRE, err, newline)
		return
	}
	if err = child.Send("admin" + newline); err != nil {
		fmt.Println(err)
		return
	}

	if _, _, err = child.Expect(passwordRE, time.Second); err != nil {
		fmt.Printf("%s: %s%s", passwordRE, err, newline)
		return
	}
	if err = child.Send("123456." + newline); err != nil {
		fmt.Println(err)
		return
	}
	if runtime.GOOS == "windows" {
		if _, _, err = child.Expect(interpreterRE, time.Second); err != nil {
			fmt.Printf("%s: %s%s", interpreterRE, err, newline)
			return
		}
		if err = child.Send("y" + newline); err != nil {
			fmt.Println(err)
			return
		}
		for lang, interpreter := range langs {
			if _, _, err = child.Expect(regexp.MustCompile(lang), time.Second); err != nil {
				fmt.Printf("%s: %s%s", lang, err, newline)
				return
			}
			if err = child.Send(interpreter + newline); err != nil {
				fmt.Println(err)
				return
			}
		}
	}
	if _, _, err = child.Expect(successRE, 5*time.Second); err != nil {
		fmt.Printf("%s: %s%s", successRE, err, newline)
		return
	}

	fmt.Println("Success")
}

func main() {
	if runtime.GOOS == "windows" {
		newline = "\r\n"
	}
	testSet("2")
	testSet("1")
}
