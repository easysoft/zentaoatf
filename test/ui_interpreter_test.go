package main

import (
	"testing"

	playwright "github.com/playwright-community/playwright-go"
)

var interpreterBrowser playwright.Browser

func CreateInterpreter(t *testing.T) {
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()

	}
	headless := false
	var slowMo float64 = 100
	if interpreterBrowser == nil || !interpreterBrowser.IsConnected() {
		interpreterBrowser, err = pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	}
	if err != nil {
		t.Errorf("Fail to launch the web interpreterBrowser: %v", err)
		t.FailNow()

	}
	page, err := interpreterBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()

	}
	if _, err = page.Goto("http://127.0.0.1:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
		t.FailNow()

	}

	err = page.Click("#navbar>>[title=\"设置\"]")
	if err != nil {
		t.Errorf("The Click interpreter nav fail: %v", err)
		t.FailNow()

	}
	err = page.Click("text=新建运行环境")
	if err != nil {
		t.Errorf("The Click create interpreter fail: %v", err)
		t.FailNow()

	}
	Locator, err := page.Locator("#interpreterFormModal select")
	if err != nil {
		t.Errorf("Find create interpreter input fail: %v", err)
		t.FailNow()

	}
	langSelect, err := Locator.Nth(0)
	if err != nil {
		t.Errorf("Find lang select fail: %v", err)
		t.FailNow()
	}
	_, err = langSelect.SelectOption(playwright.SelectOptionValues{Values: &[]string{"python"}})
	if err != nil {
		t.Errorf("Select lang select fail: %v", err)
		t.FailNow()
	}
	page.WaitForTimeout(200)
	pathSelect, err := Locator.Nth(1)
	if err != nil {
		t.Errorf("Find address input fail: %v", err)
		t.FailNow()

	}
	_, err = pathSelect.SelectOption(playwright.SelectOptionValues{Indexes: &[]int{1}})
	if err != nil {
		t.Errorf("Fil address input fail: %v", err)
		t.FailNow()

	}

	err = page.Click("#interpreterFormModal>>text=确定")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		t.FailNow()

	}
	Locator, err = page.Locator("#settingModal .z-tbody-td", playwright.PageLocatorOptions{HasText: "Python"})
	c, err := Locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find created interpreter fail: %v", err)
		t.FailNow()
	}
}
func EditInterpreter(t *testing.T) {
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()

	}
	headless := false
	var slowMo float64 = 100
	if interpreterBrowser == nil || !interpreterBrowser.IsConnected() {
		interpreterBrowser, err = pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	}
	defer interpreterBrowser.Close()
	defer pw.Stop()
	if err != nil {
		t.Errorf("Fail to launch the web interpreterBrowser: %v", err)
		t.FailNow()

	}
	page, err := interpreterBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()

	}
	if _, err = page.Goto("http://127.0.0.1:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
		t.FailNow()

	}

	err = page.Click("#navbar>>[title=\"设置\"]")
	if err != nil {
		t.Errorf("The Click interpreter nav fail: %v", err)
		t.FailNow()

	}
	Locator, err := page.Locator("#settingModal .z-tbody-tr", playwright.PageLocatorOptions{HasText: "Python"})
	if err != nil {
		t.Errorf("Find python tr fail: %v", err)
		t.FailNow()

	}
	Locator, err = Locator.Locator("text=编辑")
	if err != nil {
		t.Errorf("Find python edit btn fail: %v", err)
		t.FailNow()

	}
	err = Locator.Click()
	if err != nil {
		t.Errorf("The Click update site fail: %v", err)
		t.FailNow()

	}
	Locator, err = page.Locator("#interpreterFormModal select")
	if err != nil {
		t.Errorf("Find create interpreter input fail: %v", err)
		t.FailNow()

	}
	langSelect, err := Locator.Nth(0)
	if err != nil {
		t.Errorf("Find lang select fail: %v", err)
		t.FailNow()
	}
	_, err = langSelect.SelectOption(playwright.SelectOptionValues{Values: &[]string{"python"}})
	if err != nil {
		t.Errorf("Select lang select fail: %v", err)
		t.FailNow()
	}
	page.WaitForTimeout(200)
	pathSelect, err := Locator.Nth(1)
	if err != nil {
		t.Errorf("Find address input fail: %v", err)
		t.FailNow()

	}
	_, err = pathSelect.SelectOption(playwright.SelectOptionValues{Indexes: &[]int{1}})
	if err != nil {
		t.Errorf("Fil address input fail: %v", err)
		t.FailNow()

	}

	err = page.Click("#interpreterFormModal>>text=确定")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		t.FailNow()

	}
	Locator, err = page.Locator("#settingModal .z-tbody-td", playwright.PageLocatorOptions{HasText: "Python"})
	c, err := Locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find created interpreter fail: %v", err)
		t.FailNow()
	}
}
func DeleteInterpreter(t *testing.T) {
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()

	}
	headless := false
	var slowMo float64 = 100
	if interpreterBrowser == nil || !interpreterBrowser.IsConnected() {
		interpreterBrowser, err = pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	}
	defer interpreterBrowser.Close()
	defer pw.Stop()
	if err != nil {
		t.Errorf("Fail to launch the web interpreterBrowser: %v", err)
		t.FailNow()

	}
	page, err := interpreterBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()

	}
	if _, err = page.Goto("http://127.0.0.1:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
		t.FailNow()

	}

	err = page.Click("#navbar>>[title=\"设置\"]")
	if err != nil {
		t.Errorf("The Click interpreter nav fail: %v", err)
		t.FailNow()

	}
	Locator, err := page.Locator("#settingModal .z-tbody-tr", playwright.PageLocatorOptions{HasText: "Python"})
	if err != nil {
		t.Errorf("Find python tr fail: %v", err)
		t.FailNow()

	}
	Locator, err = Locator.Locator("text=删除")
	if err != nil {
		t.Errorf("Find python edit btn fail: %v", err)
		t.FailNow()

	}
	err = Locator.Click()
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
	Locator, err = page.Locator("#settingModal .z-tbody-tr", playwright.PageLocatorOptions{HasText: "Python"})
	c, err := Locator.Count()
	if err != nil || c > 0 {
		t.Errorf("Delete interpreter fail: %v", err)
		t.FailNow()

	}

	if err = interpreterBrowser.Close(); err != nil {
		t.Errorf("The interpreterBrowser cannot be closed: %v", err)
		t.FailNow()

	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()

	}
}

func TestUiInterpreter(t *testing.T) {
	t.Run("CreateInterpreter", CreateInterpreter)
	t.Run("EditInterpreter", EditInterpreter)
	t.Run("DeleteInterpreter", DeleteInterpreter)
}
