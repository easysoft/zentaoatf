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
	"testing"
	"time"

	expect "github.com/easysoft/zentaoatf/pkg/lib/expect"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

var (
	languageRe    = regexp.MustCompile("Enter the language you want to use|请输入你期望使用的语言")
	configRe      = regexp.MustCompile("Do you want to config Zentao site|现在开始配置禅道系统的同步参数")
	urlRe         = regexp.MustCompile("Zentao site url|请输入禅道站点网址")
	accountRe     = regexp.MustCompile("Zentao account|请输入登录账号")
	passwordRe    = regexp.MustCompile("Zentao password|请输入账号密码")
	interpreterRe = regexp.MustCompile("Do you want to config script interpreter|现在配置脚本的解释程序")
	successRe     = regexp.MustCompile("Success")
	langMap       = map[string]string{
		"php":        "D:\\Program Files\\phpstudy_pro\\Extensions\\php\\php7.4.3nts\\php.exe",
		"javascript": "",
		"lua":        "",
		"perl":       "",
		"python":     "",
		"ruby":       "",
		"tcl":        "",
		// "go":         "",
	}
	langArray = []string{"php", "javascript", "lua", "perl", "python", "ruby", "tcl"}
	newline   = "\n"
)

type SetSuite struct {
	suite.Suite
}

func (s *SetSuite) BeforeEach(t provider.T) {
	t.ID("1579")
	t.AddSubSuite("命令行-set")
}

func (s *SetSuite) TestChSetSuite(t provider.T) {
	t.Require().Equal("Success", testSet("2"))
}

func (s *SetSuite) TestEnSetSuite(t provider.T) {
	t.Require().Equal("Success", testSet("1"))
}

func testSet(language string) (ret string) {
	cmd := "ztf set"
	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		return err.Error()
	}
	defer child.Close()
	if _, err = child.Expect(languageRe, 3*time.Second*5); err != nil {
		return fmt.Sprintf("expect %s, actual %s", languageRe, err.Error())
	}

	if err = child.Send(language + newline); err != nil {
		return err.Error()
	}
	if _, err := child.Expect(configRe, time.Second*5); err != nil {
		return fmt.Sprintf("expect %s, actual %s", configRe, err.Error())
	}
	if err = child.Send("y" + newline); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(urlRe, time.Second*5); err != nil {
		return fmt.Sprintf("expect %s, actual %s", urlRe, err.Error())
	}
	if err = child.Send("http://127.0.0.1:8081/" + newline); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(accountRe, time.Second*5); err != nil {
		return fmt.Sprintf("expect %s, actual %s", accountRe, err.Error())
	}
	if err = child.Send("admin" + newline); err != nil {
		return err.Error()
	}

	if _, err = child.Expect(passwordRe, time.Second*5); err != nil {
		return fmt.Sprintf("expect %s, actual %s", passwordRe, err.Error())
	}
	if err = child.Send("Test123456." + newline); err != nil {
		return err.Error()
	}
	if runtime.GOOS == "windows" {
		if _, err = child.Expect(interpreterRe, time.Second*5); err != nil {
			return fmt.Sprintf("expect %s, actual %s", interpreterRe, err)
		}
		if err = child.Send("y" + newline); err != nil {
			return err.Error()
		}
		for _, lang := range langArray {
			if _, err = child.Expect(regexp.MustCompile(lang), time.Second*5); err != nil {
				return fmt.Sprintf("expect %s, actual %s", lang, err.Error())
			}
			if err = child.Send(langMap[lang] + newline); err != nil {
				return err.Error()
			}
		}
	}
	if _, err = child.Expect(successRe, 5*time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", successRe, err.Error())
	}
	child.Close()
	return "Success"
}

func TestCliSet(t *testing.T) {
	if runtime.GOOS == "windows" {
		newline = "\r\n"
	}
	suite.RunSuite(t, new(SetSuite))
}
