package main

import (
	"testing"

	playwright "github.com/playwright-community/playwright-go"
)

func CreateSite(t *testing.T) {
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		return
	}
	headless := true
	var slowMo float64 = 100
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web browser: %v", err)
		return
	}
	page, err := browser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		return
	}
	if _, err = page.Goto("http://app.me:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
		return
	}
	// page.WaitForSelector(".tree")
	Locator, err := page.Locator("#siteMenuToggle")
	if err != nil {
		t.Errorf("The siteMenuToggle is missing: %v", err)
		return
	}
	err = Locator.Click()
	if err != nil {
		t.Errorf("The Click is fail: %v", err)
		return
	}
	_, err = page.WaitForSelector("#navbar .list-item")
	if err != nil {
		t.Errorf("Wait for site list nav fail: %v", err)
		return
	}
	err = page.Click("text=禅道站点管理")
	if err != nil {
		t.Errorf("The Click site nav fail: %v", err)
		return
	}
	// page.WaitForSelector(".modal-title")
	err = page.Click("text=新建站点")
	if err != nil {
		t.Errorf("The Click create site fail: %v", err)
		return
	}
	Locator, err = page.Locator("#siteFormModal input")
	if err != nil {
		t.Errorf("Find create site input fail: %v", err)
		return
	}
	titleInput, err := Locator.Nth(0)
	if err != nil {
		t.Errorf("Find title input fail: %v", err)
		return
	}
	err = titleInput.Fill("单元测试站点")
	if err != nil {
		t.Errorf("Fil title input fail: %v", err)
		return
	}
	addressInput, err := Locator.Nth(1)
	if err != nil {
		t.Errorf("Find address input fail: %v", err)
		return
	}
	err = addressInput.Fill("http://pms.test/")
	if err != nil {
		t.Errorf("Fil address input fail: %v", err)
		return
	}
	nameInput, err := Locator.Nth(2)
	if err != nil {
		t.Errorf("Find name input fail: %v", err)
		return
	}
	err = nameInput.Fill("admin")
	if err != nil {
		t.Errorf("Fil name input fail: %v", err)
		return
	}
	pwdInput, err := Locator.Nth(3)
	if err != nil {
		t.Errorf("Find passwd input fail: %v", err)
		return
	}
	err = pwdInput.Fill("123456.")
	if err != nil {
		t.Errorf("Fil passwd input fail: %v", err)
		return
	}
	err = page.Click("text=确定")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		return
	}
	page.WaitForTimeout(1000)
	Locator, err = page.Locator(".list-item-content span", playwright.PageLocatorOptions{HasText: "单元测试站点"})
	c, err := Locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find created site fail: %v", err)
		return
	}
	if _, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String("site_create.png"),
	}); err != nil {
		t.Errorf("screenshot cannot be created: %v", err)
		return
	}
	// page.WaitForTimeout(1000000)
	if err = browser.Close(); err != nil {
		t.Errorf("The browser cannot be closed: %v", err)
		return
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		return
	}
}
func EditSite(t *testing.T) {
	// var timeout float64 = 5000
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		return
	}
	headless := true
	var slowMo float64 = 100
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web browser: %v", err)
		return
	}
	page, err := browser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		return
	}
	if _, err = page.Goto("http://app.me:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
		return
	}
	// page.WaitForSelector(".tree")
	Locator, err := page.Locator("#siteMenuToggle")
	if err != nil {
		t.Errorf("The siteMenuToggle is missing: %v", err)
		return
	}
	err = Locator.Click()
	if err != nil {
		t.Errorf("The Click is fail: %v", err)
		return
	}
	_, err = page.WaitForSelector("#navbar .list-item")
	if err != nil {
		t.Errorf("Wait for site list nav fail: %v", err)
		return
	}
	err = page.Click("text=禅道站点管理")
	if err != nil {
		t.Errorf("The Click site nav fail: %v", err)
		return
	}
	// page.WaitForSelector(".modal-title")
	Locator, err = page.Locator(".list-item", playwright.PageLocatorOptions{HasText: "单元测试站点"})
	if err != nil {
		t.Errorf("Find test site fail: %v", err)
		return
	}
	err = page.Click("text=编辑")
	if err != nil {
		t.Errorf("The Click update site fail: %v", err)
		return
	}
	Locator, err = page.Locator("#siteFormModal input")
	if err != nil {
		t.Errorf("Find update site input fail: %v", err)
		return
	}
	titleInput, err := Locator.Nth(0)
	if err != nil {
		t.Errorf("Find title input fail: %v", err)
		return
	}
	err = titleInput.Fill("单元测试站点-update")
	if err != nil {
		t.Errorf("Fil title input fail: %v", err)
		return
	}
	addressInput, err := Locator.Nth(1)
	if err != nil {
		t.Errorf("Find address input fail: %v", err)
		return
	}
	err = addressInput.Fill("http://pms.test/")
	if err != nil {
		t.Errorf("Fil address input fail: %v", err)
		return
	}
	nameInput, err := Locator.Nth(2)
	if err != nil {
		t.Errorf("Find name input fail: %v", err)
		return
	}
	err = nameInput.Fill("admin")
	if err != nil {
		t.Errorf("Fil name input fail: %v", err)
		return
	}
	pwdInput, err := Locator.Nth(3)
	if err != nil {
		t.Errorf("Find passwd input fail: %v", err)
		return
	}
	err = pwdInput.Fill("123456.")
	if err != nil {
		t.Errorf("Fil passwd input fail: %v", err)
		return
	}
	err = page.Click("#siteFormModal>>.modal-action>>span:has-text(\"确定\")")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		return
	}
	page.WaitForTimeout(1000)
	Locator, err = page.Locator(".list-item-content", playwright.PageLocatorOptions{HasText: "单元测试站点-update"})
	c, err := Locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find update site fail: %v", err)
		return
	}
	if _, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String("site_update.png"),
	}); err != nil {
		t.Errorf("screenshot cannot be created: %v", err)
		return
	}
	// page.WaitForTimeout(1000000)
	if err = browser.Close(); err != nil {
		t.Errorf("The browser cannot be closed: %v", err)
		return
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		return
	}
}
func DeleteSite(t *testing.T) {
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		return
	}
	headless := true
	var slowMo float64 = 100
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web browser: %v", err)
		return
	}
	page, err := browser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		return
	}
	if _, err = page.Goto("http://app.me:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
		return
	}
	// page.WaitForSelector(".tree")
	Locator, err := page.Locator("#siteMenuToggle")
	if err != nil {
		t.Errorf("The siteMenuToggle is missing: %v", err)
		return
	}
	err = Locator.Click()
	if err != nil {
		t.Errorf("The Click is fail: %v", err)
		return
	}
	_, err = page.WaitForSelector("#navbar .list-item")
	if err != nil {
		t.Errorf("Wait for site list nav fail: %v", err)
		return
	}
	err = page.Click("text=禅道站点管理")
	if err != nil {
		t.Errorf("The Click site nav fail: %v", err)
		return
	}
	// page.WaitForSelector(".modal-title")
	Locator, err = page.Locator(".list-item", playwright.PageLocatorOptions{HasText: "单元测试站点"})
	if err != nil {
		t.Errorf("Find test site fail: %v", err)
		return
	}
	err = page.Click("text=删除")
	if err != nil {
		t.Errorf("The Click delete site fail: %v", err)
		return
	}
	page.WaitForTimeout(1000)
	err = page.Click(":nth-match(.modal-action > button, 1)")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		return
	}
	page.WaitForTimeout(1000)
	Locator, err = page.Locator(".list-item-content", playwright.PageLocatorOptions{HasText: "单元测试站点"})
	c, err := Locator.Count()
	if err != nil || c > 0 {
		t.Errorf("Delete site fail: %v", err)
		return
	}
	if _, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String("site_delete.png"),
	}); err != nil {
		t.Errorf("screenshot cannot be created: %v", err)
		return
	}
	// page.WaitForTimeout(1000000)
	if err = browser.Close(); err != nil {
		t.Errorf("The browser cannot be closed: %v", err)
		return
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		return
	}
}

func TestUiSite(t *testing.T) {
	t.Run("CreateSite", CreateSite)
	t.Run("EditSite", EditSite)
	t.Run("DeleteSite", DeleteSite)
}
