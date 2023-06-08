package main

import (
	"fmt"
	"testing"

	constTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/conf"
	plwConf "github.com/easysoft/zentaoatf/cmd/test/ui/conf"
	plwHelper "github.com/easysoft/zentaoatf/cmd/test/ui/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	playwright "github.com/playwright-community/playwright-go"
)

func CreateSite(t provider.T) {
	t.ID("5466")
	t.AddParentSuite("配置禅道站点")

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	locator := webpage.Locator("#siteMenuToggle")
	locator.Click()
	webpage.WaitForSelector("#navbar .list-item")
	webpage.Click("text=禅道站点管理")

	webpage.Click("text=新建站点")
	locator = webpage.Locator("#siteFormModal input")
	locator.FillNth(0, constTestHelper.SiteName)
	locator.FillNth(1, constTestHelper.ZentaoSiteUrl)
	locator.FillNth(2, constTestHelper.ZentaoUsername)
	locator.FillNth(3, constTestHelper.ZentaoPassword)
	webpage.Click("text=确定")

	webpage.WaitForSelector(fmt.Sprintf(".list-item-content span:has-text('%s')", constTestHelper.SiteName))
	locator = webpage.Locator(".list-item-content span", playwright.PageLocatorOptions{HasText: constTestHelper.SiteName})
}

func CreateSiteWithFullUrl(t provider.T) {
	t.ID("7575")
	t.AddParentSuite("禅道站点URL格式兼容")

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	locator := webpage.Locator("#siteMenuToggle")
	locator.Click()
	webpage.WaitForSelector("#navbar .list-item")
	webpage.Click("text=禅道站点管理")
	webpage.Click("text=新建站点")

	locator = webpage.Locator("#siteFormModal input")
	locator.FillNth(0, constTestHelper.SiteName)
	locator.FillNth(1, constTestHelper.ZentaoSiteUrl+"/my.php")
	locator.FillNth(2, constTestHelper.ZentaoUsername)
	locator.FillNth(3, constTestHelper.ZentaoPassword)
	webpage.Click("text=确定")

	webpage.WaitForSelector("fmt.Sprintf(".list-item-content span:has-text('%s')", constTestHelper.SiteName)")
	locator = webpage.Locator(".list-item-content span", playwright.PageLocatorOptions{HasText: constTestHelper.SiteName})
}

func EditSite(t provider.T) {
	t.ID("5466")
	t.AddParentSuite("配置禅道站点")

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	locator := webpage.Locator("#siteMenuToggle")
	locator.Click()
	webpage.WaitForSelector("#navbar .list-item")
	webpage.Click("text=禅道站点管理")

	plwConf.DisableErr()
	locator = webpage.Locator(".list-item", playwright.PageLocatorOptions{HasText: constTestHelper.SiteName})
	c := locator.Count()
	if c == 0 {
		CreateSite(t)
		EditSite(t)
		plwConf.EnableErr()
		return
	}
	plwConf.EnableErr()

	locator = webpage.Locator(".list-item", playwright.PageLocatorOptions{HasText: constTestHelper.SiteName})
	webpage.Click("text=编辑")

	locator = webpage.Locator("#siteFormModal input")
	locator.FillNth(0, fmt.Sprintf("%s-update", constTestHelper.SiteName))
	locator.FillNth(1, constTestHelper.ZentaoSiteUrl)
	locator.FillNth(2, constTestHelper.ZentaoUsername)
	locator.FillNth(3, constTestHelper.ZentaoPassword)
	webpage.Click("#siteFormModal>>.modal-action>>span:has-text(\"确定\")")

	webpage.WaitForSelector(fmt.Sprintf(".list-item-content span:has-text('%s-update')", constTestHelper.SiteName))
	locator = webpage.Locator(".list-item-content span", playwright.PageLocatorOptions{HasText: constTestHelper.SiteName+"-update"})
}
func DeleteSite(t provider.T) {
	t.ID("5466")
	t.AddParentSuite("配置禅道站点")

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	locator := webpage.Locator("#siteMenuToggle")
	locator.Click()
	webpage.WaitForSelector("#navbar .list-item")
	webpage.Click("text=禅道站点管理")

	locator = webpage.Locator(fmt.Sprintf(".list-item:has-text('%s')", constTestHelper.SiteName))
	webpage.Click("text=删除")
	webpage.WaitForTimeout(1000)

	webpage.Click(":nth-match(.modal-action > button, 1)")
	webpage.WaitForSelector("fmt.Sprintf(".list-item-content span:has-text('%s')", constTestHelper.SiteName)", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})

	plwConf.DisableErr()
	defer plwConf.EnableErr()
	locator = webpage.Locator(fmt.Sprintf(".list-item-content:has-text('%s')", constTestHelper.SiteName))
	c := locator.Count()
	if c > 0 {
		t.Errorf("Delete site fail")
		t.FailNow()
	}
}

func TestUiSite(t *testing.T) {
	runner.Run(t, "客户端-编辑禅道站点", EditSite)
	runner.Run(t, "客户端-删除禅道站点", DeleteSite)
	runner.Run(t, "客户端-创建禅道站点", CreateSite)
}
