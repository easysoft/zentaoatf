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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	languageRe    = regexp.MustCompile("Enter the language you want to use|请输入你期望使用的语言")
	configRe      = regexp.MustCompile("Do you want to config Zentao site|现在开始配置禅道系统的同步参数")
	urlRe         = regexp.MustCompile("Zentao site url|请输入禅道站点网址")
	accountRe     = regexp.MustCompile("Zentao account|请输入登录账号")
	passwordRe    = regexp.MustCompile("Zentao password|请输入账号密码")
	interpreterRe = regexp.MustCompile("Do you want to config script interpreter|现在配置脚本的解释程序")
	successRe     = regexp.MustCompile("Success")
	langsMap      = map[string]string{
		"php":        "D:\\Program Files\\phpstudy_pro\\Extensions\\php\\php7.4.3nts\\php.exe",
		"javascript": "",
		"lua":        "",
		"perl":       "",
		"python":     "",
		"ruby":       "",
		"tcl":        "",
		// "go":         "",
	}
	langs   = []string{"php", "javascript", "lua", "perl", "python", "ruby", "tcl"}
	newline = "\n"
)

type SetSuit struct {
	suite.Suite
	testCount uint32
}

func (s *SetSuit) TestChSetSuite() {
	assert.Equal(s.Suite.T(), "Success", testSet("2"))
}

func (s *SetSuit) TestEnSetSuite() {
	assert.Equal(s.Suite.T(), "Success", testSet("1"))
}

func testSet(language string) (ret string) {
	cmd := "ztf set"
	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		return err.Error()
	}
	defer child.Close()
	if _, err = child.Expect(languageRe, 3*time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", languageRe, err.Error())
	}

	if err = child.Send(language + newline); err != nil {
		return err.Error()
	}
	if _, err := child.Expect(configRe, time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", configRe, err.Error())
	}
	if err = child.Send("y" + newline); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(urlRe, time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", urlRe, err.Error())
	}
	if err = child.Send("http://127.0.0.1:81/zentao/" + newline); err != nil {
		return err.Error()
	}
	if _, err = child.Expect(accountRe, time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", accountRe, err.Error())
	}
	if err = child.Send("admin" + newline); err != nil {
		return err.Error()
	}

	if _, err = child.Expect(passwordRe, time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", passwordRe, err.Error())
	}
	if err = child.Send("123456." + newline); err != nil {
		return err.Error()
	}
	if runtime.GOOS == "windows" {
		if _, err = child.Expect(interpreterRe, time.Second); err != nil {
			return fmt.Sprintf("expect %s, actual %s", interpreterRe, err.Error())
		}
		if err = child.Send("y" + newline); err != nil {
			return err.Error()
		}
		for _, lang := range langs {
			if _, err = child.Expect(regexp.MustCompile(lang), time.Second); err != nil {
				return fmt.Sprintf("expect %s, actual %s", lang, err.Error())
			}
			if err = child.Send(langsMap[lang] + newline); err != nil {
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

func TestSet(t *testing.T) {
	if runtime.GOOS == "windows" {
		newline = "\r\n"
	}
	suite.Run(t, new(SetSuit))
}
