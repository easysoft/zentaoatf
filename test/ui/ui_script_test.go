package main

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	playwright "github.com/playwright-community/playwright-go"
)

var scriptBrowser playwright.Browser

func SaveScript(t provider.T) {
	t.ID("5470")
	t.AddParentSuite("禅道站点脚本")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	workspaceBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web workspaceBrowser: %v", err)
		t.FailNow()
	}
	page, err := workspaceBrowser.NewPage()
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
	locator, err = page.Locator(".view-line>>text=title=check string matches pattern")
	if err != nil {
		t.Errorf("Find title fail: %v", err)
		t.FailNow()
	}
	var positionX, positionY float64 = 400, 0
	force := true
	err = locator.Click(playwright.PageClickOptions{Force: &force, Position: &playwright.PageClickOptionsPosition{X: &positionX, Y: &positionY}})
	if err != nil {
		t.Errorf("Click title fail: %v", err)
		t.FailNow()
	}
	err = locator.Type("-test")
	if err != nil {
		t.Errorf("Type code fail: %v", err)
		t.FailNow()
	}
	err = page.Click(".tabs-nav-toolbar>>[title=\"Save\"]")
	if err != nil {
		t.Errorf("Click script fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector(".toast-notification-close")
	if err != nil {
		t.Errorf("Wait toast-notification-close fail: %v", err)
		t.FailNow()
	}
	locator, err = page.Locator(".toast-notification-container", playwright.PageLocatorOptions{HasText: "保存成功"})
	c, err = locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Save fail: %v", err)
		t.FailNow()
	}
	err = scriptLocator.Click()
	if err != nil {
		t.Errorf("Click script fail: %v", err)
		t.FailNow()
	}

	locator, err = page.Locator(".view-line>>:has-text('title=check string matches pattern')")
	locator.Click(playwright.PageClickOptions{Force: &force, Position: &playwright.PageClickOptionsPosition{X: &positionX, Y: &positionY}})
	locator.Press("Backspace")
	locator.Press("Backspace")
	locator.Press("Backspace")
	locator.Press("Backspace")
	locator.Press("Backspace")
	page.Click(".tabs-nav-toolbar>>[title=\"Save\"]")
	if err = workspaceBrowser.Close(); err != nil {
		t.Errorf("The workspaceBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}

func ViewScript(t provider.T) {
	t.ID("5469")
	t.AddParentSuite("禅道站点脚本")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	workspaceBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web workspaceBrowser: %v", err)
		t.FailNow()
	}
	page, err := workspaceBrowser.NewPage()
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
	locator, err = page.Locator(".view-line>>text=title=check string matches pattern")
	if err != nil {
		t.Errorf("Find title fail: %v", err)
		t.FailNow()
	}

	if err = workspaceBrowser.Close(); err != nil {
		t.Errorf("The workspaceBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}

func TestUiScript(t *testing.T) {
	runner.Run(t, "客户端-编辑保存禅道站点脚本", SaveScript)
	runner.Run(t, "客户端-显示禅道站点脚本", ViewScript)
}
