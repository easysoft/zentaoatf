package plw

import (
	"errors"
	"fmt"
	"time"

	constTestHelper "github.com/easysoft/zentaoatf/test/helper/conf"
	"github.com/easysoft/zentaoatf/test/ui/conf"
	"github.com/easysoft/zentaoatf/test/ui/utils"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	playwright "github.com/playwright-community/playwright-go"
)

func OpenUrl(url string, t provider.T) (ret Webpage, err error) {
	pw, err := playwright.Run()
	utils.PrintErrOrNot(err, t)

	headless := conf.Headless
	var slowMo float64 = 100
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	utils.PrintErrOrNot(err, t)

	page, err := browser.NewPage(playwright.BrowserNewContextOptions{Locale: &conf.Locale})
	utils.PrintErrOrNot(err, t)

	if _, err = page.Goto(url, playwright.PageGotoOptions{WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		page.Close()
		browser.Close()
		utils.PrintErrOrNot(err, t)
	}

	page.SetDefaultTimeout(conf.Timeout)

	ret = Webpage{
		Browser: &browser,
		Pw:      pw,
		Page:    page,
		T:       t,
	}
	return
}

func (p *Webpage) Close() {
	t := p.T
	if p.Page != nil {
		if err := p.Page.Close(); err != nil {
			t.Errorf("The page cannot be closed: %v", err)
			t.FailNow()
		}
	}

	if p.Browser != nil {
		if err := (*p.Browser).Close(); err != nil {
			t.Errorf("The browser cannot be closed: %v", err)
			t.FailNow()
		}
	}

	if p.Pw != nil {
		if err := p.Pw.Stop(); err != nil {
			t.Errorf("The playwright cannot be stopped: %v", err)
			t.FailNow()
		}
	}

	return
}

func (p *Webpage) Goto(url string, options ...playwright.PageGotoOptions) (ret playwright.Response, err error) {
	t := p.T
	ret, err = p.Page.Goto(url, options...)
	if err != nil {
		utils.PrintErrOrNot(err, t)
	}

	return
}

func (p *Webpage) Locator(selector string, options ...playwright.PageLocatorOptions) (ret MyLocator) {
	t := p.T
	locator, err := p.Page.Locator(selector, options...)
	c, err := locator.Count()
	if err == nil && c == 0 {
		err = errors.New("not found")
	}

	if err != nil {
		p.ScreenShot()
		err = errors.New(fmt.Sprintf("Get Locator %s fail: %s", selector, err.Error()))
		utils.PrintErrOrNot(err, t)
	}

	ret = MyLocator{
		PlwLocator: locator,
		Selector:   selector,
		T:          t,
		Page:       p.Page,
	}

	return
}

func (p *Webpage) QuerySelectorAll(selector string) (ret MyElementHandle) {
	t := p.T
	ElementHandles, err := p.Page.QuerySelectorAll(selector)
	if err != nil {
		p.ScreenShot()
		err = errors.New("not found")
	}

	if err != nil {
		p.ScreenShot()
		err = errors.New(fmt.Sprintf("Get Selector %s fail: %s", selector, err.Error()))
		utils.PrintErrOrNot(err, t)
	}

	ret = MyElementHandle{
		ElementHandles: ElementHandles,
		Selector:       selector,
		T:              t,
	}

	return
}

func (p *Webpage) WaitForSelector(selector string, options ...playwright.PageWaitForSelectorOptions) (err error) {
	return p.WaitForSelectorTimeout(selector, 10000, options...)
}

func (p *Webpage) WaitForSelectorTimeout(selector string, timeout float64, options ...playwright.PageWaitForSelectorOptions) (err error) {
	t := p.T
	timeoutOption := playwright.PageWaitForSelectorOptions{Timeout: &timeout}
	options = append(options, timeoutOption)
	_, err = p.Page.WaitForSelector(selector, options...)
	if err != nil {
		p.ScreenShot()
		err = errors.New(fmt.Sprintf("Wait for %s selector fail: %s", selector, err.Error()))
		utils.PrintErrOrNot(err, t)
	}

	return
}

func (p *Webpage) WaitForTimeout(timeout float64) {
	p.Page.WaitForTimeout(timeout)
	return
}

func (p *Webpage) Title() string {
	title, _ := p.Page.Title()

	return title
}

func (p *Webpage) Click(selector string) {
	t := p.T
	err := p.Page.Click(selector, playwright.PageClickOptions{Timeout: &conf.Timeout})
	if err != nil {
		p.ScreenShot()
		err = errors.New(fmt.Sprintf("Click %s fail: %s", selector, err.Error()))
		utils.PrintErrOrNot(err, t)
	}
	return
}

func (p *Webpage) Check(selector string, options ...playwright.FrameCheckOptions) {
	t := p.T
	err := p.Page.Check(selector, options...)
	if err != nil {
		p.ScreenShot()
		err = errors.New(fmt.Sprintf("Check %s fail: %s", selector, err.Error()))
		utils.PrintErrOrNot(err, t)
	}
	return
}

func (p *Webpage) RightClick(selector string) {
	t := p.T
	err := p.Page.Click(selector, playwright.PageClickOptions{Button: playwright.MouseButtonRight, Timeout: &conf.Timeout})
	if err != nil {
		p.ScreenShot()
		err = errors.New(fmt.Sprintf("Rigth click %s fail: %s", selector, err.Error()))
		utils.PrintErrOrNot(err, t)
	}
	return
}

func (p *Webpage) IsHidden(selector string) bool {
	t := p.T
	isHidden, err := p.Page.IsHidden(selector)
	if err != nil {
		p.ScreenShot()
		err = errors.New(fmt.Sprintf("Get %s isHidden fail: %s", selector, err.Error()))
		utils.PrintErrOrNot(err, t)
	}
	return isHidden
}

func (p *Webpage) InnerText(selector string) string {
	t := p.T
	text, err := p.Page.InnerText(selector, playwright.PageInnerTextOptions{Timeout: &conf.Timeout})
	if err != nil {
		p.ScreenShot()
		err = errors.New(fmt.Sprintf("Get %s InnerText fail: %s", selector, err.Error()))
		utils.PrintErrOrNot(err, t)
	}
	return text
}

func (p *Webpage) GetAttribute(selector string, name string, options ...playwright.PageGetAttributeOptions) string {
	t := p.T
	attr, err := p.Page.GetAttribute(selector, name, options...)
	if err != nil {
		p.ScreenShot()
		err = errors.New(fmt.Sprintf("Get %s attribute fail: %s", selector, err.Error()))
		utils.PrintErrOrNot(err, t)
	}
	return attr
}

func (p *Webpage) ScreenShot() {
	if !conf.ShowErr && !conf.ExitAllOnError {
		return
	}
	var screenshotPath = fmt.Sprintf("%stest/screenshot/%v.png", constTestHelper.RootPath, time.Now().Unix())
	p.Page.Screenshot(playwright.PageScreenshotOptions{Path: &screenshotPath})
}

func (p *Webpage) WaitForResponse(url string) (resp playwright.Response) {
	resp = p.Page.WaitForResponse(url, playwright.PageWaitForResponseOptions{Timeout: &conf.Timeout})
	return
}
