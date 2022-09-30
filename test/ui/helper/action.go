package helper

import (
	playwright "github.com/playwright-community/playwright-go"
)

var (
	page playwright.Page
)

func OpenUrl(url string) (err error) {
	pw, err := playwright.Run()
	if err != nil {
		panic(err)
	}
	headless := true
	var slowMo float64 = 100
	runBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		panic(err)
	}
	page, err = runBrowser.NewPage()
	if err != nil {
		panic(err)
	}

	if _, err = page.Goto(url, playwright.PageGotoOptions{WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		panic(err)
	}

	return
}
