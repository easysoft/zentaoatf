package main

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	playwright "github.com/playwright-community/playwright-go"
)

var bugBrowser playwright.Browser

func ScriptBug(t provider.T) {
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	bugBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web bugBrowser: %v", err)
		t.FailNow()
	}
	page, err := bugBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
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
	page.WaitForTimeout(200)
	err = page.Click(".statistic>>span>>nth=3")
	if err != nil {
		t.Errorf("Click bug btn fail: %v", err)
		t.FailNow()
	}
	elements, _ := page.QuerySelectorAll("#bugsModal>>tr")

	bugTimesInt := len(elements)
	if bugTimesInt < 2 {
		t.Error("statistic error")
		t.FailNow()
	}
	if err = bugBrowser.Close(); err != nil {
		t.Errorf("The bugBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}

func ScriptsBug(t provider.T) {
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	bugBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web bugBrowser: %v", err)
		t.FailNow()
	}
	page, err := bugBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
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
	err = page.Click(`[title="批量选择"]`)
	if err != nil {
		t.Errorf("The Click select btn fail: %v", err)
		t.FailNow()
	}
	scriptLocator, err := locator.Locator(".tree-node-item:has-text('1_string_match.php')")
	if err != nil {
		t.Errorf("Find 1_string_match.php fail: %v", err)
		t.FailNow()
	}
	scriptLocator, err = scriptLocator.Locator(".tree-node-check")
	if err != nil {
		t.Errorf("Find 1_string_match.php checkbox fail: %v", err)
		t.FailNow()
	}
	err = scriptLocator.Click()
	if err != nil {
		t.Errorf("Click 1_string_match.php checkbox fail: %v", err)
		t.FailNow()
	}
	scriptLocator, err = locator.Locator(".tree-node-item:has-text('2_webpage_extract.php')")
	if err != nil {
		t.Errorf("Find 2_webpage_extract.php fail: %v", err)
		t.FailNow()
	}
	scriptLocator, err = scriptLocator.Locator(".tree-node-check")
	if err != nil {
		t.Errorf("Find 2_webpage_extract.php checkbox fail: %v", err)
		t.FailNow()
	}
	err = scriptLocator.Click()
	if err != nil {
		t.Errorf("Click 2_webpage_extract.php checkbox fail: %v", err)
		t.FailNow()
	}
	err = page.Click(`[title="禅道BUG"]`)
	if err != nil {
		t.Errorf("The Click bug btn fail: %v", err)
		t.FailNow()
	}
	elements, _ := page.QuerySelectorAll("#bugsModal>>tr")

	bugTimesInt := len(elements)
	if bugTimesInt < 2 {
		t.Error("statistic error")
		t.FailNow()
	}
	if err = bugBrowser.Close(); err != nil {
		t.Errorf("The bugBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}
func TestUiBug(t *testing.T) {
	runner.Run(t, "客户端-查看单个脚本bug列表", ScriptBug)
	runner.Run(t, "客户端-查看选中脚本bug列表", ScriptsBug)
}
