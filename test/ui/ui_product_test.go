package main

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	playwright "github.com/playwright-community/playwright-go"
)

var productBrowser playwright.Browser

func SwitchProduct(t provider.T) {
	t.ID("5496")
	t.AddParentSuite("切换禅道产品")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	productBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web productBrowser: %v", err)
		t.FailNow()
	}
	page, err := productBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
	defer func() {
		if err = productBrowser.Close(); err != nil {
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

	locator, err := page.Locator("#siteMenuToggle")
	if err != nil {
		t.Errorf("The siteMenuToggle is missing: %v", err)
		t.FailNow()
	}
	err = locator.Click()
	if err != nil {
		t.Errorf("Click is fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector("#navbar .list-item")
	if err != nil {
		t.Errorf("Wait for site list nav fail: %v", err)
		t.FailNow()
	}
	err = page.Click(".list-item-title>>text=单元测试站点")
	if err != nil {
		t.Errorf("Click site nav fail: %v", err)
		t.FailNow()
	}
	err = page.Click("#productMenuToggle")
	if err != nil {
		t.Errorf("Click site nav fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector("#navbar .list-item")
	if err != nil {
		t.Errorf("Wait for product list nav fail: %v", err)
		t.FailNow()
	}
	err = page.Click("#navbar .list-item>>text=test")
	page.WaitForTimeout(100)
	productName, err := page.InnerText("#productMenuToggle>>span")
	if productName != "test" {
		t.Error("Switch product fail")
		t.FailNow()
	}
}

func TestUiProduct(t *testing.T) {
	runner.Run(t, "客户端-切换禅道产品", SwitchProduct)
}
