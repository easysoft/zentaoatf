package main

import (
	"testing"

	plwConf "github.com/easysoft/zentaoatf/test/ui/conf"
	plwHelper "github.com/easysoft/zentaoatf/test/ui/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	playwright "github.com/playwright-community/playwright-go"
)

func CreateInterpreter(t provider.T) {
	t.ID("5744")
	t.AddParentSuite("管理解析器")
	webpage, _ := plwHelper.OpenUrl("http://127.0.0.1:8000/", t)
	defer webpage.Close()
	webpage.Click("#navbar>>[title=\"设置\"]")
	webpage.Click("#proxyTable>>tr:has-text('本地节点')>>button:has-text('运行环境')")
	webpage.Click("text=新建运行环境")
	locator := webpage.Locator("#interpreterFormModal select")
	locator.Click()
	locator.SelectNth(0, playwright.SelectOptionValues{Values: &[]string{"python"}})
	webpage.WaitForTimeout(200)
	locator.SelectNth(1, playwright.SelectOptionValues{Indexes: &[]int{1}})
	webpage.Click("#interpreterFormModal>>text=确定")
	_ = webpage.Locator("#settingModal .z-tbody-td", playwright.PageLocatorOptions{HasText: "Python"})
}

func EditInterpreter(t provider.T) {
	t.ID("5745")
	t.AddParentSuite("管理解析器")
	webpage, _ := plwHelper.OpenUrl("http://127.0.0.1:8000/", t)
	defer webpage.Close()
	webpage.Click("#navbar>>[title=\"设置\"]")
	webpage.Click("#proxyTable>>tr:has-text('本地节点')>>button:has-text('运行环境')")
	locator := webpage.Locator("#settingModal .z-tbody-tr:has-text('Python')>>text=编辑")
	locator.Click()
	locator = webpage.Locator("#interpreterFormModal select")
	locator.SelectNth(0, playwright.SelectOptionValues{Values: &[]string{"python"}})
	webpage.WaitForTimeout(200)
	locator.SelectNth(1, playwright.SelectOptionValues{Indexes: &[]int{1}})
	webpage.Click("#interpreterFormModal>>text=确定")
	locator = webpage.Locator("#settingModal .z-tbody-td", playwright.PageLocatorOptions{HasText: "Python"})
}

func DeleteInterpreter(t provider.T) {
	t.ID("5465")
	t.AddParentSuite("管理解析器")
	webpage, _ := plwHelper.OpenUrl("http://127.0.0.1:8000/", t)
	defer webpage.Close()
	webpage.Click("#navbar>>[title=\"设置\"]")
	webpage.Click("#proxyTable>>tr:has-text('本地节点')>>button:has-text('运行环境')")
	locator := webpage.Locator("#settingModal .z-tbody-tr", playwright.PageLocatorOptions{HasText: "Python"})
	locator = locator.Locator("text=删除")
	locator.Click()
	webpage.Click(":nth-match(.modal-action > button, 1)")
	webpage.WaitForTimeout(1000)
	plwConf.DisableErr()
	locator = webpage.Locator("#settingModal .z-tbody-tr", playwright.PageLocatorOptions{HasText: "Python"})
	c := locator.Count()
	if c > 0 {
		t.Errorf("Delete interpreter fail")
		t.FailNow()
	}
	plwConf.EnableErr()
}

func TestUiInterpreter(t *testing.T) {
	runner.Run(t, "客户端-创建语言解析器", CreateInterpreter)
	runner.Run(t, "客户端-编辑语言解析器", EditInterpreter)
	runner.Run(t, "客户端-删除语言解析器", DeleteInterpreter)
}
