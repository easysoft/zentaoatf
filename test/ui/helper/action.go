package plw

import (
	"errors"
	"fmt"
	"github.com/easysoft/zentaoatf/test/ui/utils"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	playwright "github.com/playwright-community/playwright-go"
)

var (
	browser playwright.Browser
	pw      *playwright.Playwright

	page playwright.Page
)

func OpenUrl(url string, t provider.T) (ret Webpage, err error) {
	pw, err = playwright.Run()
	utils.PrintErrOrNot(err, t)

	headless := false
	var slowMo float64 = 100

	browser, err = pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	utils.PrintErrOrNot(err, t)

	page, err = browser.NewPage()
	utils.PrintErrOrNot(err, t)

	if _, err = page.Goto(url, playwright.PageGotoOptions{WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		utils.PrintErrOrNot(err, t)
	}

	ret = Webpage{
		Page: page,
	}

	return
}

func Close(t provider.T) {

	if browser != nil {
		if err := browser.Close(); err != nil {
			t.Errorf("The resultBrowser cannot be closed: %v", err)
			t.FailNow()
		}
	}

	if pw != nil {
		if err := pw.Stop(); err != nil {
			t.Errorf("The playwright cannot be stopped: %v", err)
			t.FailNow()
		}
	}

	return
}

func (p *Webpage) GetLocator(selector string, t provider.T) (ret MyLocator) {
	locator, err := p.Page.Locator(selector)
	c, err := locator.Count()
	if err == nil && c == 0 {
		err = errors.New("not found")
	}

	err = errors.New(fmt.Sprintf("%s: %s", selector, err.Error()))
	utils.PrintErrOrNot(err, t)

	ret = MyLocator{
		Locator:  locator,
		Selector: selector,
	}

	return
}

func (p *Webpage) GetLocatorByOptions(selector string, options ...playwright.PageLocatorOptions) (elem MyLocator) {
	locator, err := p.Page.Locator(selector, options...)
	c, err := locator.Count()
	if err != nil || c == 0 {
		panic(err)
	}

	elem = MyLocator{
		Locator: locator,
	}

	return
}
