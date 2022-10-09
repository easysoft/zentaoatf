package script

import (
	"github.com/easysoft/zentaoatf/test/ui/conf"
	plw "github.com/easysoft/zentaoatf/test/ui/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	"testing"
)

func Demo(t provider.T) {
	webpage, _ := plw.OpenUrl("https://baidu.com", t)
	defer func() { plw.Close(t) }()

	webpage.GetLocator("//*[@id=\"kw1\"]", t).Click(t)
}

func TestDemo(t *testing.T) {
	conf.ExitAllOnError = true

	runner.Run(t, "设置界面语言", Demo)
}
