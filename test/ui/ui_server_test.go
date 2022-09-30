package main

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	playwright "github.com/playwright-community/playwright-go"
)

var serverBrowser playwright.Browser

func CreateServer(t provider.T) {
	t.ID("5737")
	t.AddParentSuite("管理服务器")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	if serverBrowser == nil || !serverBrowser.IsConnected() {
		serverBrowser, err = pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	}
	if err != nil {
		t.Errorf("Fail to launch the web serverBrowser: %v", err)
		t.FailNow()
	}
	page, err := serverBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
	defer func() {
		if err = serverBrowser.Close(); err != nil {
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
		t.Errorf("The Click server nav fail: %v", err)
		t.FailNow()
	}
	err = page.Click("text=新建服务器")
	if err != nil {
		t.Errorf("The Click create server fail: %v", err)
		t.FailNow()
	}
	locator, err := page.Locator("#serverFormModal input")
	if err != nil {
		t.Errorf("Find create server input fail: %v", err)
		t.FailNow()
	}
	nameInput, err := locator.Nth(0)
	if err != nil {
		t.Errorf("Find lang select fail: %v", err)
		t.FailNow()
	}
	err = nameInput.Fill("测试服务器")
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

	err = page.Click("#serverFormModal>>text=确定")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector("#settingModal .z-tbody-td:has-text('测试服务器')")
	if err != nil {
		t.Errorf("Wait created server fail: %v", err)
		t.FailNow()
	}
	locator, err = page.Locator("#settingModal .z-tbody-td", playwright.PageLocatorOptions{HasText: "测试服务器"})
	c, err := locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find created server fail: %v", err)
		t.FailNow()
	}
}
func EditServer(t provider.T) {
	t.ID("5738")
	t.AddParentSuite("管理服务器")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	if serverBrowser == nil || !serverBrowser.IsConnected() {
		serverBrowser, err = pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	}
	defer serverBrowser.Close()
	defer pw.Stop()
	if err != nil {
		t.Errorf("Fail to launch the web serverBrowser: %v", err)
		t.FailNow()
	}
	page, err := serverBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
	defer func() {
		if err = serverBrowser.Close(); err != nil {
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
		t.Errorf("The Click server nav fail: %v", err)
		t.FailNow()
	}
	locator, err := page.Locator("#settingModal .z-tbody-tr:has-text('测试服务器')>>td>>nth=-1>>text=编辑")
	if err != nil {
		t.Errorf("Find 测试服务器 edit btn fail: %v", err)
		t.FailNow()
	}
	err = locator.Click()
	if err != nil {
		t.Errorf("The Click update site fail: %v", err)
		t.FailNow()
	}
	locator, err = page.Locator("#serverFormModal input")
	if err != nil {
		t.Errorf("Find create server input fail: %v", err)
		t.FailNow()
	}
	nameInput, err := locator.Nth(0)
	if err != nil {
		t.Errorf("Find name select fail: %v", err)
		t.FailNow()
	}
	err = nameInput.Fill("测试服务器-update")
	if err != nil {
		t.Errorf("Fill name input fail: %v", err)
		t.FailNow()
	}
	err = page.Click("#serverFormModal>>text=确定")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		t.FailNow()
	}
	page.WaitForSelector("#serverFormModal", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})
	_, err = page.WaitForSelector("#settingModal .z-tbody-td:has-text('测试服务器-update')")
	if err != nil {
		t.Errorf("Wait updated server fail: %v", err)
		t.FailNow()
	}
	locator, err = page.Locator("#settingModal .z-tbody-td", playwright.PageLocatorOptions{HasText: "测试服务器-update"})
	c, err := locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find updated server fail: %v", err)
		t.FailNow()
	}
}
func DeleteServer(t provider.T) {
	t.ID("5739")
	t.AddParentSuite("管理服务器")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	if serverBrowser == nil || !serverBrowser.IsConnected() {
		serverBrowser, err = pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
		if err != nil {
			t.Errorf("Fail to launch the web serverBrowser: %v", err)
			t.FailNow()
		}
	}
	page, err := serverBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
	defer func() {
		if err = serverBrowser.Close(); err != nil {
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
		t.Errorf("The Click server nav fail: %v", err)
		t.FailNow()
	}
	locator, err := page.Locator("#settingModal .z-tbody-tr:has-text('测试服务器-update')>>td>>nth=-1")
	if err != nil {
		t.Errorf("Find 测试服务器-update tr fail: %v", err)
		t.FailNow()
	}
	locator, err = locator.Locator("text=删除")
	if err != nil {
		t.Errorf("Find 测试服务器 del btn fail: %v", err)
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
	_, err = page.WaitForSelector("#settingModal .z-tbody-td:has-text('测试服务器-update')", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})
	if err != nil {
		t.Errorf("Wait updated server fail: %v", err)
		t.FailNow()
	}
	locator, err = page.Locator("#settingModal .z-tbody-tr", playwright.PageLocatorOptions{HasText: "测试服务器-update"})
	c, err := locator.Count()
	if err != nil || c > 0 {
		t.Errorf("Delete server fail: %v", err)
		t.FailNow()
	}
}

func TestUiServer(t *testing.T) {
	runner.Run(t, "客户端-创建服务器", CreateServer)
	runner.Run(t, "客户端-编辑服务器", EditServer)
	runner.Run(t, "客户端-删除服务器", DeleteServer)
}
