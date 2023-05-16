package main

import (
	"testing"

	ztfTestHelper "github.com/easysoft/zentaoatf/test/helper/ztf"
	plwConf "github.com/easysoft/zentaoatf/test/ui/conf"
	plwHelper "github.com/easysoft/zentaoatf/test/ui/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	playwright "github.com/playwright-community/playwright-go"
	constTestHelper "github.com/easysoft/zentaoatf/test/helper/conf"
)

func CollapseLog(t provider.T) {
	t.ID("5502")
	t.AddParentSuite("脚本执行日志")
	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage)
	ztfTestHelper.ExpandWorspace(webpage)
	scriptLocator := webpage.Locator(".tree-node-title>>text=1_string_match.php")
	scriptLocator.Click()
	webpage.Click(".tabs-nav-toolbar>>[title=\"Run\"]")
	webpage.WaitForSelector("#log-list>>.msg-span>>:has-text('执行1个用例，耗时')")
	webpage.Click(".btn[title=\"展开所有\"]")
	locator := webpage.Locator("#log-list>>.show-detail>>:has-text('[期望]')")
	webpage.WaitForTimeout(100)
	webpage.Click(".btn[title=\"折叠所有\"]")
	plwConf.DisableErr()
	defer plwConf.EnableErr()
	locator = webpage.Locator("#log-list>>.show-detail>>:has-text('[期望]')")
	count := locator.Count()
	if count > 0 {
		t.Error("Find Collapsed log fail")
		t.FailNow()
	}
}

func FullScreenLog(t provider.T) {
	t.ID("5749")
	t.AddParentSuite("脚本执行日志")
	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()
	webpage.WaitForSelector(".tree-node")
	locator := webpage.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	locator.Click()
	scriptLocator := locator.Locator("text=1_string_match.php")
	scriptLocator.Click()
	webpage.Click(".tabs-nav-toolbar>>[title=\"Run\"]")
	webpage.WaitForSelector("#log-list>>.msg-span>>:has-text('执行1个用例，耗时')")
	webpage.Click(".btn[title=\"向上展开\"]")
	webpage.WaitForTimeout(100)
	isHidden := webpage.IsHidden("#tabsPane")
	if !isHidden {
		t.Errorf("Check Full Screen fail")
		t.FailNow()
	}
}

func TestUiLog(t *testing.T) {
	runner.Run(t, "客户端-展开折叠执行日志", CollapseLog)
	runner.Run(t, "客户端-最大化脚本执行日志", FullScreenLog)
}
