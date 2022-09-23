package main

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
	headless := true
	var slowMo float64 = 100
	if languageBrowser == nil || !languageBrowser.IsConnected() {
		languageBrowser, err = pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	}
	if err != nil {
		t.Errorf("Fail to launch the web languageBrowser: %v", err)
		t.FailNow()
	}
	page, err := languageBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
	if _, err = page.Goto("http://127.0.0.1:8081/", playwright.PageGotoOptions{
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
	locator, err := page.Locator(".t-card-toolbar div")
	interpreterTitle, err := locator.InnerText()
	if interpreterTitle != "Interpreter" {
		t.Error("Switch language fail")
		t.FailNow()
	}
	locator, err = page.Locator(".t-card-toolbar button")
	CreateInterpreterTitle, err := locator.InnerText()
	if CreateInterpreterTitle != "Create Interpreter" {
		t.Error("Switch language fail")
		t.FailNow()
	}
	locator, err = page.Locator("#settingModal .modal-title")
	modalTitle, err := locator.InnerText()
	if modalTitle != "Settings" {
		t.Error("Switch language fail")
		t.FailNow()
	}
	err = page.Click(`input[type="radio"]>>nth=0`)
	if err = languageBrowser.Close(); err != nil {
		t.Errorf("The languageBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}

func TestUiLanguage(t *testing.T) {
	runner.Run(t, "设置界面语言", SwitchLanguage)
}
