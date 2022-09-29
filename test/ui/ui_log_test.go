package main

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	playwright "github.com/playwright-community/playwright-go"
)

var logBrowser playwright.Browser

func CollapseLog(t provider.T) {
	t.ID("5502")
	t.AddParentSuite("脚本执行日志")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	logBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web logBrowser: %v", err)
		t.FailNow()
	}
	page, err := logBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
	defer func() {
		if err = logBrowser.Close(); err != nil {
			t.Errorf("The workspaceBrowser cannot be closed: %v", err)
			t.FailNow()
			return
		}
		if err = pw.Stop(); err != nil {
			t.Errorf("The playwright cannot be stopped: %v", err)
			t.FailNow()
			return
		}
	}()
	if _, err = page.Goto("http://127.0.0.1:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector(".tree-node")
	if err != nil {
		t.Errorf("Wait tree-node fail: %v", err)
		t.FailNow()
	}
	locator, err := page.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	c, err := locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find workspace fail: %v", err)
		t.FailNow()
	}
	err = locator.Click()
	if err != nil {
		t.Errorf("Click node fail: %v", err)
		t.FailNow()
	}
	scriptLocator, err := locator.Locator("text=1_string_match.php")
	if err != nil {
		t.Errorf("Find 1_string_match.php fail: %v", err)
		t.FailNow()
	}
	err = scriptLocator.Click()
	if err != nil {
		t.Errorf("Click script fail: %v", err)
		t.FailNow()
	}
	err = page.Click(".tabs-nav-toolbar>>[title=\"Run\"]")
	if err != nil {
		t.Errorf("Click run fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector("#log-list>>.msg-span>>:has-text('执行1个用例，耗时')")
	if err != nil {
		t.Errorf("Wait exec result fail: %v", err)
		t.FailNow()
	}
	err = page.Click(".btn[title=\"展开所有\"]")
	if err != nil {
		t.Errorf("Click expand btn fail: %v", err)
		t.FailNow()
	}
	locator, _ = page.Locator("#log-list>>.show-detail>>:has-text('[Expect]')")
	count, _ := locator.Count()
	if count == 0 {
		t.Error("Find expanded log fail")
		t.FailNow()
	}
	page.WaitForTimeout(100)
	err = page.Click(".btn[title=\"折叠所有\"]")
	if err != nil {
		t.Errorf("Click Collapse btn fail: %v", err)
		t.FailNow()
	}
	locator, _ = page.Locator("#log-list>>.show-detail>>:has-text('[Expect]')")
	count, _ = locator.Count()
	if count > 0 {
		t.Error("Find Collapsed log fail")
		t.FailNow()
	}
}

func FullScreenLog(t provider.T) {
	t.ID("5749")
	t.AddParentSuite("脚本执行日志")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	logBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web logBrowser: %v", err)
		t.FailNow()
	}
	page, err := logBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
	defer func() {
		if err = logBrowser.Close(); err != nil {
			t.Errorf("The workspaceBrowser cannot be closed: %v", err)
			t.FailNow()
			return
		}
		if err = pw.Stop(); err != nil {
			t.Errorf("The playwright cannot be stopped: %v", err)
			t.FailNow()
			return
		}
	}()
	if _, err = page.Goto("http://127.0.0.1:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector(".tree-node")
	if err != nil {
		t.Errorf("Wait tree-node fail: %v", err)
		t.FailNow()
	}
	locator, err := page.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	c, err := locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find workspace fail: %v", err)
		t.FailNow()
	}
	err = locator.Click()
	if err != nil {
		t.Errorf("Click node fail: %v", err)
		t.FailNow()
	}
	scriptLocator, err := locator.Locator("text=1_string_match.php")
	if err != nil {
		t.Errorf("Find 1_string_match.php fail: %v", err)
		t.FailNow()
	}
	err = scriptLocator.Click()
	if err != nil {
		t.Errorf("Click script fail: %v", err)
		t.FailNow()
	}
	err = page.Click(".tabs-nav-toolbar>>[title=\"Run\"]")
	if err != nil {
		t.Errorf("Click run fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector("#log-list>>.msg-span>>:has-text('执行1个用例，耗时')")
	if err != nil {
		t.Errorf("Wait exec result fail: %v", err)
		t.FailNow()
	}
	err = page.Click(".btn[title=\"向上展开\"]")
	if err != nil {
		t.Errorf("Click Collapse btn fail: %v", err)
		t.FailNow()
	}
	page.WaitForTimeout(100)
	isHidden, err := page.IsHidden("#tabsPane")
	if !isHidden {
		t.Errorf("Full Screen fail: %v", err)
		t.FailNow()
	}
}

func TestUiLog(t *testing.T) {
	runner.Run(t, "客户端-展开折叠执行日志", CollapseLog)
	runner.Run(t, "客户端-最大化脚本执行日志", FullScreenLog)
}
