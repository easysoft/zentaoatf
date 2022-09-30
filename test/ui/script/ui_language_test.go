package script

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	playwright "github.com/playwright-community/playwright-go"
)

var languageBrowser playwright.Browser

func SwitchLanguage(t provider.T) {
	t.ID("5464")
	t.AddParentSuite("设置界面语言")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := false
	var slowMo float64 = 100
	languageBrowser, err = pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web languageBrowser: %v", err)
		t.FailNow()
	}
	page, err := languageBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
	defer func() {
		if err = languageBrowser.Close(); err != nil {
			t.Errorf("The languageBrowser cannot be closed: %v", err)
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
		t.Errorf("The Click language nav fail: %v", err)
		t.FailNow()
	}
	err = page.Click(`input[type="radio"]>>nth=1`)
	locator, err := page.Locator(".t-card-toolbar div>>nth=2")
	interpreterTitle, err := locator.InnerText()
	if interpreterTitle != "Remote Server" {
		t.Error("Switch language fail, find interpreter fail")
		t.FailNow()
	}
	locator, err = page.Locator(".t-card-toolbar button")
	CreateInterpreterTitle, err := locator.InnerText()
	if CreateInterpreterTitle != "Create Remote Server" {
		t.Error("Switch language fail, find create remote server btn fail")
		t.FailNow()
	}
	locator, err = page.Locator("#settingModal .modal-title")
	modalTitle, err := locator.InnerText()
	if modalTitle != "Settings" {
		t.Error("Switch language fail, find modalTitle fail")
		t.FailNow()
	}
	err = page.Click(`input[type="radio"]>>nth=0`)
}

func TestUiLanguage(t *testing.T) {
	runner.Run(t, "设置界面语言", SwitchLanguage)
}
