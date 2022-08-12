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
	languageRe    = regexp.MustCompile("Enter the language you want to use|请输入你期望使用的语言")
	configRe      = regexp.MustCompile("Do you want to config Zentao site|现在开始配置禅道系统的同步参数")
	urlRe         = regexp.MustCompile("Zentao site url|请输入禅道站点网址")
	accountRe     = regexp.MustCompile("Zentao account|请输入登录账号")
	passwordRe    = regexp.MustCompile("Zentao password|请输入账号密码")
	interpreterRe = regexp.MustCompile("Do you want to config script interpreter|现在配置脚本的解释程序")
	successRe     = regexp.MustCompile("Success")
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
	if _, _, err = child.Expect(languageRe, time.Second); err != nil {
		fmt.Printf("%s: %s%s", languageRe, err, newline)
		return
	}

	if err = child.Send(language + newline); err != nil {
		fmt.Println(err)
		return
	}
	if _, _, err := child.Expect(configRe, time.Second); err != nil {
		fmt.Printf("%s: %s%s", configRe, err, newline)
		return
	}
	if err = child.Send("y" + newline); err != nil {
		fmt.Println(err)
		return
	}
	if _, _, err = child.Expect(urlRe, time.Second); err != nil {
		fmt.Printf("%s: %s%s", urlRe, err, newline)
		return
	}
	if err = child.Send("http://pms.test/" + newline); err != nil {
		fmt.Println(err)
		return
	}
	if _, _, err = child.Expect(accountRe, time.Second); err != nil {
		fmt.Printf("%s: %s%s", accountRe, err, newline)
		return
	}
	if err = child.Send("admin" + newline); err != nil {
		fmt.Println(err)
		return
	}

	if _, _, err = child.Expect(passwordRe, time.Second); err != nil {
		fmt.Printf("%s: %s%s", passwordRe, err, newline)
		return
	}
	if err = child.Send("123456." + newline); err != nil {
		fmt.Println(err)
		return
	}
	if runtime.GOOS == "windows" {
		if _, _, err = child.Expect(interpreterRe, time.Second); err != nil {
			fmt.Printf("%s: %s%s", interpreterRe, err, newline)
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
	if _, _, err = child.Expect(successRe, 5*time.Second); err != nil {
		fmt.Printf("%s: %s%s", successRe, err, newline)
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
