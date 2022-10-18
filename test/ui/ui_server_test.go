package main

import (
	"testing"

	plwConf "github.com/easysoft/zentaoatf/test/ui/conf"
	plwHelper "github.com/easysoft/zentaoatf/test/ui/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	playwright "github.com/playwright-community/playwright-go"
)

func CreateServer(t provider.T) {
	t.ID("5737")
	t.AddParentSuite("管理服务器")
	webpage, _ := plwHelper.OpenUrl("http://127.0.0.1:8000/", t)
	defer webpage.Close()
	webpage.Click("#navbar>>[title=\"设置\"]")
	webpage.Click("text=新建服务器")
	locator := webpage.Locator("#serverFormModal input")
	locator.FillNth(0, "测试服务器")
	webpage.WaitForTimeout(200)
	locator.FillNth(1, "http://127.0.0.1:8085")

	webpage.Click("#serverFormModal>>text=确定")
	webpage.WaitForSelector("#settingModal .z-tbody-td:has-text('测试服务器')")
	locator = webpage.Locator("#settingModal .z-tbody-td", playwright.PageLocatorOptions{HasText: "测试服务器"})
}

func EditServer(t provider.T) {
	t.ID("5738")
	t.AddParentSuite("管理服务器")
	webpage, _ := plwHelper.OpenUrl("http://127.0.0.1:8000/", t)
	defer webpage.Close()
	webpage.Click("#navbar>>[title=\"设置\"]")
	webpage.WaitForSelector("#settingModal .z-tbody-tr:has-text('测试服务器')")
	locator := webpage.Locator("#settingModal .z-tbody-tr:has-text('测试服务器')>>td>>nth=-1>>text=编辑")
	locator.Click()
	locator = webpage.Locator("#serverFormModal input")
	locator.FillNth(0, "测试服务器-update")
	webpage.Click("#serverFormModal>>text=确定")
	webpage.WaitForSelector("#serverFormModal", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})
	webpage.WaitForSelector("#settingModal .z-tbody-td:has-text('测试服务器-update')")
	locator = webpage.Locator("#settingModal .z-tbody-td", playwright.PageLocatorOptions{HasText: "测试服务器-update"})
}

func DeleteServer(t provider.T) {
	t.ID("5739")
	t.AddParentSuite("管理服务器")
	webpage, _ := plwHelper.OpenUrl("http://127.0.0.1:8000/", t)
	defer webpage.Close()
	webpage.Click("#navbar>>[title=\"设置\"]")
	webpage.WaitForSelector("#settingModal .z-tbody-tr:has-text('测试服务器')")
	locator := webpage.Locator("#settingModal .z-tbody-tr:has-text('测试服务器-update')>>td>>nth=-1")
	locator = locator.Locator("text=删除")
	locator.Click()

	webpage.Click(":nth-match(.modal-action > button, 1)")
	webpage.WaitForSelector("#settingModal .z-tbody-td:has-text('测试服务器-update')", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})
	plwConf.DisableErr()
	defer plwConf.EnableErr()
	c := locator.Count()
	locator = webpage.Locator("#settingModal .z-tbody-tr:has-text('测试服务器-update')")
	if c > 0 {
		t.Errorf("Delete server fail")
		t.FailNow()
	}
}

func TestUiServer(t *testing.T) {
	runner.Run(t, "客户端-创建服务器", CreateServer)
	runner.Run(t, "客户端-编辑服务器", EditServer)
	runner.Run(t, "客户端-删除服务器", DeleteServer)
}
