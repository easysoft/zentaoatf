package uiTest

import (
	playwright "github.com/playwright-community/playwright-go"
)

func login() {

}

func goToLastUnitTestResult() {

}

func checkUnitTestResult() {

}

func InitZentaoData() (err error) {
	pw, err := playwright.Run()
	if err != nil {
		return
	}
	headless := true
	var slowMo float64 = 100
	runBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		return
	}
	page, err := runBrowser.NewPage()
	if err != nil {
		return
	}
	if _, err = page.Goto("http://127.0.0.1:8081/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		return
	}
	page.WaitForTimeout(10000000)
	return
}
