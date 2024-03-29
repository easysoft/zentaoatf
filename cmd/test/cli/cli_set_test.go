package main

/**
$>ztf set   根据系统提示，设置语言、禅道地址、账号等，Windows下会提示输入语言解释程序。
*/
import (
	"fmt"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"

	commonTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/common"
	constTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/conf"
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

	langMap = map[string]string{
		"php":        "D:\\Program Files\\phpstudy_pro\\Extensions\\php\\php7.4.3nts\\php.exe",
		"javascript": "",
		"lua":        "",
		"perl":       "",
		"python":     "",
		"ruby":       "",
		"tcl":        "",
		"go":         "",
	}
	langArray = []string{"go", "php", "ruby", "javascript", "lua", "perl", "python", "tcl"}
)

type SetSuite struct {
	suite.Suite
}

func (s *SetSuite) BeforeEach(t provider.T) {
	t.ID("1579")
	commonTestHelper.ReplaceLabel(t, "命令行-设置")
}

func (s *SetSuite) TestChSetSuite(t provider.T) {
	t.Require().Equal("Success", testSet("2"))
}

func (s *SetSuite) TestEnSetSuite(t provider.T) {
	t.Require().Equal("Success", testSet("1"))
}

func testSet(language string) (ret string) {
	cmd := commonTestHelper.GetZtfPath() + " set"
	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		return err.Error()
	}
	defer child.Close()
	if _, err = child.Expect(languageRe, 3*time.Second*5); err != nil {
		return fmt.Sprintf("expect %s, actual %s", languageRe, err.Error())
	}
	if err = child.Send(language + constTestHelper.NewLine); err != nil {
		return err.Error()
	}

	if _, err := child.Expect(configRe, time.Second*5); err != nil {
		return fmt.Sprintf("expect %s, actual %s", configRe, err.Error())
	}
	if err = child.Send("y" + constTestHelper.NewLine); err != nil {
		return err.Error()
	}

	if _, err = child.Expect(urlRe, time.Second*5); err != nil {
		return fmt.Sprintf("expect %s, actual %s", urlRe, err.Error())
	}
	if err = child.Send(constTestHelper.ZentaoSiteUrl + constTestHelper.NewLine); err != nil {
		return err.Error()
	}

	if _, err = child.Expect(accountRe, time.Second*5); err != nil {
		return fmt.Sprintf("expect %s, actual %s", accountRe, err.Error())
	}
	if err = child.Send(constTestHelper.ZentaoUsername + constTestHelper.NewLine); err != nil {
		return err.Error()
	}

	if _, err = child.Expect(passwordRe, time.Second*5); err != nil {
		return fmt.Sprintf("expect %s, actual %s", passwordRe, err.Error())
	}
	if err = child.Send(constTestHelper.ZentaoPassword + constTestHelper.NewLine); err != nil {
		return err.Error()
	}

	if runtime.GOOS == "windows" {
		if _, err = child.Expect(interpreterRe, time.Second*5); err != nil {
			return fmt.Sprintf("expect %s, actual %s", interpreterRe, err)
		}
		if err = child.Send("y" + constTestHelper.NewLine); err != nil {
			return err.Error()
		}

		for _, lang := range langArray {
			out, err := child.Expect(regexp.MustCompile("Please set script|请设置"), time.Second*5)
			if err != nil {
				return fmt.Sprintf("expect Please set script|请设置%s, actual %s", lang, err.Error())
			}

			sendMsg := ""
			if strings.Contains(out, "php") {
				sendMsg = "D:\\Program Files\\phpstudy_pro\\Extensions\\php\\php7.4.3nts\\php.exe"
			}

			if err = child.Send(sendMsg + constTestHelper.NewLine); err != nil {
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
	suite.RunSuite(t, new(SetSuite))
}
