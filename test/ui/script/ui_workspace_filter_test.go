package main

import (
	"testing"

	ztfTestHelper "github.com/easysoft/zentaoatf/test/helper/ztf"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	playwright "github.com/playwright-community/playwright-go"
)

var workspaceFilterPage playwright.Page

func FilterDir(t provider.T) {
	t.ID("5494")
	t.AddParentSuite("管理禅道站点下工作目录")
	locator, err := workspaceFilterPage.Locator("#siteMenuToggle")
	if err != nil {
		t.Errorf("The siteMenuToggle is missing: %v", err)
		t.FailNow()
		return
	}
	err = locator.Click()
	if err != nil {
		t.Errorf("The Click is fail: %v", err)
		t.FailNow()
		return
	}
	_, err = workspaceFilterPage.WaitForSelector("#navbar .list-item")
	if err != nil {
		t.Errorf("Wait for workspace list nav fail: %v", err)
		t.FailNow()
		return
	}
	err = workspaceFilterPage.Click(".list-item-title>>text=单元测试站点")
	if err != nil {
		t.Errorf("The Click workspace nav fail: %v", err)
		t.FailNow()
		return
	}
	err = workspaceFilterPage.Click(`[title="筛选"]`)
	if err != nil {
		t.Errorf("The Click create workspace fail: %v", err)
		t.FailNow()
		return
	}
	_, err = workspaceFilterPage.WaitForSelector("#filterModal")
	if err != nil {
		t.Errorf("Wait filter modal fail: %v", err)
		t.FailNow()
		return
	}

	err = workspaceFilterPage.Click("#filterModal>>.list-item-title:has-text(\"单元测试工作目录\")")
	if err != nil {
		t.Errorf("The Click php filter fail: %v", err)
		t.FailNow()
		return
	}
	eleArr, err := workspaceFilterPage.QuerySelectorAll("#leftPane .tree .tree-node")
	if len(eleArr) < 1 {
		t.Errorf("Filter valid fail: %v", err)
		t.FailNow()
		return
	}
}
func FilterSuite(t provider.T) {
	t.ID("5495")
	t.AddParentSuite("管理禅道站点下工作目录")
	err := workspaceFilterPage.Click(`[title="筛选"]`)
	if err != nil {
		t.Errorf("The Click create workspace fail: %v", err)
		t.FailNow()
		return
	}
	_, err = workspaceFilterPage.WaitForSelector("#filterModal")
	if err != nil {
		t.Errorf("Wait filter modal fail: %v", err)
		t.FailNow()
		return
	}
	workspaceFilterPage.WaitForTimeout(1000)
	err = workspaceFilterPage.Click("#filterModal>>.tab-nav:has-text(\"按套件\")")
	if err != nil {
		t.Errorf("The Click by suite fail: %v", err)
		t.FailNow()
		return
	}
	workspaceFilterPage.WaitForSelector("#filterModal>>.list-item-title:has-text(\"test_suite\")")
	err = workspaceFilterPage.Click("#filterModal>>.list-item-title:has-text(\"test_suite\")")
	if err != nil {
		t.Errorf("The Click test_suite filter fail: %v", err)
		t.FailNow()
		return
	}
	workspaceFilterPage.WaitForTimeout(200)
	workspaceFilterPage.WaitForSelector(".toolbar:has-text(\"按套件\")")
	ztfTestHelper.ExpandWorspace(workspaceFilterPage)
	workspaceFilterPage.WaitForTimeout(200)
	scriptLocator, err := workspaceFilterPage.Locator(".tree-node>>text=1_string_match.php")
	c, err := scriptLocator.Count()
	if err != nil || c == 0 {
		t.Errorf("Filter suite fail: %v", err)
		t.FailNow()
		return
	}
}
func ByModule(t provider.T) {
	t.ID("5493")
	t.AddParentSuite("管理禅道站点下工作目录")
	err := workspaceFilterPage.Click("#displayByMenuToggle")
	if err != nil {
		t.Errorf("The Click byModule btn fail: %v", err)
		t.FailNow()
	}
	workspaceFilterPage.WaitForTimeout(1000)
	err = workspaceFilterPage.Click(".dropdown-menu>>.list-item-content:has-text(\"按模块\")")
	if err != nil {
		t.Errorf("The Click by module fail: %v", err)
		t.FailNow()
	}
	err = workspaceFilterPage.Click(".tree-node-title:has-text(\"module1\")")
	if err != nil {
		t.Errorf("The Click module1 dir fail: %v", err)
		t.FailNow()
	}
	// scriptLocator, err := workspaceFilterPage.Locator(".tree-node>>:has-text(\"check string matches pattern\")")
	// c, err := scriptLocator.Count()
	// if err != nil || c == 0 {
	// 	t.Errorf("Filter suite fail: %v", err)
	// 	t.FailNow()
	// }
}
func FilterTask(t provider.T) {
	t.ID("5496")
	t.AddParentSuite("管理禅道站点下工作目录")
	locator, err := workspaceFilterPage.Locator("#siteMenuToggle")
	if err != nil {
		t.Errorf("The siteMenuToggle is missing: %v", err)
		t.FailNow()
		return
	}
	err = locator.Click()
	if err != nil {
		t.Errorf("The Click is fail: %v", err)
		t.FailNow()
		return
	}
	_, err = workspaceFilterPage.WaitForSelector("#navbar .list-item")
	if err != nil {
		t.Errorf("Wait for workspace list nav fail: %v", err)
		t.FailNow()
		return
	}
	err = workspaceFilterPage.Click(".list-item-title>>text=单元测试站点")
	if err != nil {
		t.Errorf("The Click workspace nav fail: %v", err)
		t.FailNow()
		return
	}
	err = workspaceFilterPage.Click(`[title="筛选"]`)
	if err != nil {
		t.Errorf("The Click create workspace fail: %v", err)
		t.FailNow()
		return
	}
	_, err = workspaceFilterPage.WaitForSelector("#filterModal")
	if err != nil {
		t.Errorf("Wait filter modal fail: %v", err)
		t.FailNow()
		return
	}
	workspaceFilterPage.WaitForTimeout(1000)
	err = workspaceFilterPage.Click("#filterModal>>.tab-nav:has-text(\"按测试单\")")
	if err != nil {
		t.Errorf("The Click by suite fail: %v", err)
		t.FailNow()
		return
	}
	workspaceFilterPage.WaitForSelector("#filterModal>>.list-item-title:has-text(\"企业网站第一期测试任务\")")
	err = workspaceFilterPage.Click("#filterModal>>.list-item-title:has-text(\"企业网站第一期测试任务\")")
	workspaceFilterPage.WaitForTimeout(200)
	if err != nil {
		t.Errorf("The Click test_task filter fail: %v", err)
		t.FailNow()
		return
	}
	workspaceFilterPage.WaitForSelector(".toolbar:has-text(\"按测试单\")")
	ztfTestHelper.ExpandWorspace(workspaceFilterPage)
	scriptLocator, err := workspaceFilterPage.Locator(".tree-node>>text=1_string_match.php")
	c, err := scriptLocator.Count()
	if err != nil || c == 0 {
		t.Errorf("Filter task fail: %v", err)
		t.FailNow()
		return
	}
}

func TestUiWorkspaceFilter(t *testing.T) {
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
		return
	}
	headless := false
	var slowMo float64 = 100
	workspaceBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web workspaceBrowser: %v", err)
		t.FailNow()
		return
	}
	workspaceFilterPage, err = workspaceBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
		return
	}
	if _, err := workspaceFilterPage.Goto("http://127.0.0.1:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
		t.FailNow()
		return
	}
	ztfTestHelper.ExpandWorspace(workspaceFilterPage)
	runner.Run(t, "客户端-按目录过滤禅道用例脚本", FilterDir)
	runner.Run(t, "客户端-按套件过滤禅道用例脚本", FilterSuite)
	runner.Run(t, "客户端-按测试单过滤禅道用例脚本", FilterTask)
	runner.Run(t, "客户端-按模块展示禅道用例脚本", ByModule)
	if err = workspaceBrowser.Close(); err != nil {
		t.Errorf("The workspaceBrowser cannot be closed: %v", err)
		t.FailNow()
		return
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
		return
	}
}
