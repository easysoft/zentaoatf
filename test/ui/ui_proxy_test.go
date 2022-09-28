package main

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	playwright "github.com/playwright-community/playwright-go"
)

var proxyBrowser playwright.Browser

func CreateProxy(t provider.T) {
	t.ID("5465")
	t.AddParentSuite("设置界面语言")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	if proxyBrowser == nil || !proxyBrowser.IsConnected() {
		proxyBrowser, err = pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	}
	if err != nil {
		t.Errorf("Fail to launch the web proxyBrowser: %v", err)
		t.FailNow()
	}
	page, err := proxyBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
	defer func() {
		if err = proxyBrowser.Close(); err != nil {
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

	err = page.Click("#navbar>>[title=\"设置\"]")
	if err != nil {
		t.Errorf("The Click proxy nav fail: %v", err)
		t.FailNow()
	}
	err = page.Click("#serverTable>>button:has-text('新建执行节点')")
	if err != nil {
		t.Errorf("Click open interpreter modal fail: %v", err)
		t.FailNow()
	}
	locator, err := page.Locator("#proxyFormModal input")
	if err != nil {
		t.Errorf("Find create proxy input fail: %v", err)
		t.FailNow()
	}
	nameInput, err := locator.Nth(0)
	if err != nil {
		t.Errorf("Find lang select fail: %v", err)
		t.FailNow()
	}
	err = nameInput.Fill("测试执行节点")
	if err != nil {
		t.Errorf("Fill name fail: %v", err)
		t.FailNow()
	}
	page.WaitForTimeout(200)
	pathSelect, err := locator.Nth(1)
	if err != nil {
		t.Errorf("Find path input fail: %v", err)
		t.FailNow()
	}
	err = pathSelect.Fill("http://127.0.0.1:8085")
	if err != nil {
		t.Errorf("Fil path input fail: %v", err)
		t.FailNow()
	}

	err = page.Click("#proxyFormModal>>text=确定")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		t.FailNow()
	}
	page.WaitForSelector("#proxyFormModal", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})
	page.WaitForTimeout(1000)
	locator, err = page.Locator("#proxyTable .z-tbody-td >> :scope:has-text('测试执行节点')")
	c, err := locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find created proxy fail: %v", err)
		t.FailNow()
	}
}
func EditProxy(t provider.T) {
	t.ID("5465")
	t.AddParentSuite("设置界面语言")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	if proxyBrowser == nil || !proxyBrowser.IsConnected() {
		proxyBrowser, err = pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	}
	defer proxyBrowser.Close()
	defer pw.Stop()
	if err != nil {
		t.Errorf("Fail to launch the web proxyBrowser: %v", err)
		t.FailNow()
	}
	page, err := proxyBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
	defer func() {
		if err = proxyBrowser.Close(); err != nil {
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

	err = page.Click("#navbar>>[title=\"设置\"]")
	if err != nil {
		t.Errorf("The Click proxy nav fail: %v", err)
		t.FailNow()
	}
	locator, err := page.Locator("#proxyTable", playwright.PageLocatorOptions{HasText: "测试执行节点"})
	if err != nil {
		t.Errorf("Find 测试执行节点 tr fail: %v", err)
		t.FailNow()
	}
	locator, err = locator.Locator("text=编辑")
	if err != nil {
		t.Errorf("Find 测试执行节点 edit btn fail: %v", err)
		t.FailNow()
	}
	err = locator.Click()
	if err != nil {
		t.Errorf("The Click update site fail: %v", err)
		t.FailNow()
	}
	locator, err = page.Locator("#proxyFormModal input")
	if err != nil {
		t.Errorf("Find create proxy input fail: %v", err)
		t.FailNow()
	}
	nameInput, err := locator.Nth(0)
	if err != nil {
		t.Errorf("Find name select fail: %v", err)
		t.FailNow()
	}
	err = nameInput.Fill("测试执行节点-update")
	if err != nil {
		t.Errorf("Fill name input fail: %v", err)
		t.FailNow()
	}
	err = page.Click("#proxyFormModal>>text=确定")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		t.FailNow()
	}
	page.WaitForSelector("#proxyFormModal", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})
	page.WaitForTimeout(1000)
	locator, err = page.Locator("#proxyTable .z-tbody-td >> :scope:has-text('测试执行节点')")
	c, err := locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find updated proxy fail: %v", err)
		t.FailNow()
	}
}
func DeleteProxy(t provider.T) {
	t.ID("5465")
	t.AddParentSuite("设置界面语言")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	if proxyBrowser == nil || !proxyBrowser.IsConnected() {
		proxyBrowser, err = pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
		if err != nil {
			t.Errorf("Fail to launch the web proxyBrowser: %v", err)
			t.FailNow()
		}
	}
	page, err := proxyBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
	defer func() {
		if err = proxyBrowser.Close(); err != nil {
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

	err = page.Click("#navbar>>[title=\"设置\"]")
	if err != nil {
		t.Errorf("The Click proxy nav fail: %v", err)
		t.FailNow()
	}
	locator, err := page.Locator("#proxyTable", playwright.PageLocatorOptions{HasText: "测试执行节点-update"})
	if err != nil {
		t.Errorf("Find python tr fail: %v", err)
		t.FailNow()
	}
	locator, err = locator.Locator("text=删除")
	if err != nil {
		t.Errorf("Find 测试执行节点 del btn fail: %v", err)
		t.FailNow()
	}
	err = locator.Click()
	if err != nil {
		t.Errorf("The Click update site fail: %v", err)
		t.FailNow()
	}

	err = page.Click(":nth-match(.modal-action > button, 1)")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		t.FailNow()
	}
	page.WaitForTimeout(1000)
	locator, err = page.Locator("#settingModal .z-tbody-tr", playwright.PageLocatorOptions{HasText: "测试执行节点-update"})
	c, err := locator.Count()
	if err != nil || c > 0 {
		t.Errorf("Delete proxy fail: %v", err)
		t.FailNow()
	}
}

func TestUiProxy(t *testing.T) {
	runner.Run(t, "客户端-创建解析器", CreateProxy)
	runner.Run(t, "客户端-编辑解析器", EditProxy)
	runner.Run(t, "客户端-删除解析器", DeleteProxy)
}
