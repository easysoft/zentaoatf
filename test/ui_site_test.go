package main

import (
	"testing"

	playwright "github.com/playwright-community/playwright-go"
)

var siteBrowser playwright.Browser

func CreateSite(t *testing.T) {
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	if siteBrowser == nil || !siteBrowser.IsConnected() {
		siteBrowser, err = pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	}
	if err != nil {
		t.Errorf("Fail to launch the web siteBrowser: %v", err)
		t.FailNow()
	}
	page, err := siteBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
	if _, err = page.Goto("http://127.0.0.1:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
		t.FailNow()
	}
	// page.WaitForSelector(".tree")
	locator, err := page.Locator("#siteMenuToggle")
	if err != nil {
		t.Errorf("The siteMenuToggle is missing: %v", err)
		t.FailNow()
	}
	err = locator.Click()
	if err != nil {
		t.Errorf("The Click is fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector("#navbar .list-item")
	if err != nil {
		t.Errorf("Wait for site list nav fail: %v", err)
		t.FailNow()
	}
	err = page.Click("text=禅道站点管理")
	if err != nil {
		t.Errorf("The Click site nav fail: %v", err)
		t.FailNow()
	}
	// page.WaitForSelector(".modal-title")
	err = page.Click("text=新建站点")
	if err != nil {
		t.Errorf("The Click create site fail: %v", err)
		t.FailNow()
	}
	locator, err = page.Locator("#siteFormModal input")
	if err != nil {
		t.Errorf("Find create site input fail: %v", err)
		t.FailNow()
	}
	titleInput, err := locator.Nth(0)
	if err != nil {
		t.Errorf("Find title input fail: %v", err)
		t.FailNow()
	}
	err = titleInput.Fill("单元测试站点")
	if err != nil {
		t.Errorf("Fil title input fail: %v", err)
		t.FailNow()
	}
	addressInput, err := locator.Nth(1)
	if err != nil {
		t.Errorf("Find address input fail: %v", err)
		t.FailNow()
	}
	err = addressInput.Fill("http://127.0.0.1/zentao/")
	if err != nil {
		t.Errorf("Fil address input fail: %v", err)
		t.FailNow()
	}
	nameInput, err := locator.Nth(2)
	if err != nil {
		t.Errorf("Find name input fail: %v", err)
		t.FailNow()
	}
	err = nameInput.Fill("admin")
	if err != nil {
		t.Errorf("Fil name input fail: %v", err)
		t.FailNow()
	}
	pwdInput, err := locator.Nth(3)
	if err != nil {
		t.Errorf("Find passwd input fail: %v", err)
		t.FailNow()
	}
	err = pwdInput.Fill("123456.")
	if err != nil {
		t.Errorf("Fil passwd input fail: %v", err)
		t.FailNow()
	}
	err = page.Click("text=确定")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		t.FailNow()
	}
	page.WaitForTimeout(1000)
	locator, err = page.Locator(".list-item-content span", playwright.PageLocatorOptions{HasText: "单元测试站点"})
	c, err := locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find created site fail: %v", err)
		t.FailNow()
	}

	if err = siteBrowser.Close(); err != nil {
		t.Errorf("The siteBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}
func EditSite(t *testing.T) {
	// var timeout float64 = 5000
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	if siteBrowser == nil || !siteBrowser.IsConnected() {
		siteBrowser, err = pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	}
	if err != nil {
		t.Errorf("Fail to launch the web siteBrowser: %v", err)
		t.FailNow()
	}
	page, err := siteBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
	if _, err = page.Goto("http://127.0.0.1:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
		t.FailNow()
	}
	locator, err := page.Locator("#siteMenuToggle")
	if err != nil {
		t.Errorf("The siteMenuToggle is missing: %v", err)
		t.FailNow()
	}
	err = locator.Click()
	if err != nil {
		t.Errorf("The Click is fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector("#navbar .list-item")
	if err != nil {
		t.Errorf("Wait for site list nav fail: %v", err)
		t.FailNow()
	}
	err = page.Click("text=禅道站点管理")
	if err != nil {
		t.Errorf("The Click site nav fail: %v", err)
		t.FailNow()
	}
	// page.WaitForSelector(".modal-title")
	locator, err = page.Locator(".list-item", playwright.PageLocatorOptions{HasText: "单元测试站点"})
	c, err := locator.Count()
	if err != nil || c == 0 {
		siteBrowser.Close()
		CreateSite(t)
		EditSite(t)
		return
	}
	locator, err = page.Locator(".list-item", playwright.PageLocatorOptions{HasText: "单元测试站点"})
	if err != nil {
		t.Errorf("Find test site fail: %v", err)
		t.FailNow()
	}
	err = page.Click("text=编辑")
	if err != nil {
		t.Errorf("The Click update site fail: %v", err)
		t.FailNow()
	}
	locator, err = page.Locator("#siteFormModal input")
	if err != nil {
		t.Errorf("Find update site input fail: %v", err)
		t.FailNow()
	}
	titleInput, err := locator.Nth(0)
	if err != nil {
		t.Errorf("Find title input fail: %v", err)
		t.FailNow()
	}
	err = titleInput.Fill("单元测试站点-update")
	if err != nil {
		t.Errorf("Fil title input fail: %v", err)
		t.FailNow()
	}
	addressInput, err := locator.Nth(1)
	if err != nil {
		t.Errorf("Find address input fail: %v", err)
		t.FailNow()
	}
	err = addressInput.Fill("http://127.0.0.1/zentao/")
	if err != nil {
		t.Errorf("Fil address input fail: %v", err)
		t.FailNow()
	}
	nameInput, err := locator.Nth(2)
	if err != nil {
		t.Errorf("Find name input fail: %v", err)
		t.FailNow()
	}
	err = nameInput.Fill("admin")
	if err != nil {
		t.Errorf("Fil name input fail: %v", err)
		t.FailNow()
	}
	pwdInput, err := locator.Nth(3)
	if err != nil {
		t.Errorf("Find passwd input fail: %v", err)
		t.FailNow()
	}
	err = pwdInput.Fill("123456.")
	if err != nil {
		t.Errorf("Fil passwd input fail: %v", err)
		t.FailNow()
	}
	err = page.Click("#siteFormModal>>.modal-action>>span:has-text(\"确定\")")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		t.FailNow()
	}
	page.WaitForTimeout(1000)
	locator, err = page.Locator(".list-item-content", playwright.PageLocatorOptions{HasText: "单元测试站点-update"})
	c, err = locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find update site fail: %v", err)
		t.FailNow()
	}

	if err = siteBrowser.Close(); err != nil {
		t.Errorf("The siteBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}
func DeleteSite(t *testing.T) {
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	if siteBrowser == nil || !siteBrowser.IsConnected() {
		siteBrowser, err = pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	}
	if err != nil {
		t.Errorf("Fail to launch the web siteBrowser: %v", err)
		t.FailNow()
	}
	page, err := siteBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
	if _, err = page.Goto("http://127.0.0.1:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
		t.FailNow()
	}
	// page.WaitForSelector(".tree")
	locator, err := page.Locator("#siteMenuToggle")
	if err != nil {
		t.Errorf("The siteMenuToggle is missing: %v", err)
		t.FailNow()
	}
	err = locator.Click()
	if err != nil {
		t.Errorf("The Click is fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector("#navbar .list-item")
	if err != nil {
		t.Errorf("Wait for site list nav fail: %v", err)
		t.FailNow()
	}
	err = page.Click("text=禅道站点管理")
	if err != nil {
		t.Errorf("The Click site nav fail: %v", err)
		t.FailNow()
	}
	// page.WaitForSelector(".modal-title")
	locator, err = page.Locator(".list-item", playwright.PageLocatorOptions{HasText: "单元测试站点"})
	if err != nil {
		t.Errorf("Find test site fail: %v", err)
		t.FailNow()
	}
	err = page.Click("text=删除")
	if err != nil {
		t.Errorf("The Click delete site fail: %v", err)
		t.FailNow()
	}
	page.WaitForTimeout(1000)
	err = page.Click(":nth-match(.modal-action > button, 1)")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		t.FailNow()
	}
	page.WaitForTimeout(1000)
	locator, err = page.Locator(".list-item-content", playwright.PageLocatorOptions{HasText: "单元测试站点"})
	c, err := locator.Count()
	if err != nil || c > 0 {
		t.Errorf("Delete site fail: %v", err)
		t.FailNow()
	}

	if err = siteBrowser.Close(); err != nil {
		t.Errorf("The siteBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}

func TestUiSite(t *testing.T) {
	t.Run("EditSite", EditSite)
	t.Run("DeleteSite", DeleteSite)
	t.Run("CreateSite", CreateSite)
}
