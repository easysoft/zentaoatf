package plw

import (
	"errors"
	"fmt"
	"time"

	constTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/conf"
	"github.com/easysoft/zentaoatf/cmd/test/ui/conf"

	"github.com/easysoft/zentaoatf/cmd/test/ui/utils"
	playwright "github.com/playwright-community/playwright-go"
)

func (l MyLocator) Click() (err error) {
	err = l.PlwLocator.Click(playwright.PageClickOptions{Timeout: &conf.Timeout})
	t := l.T
	if err != nil {
		l.ScreenShot()
		err = errors.New(fmt.Sprintf("Click %s fail: %s", l.Selector, err.Error()))
		utils.PrintErrOrNot(err, t)
	}

	return
}

func (l MyLocator) ClickWithOption(option playwright.PageClickOptions) (err error) {
	option.Timeout = &conf.Timeout
	err = l.PlwLocator.Click(option)
	t := l.T
	if err != nil {
		l.ScreenShot()
		err = errors.New(fmt.Sprintf("Click %s fail: %s", l.Selector, err.Error()))
		utils.PrintErrOrNot(err, t)
	}

	return
}

func (l MyLocator) RightClick() {
	t := l.T
	err := l.PlwLocator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight, Timeout: &conf.Timeout})
	if err != nil {
		l.ScreenShot()
		err = errors.New(fmt.Sprintf("Rigth click %s fail: %s", l.Selector, err.Error()))
		utils.PrintErrOrNot(err, t)
	}
	return
}

func (l MyLocator) Type(text string, options ...playwright.PageTypeOptions) (err error) {
	err = l.PlwLocator.Type(text, options...)
	t := l.T
	if err != nil {
		l.ScreenShot()
		err = errors.New(fmt.Sprintf("Type %s fail: %s", l.Selector, err.Error()))
		utils.PrintErrOrNot(err, t)
	}

	return
}

func (l MyLocator) Press(text string, options ...playwright.PagePressOptions) (err error) {
	err = l.PlwLocator.Press(text, options...)
	t := l.T
	if err != nil {
		l.ScreenShot()
		err = errors.New(fmt.Sprintf("Press %s fail: %s", l.Selector, err.Error()))
		utils.PrintErrOrNot(err, t)
	}

	return
}

func (l MyLocator) FillNth(nth int, value string, options ...playwright.FrameFillOptions) (err error) {
	t := l.T
	input, err := l.PlwLocator.Nth(nth)
	if err != nil {
		l.ScreenShot()
		err = errors.New(fmt.Sprintf("Fill %s fail: %s", l.Selector, err.Error()))
		utils.PrintErrOrNot(err, t)
	}

	err = input.Fill(value, options...)
	if err != nil {
		l.ScreenShot()
		err = errors.New(fmt.Sprintf("Fill %s fail: %s", l.Selector, err.Error()))
		utils.PrintErrOrNot(err, t)
	}

	return
}

func (l MyLocator) Fill(value string, options ...playwright.FrameFillOptions) (err error) {
	t := l.T
	err = l.PlwLocator.Fill(value, options...)
	if err != nil {
		l.ScreenShot()
		err = errors.New(fmt.Sprintf("Fill %s fail: %s", l.Selector, err.Error()))
		utils.PrintErrOrNot(err, t)
	}

	return
}

func (l MyLocator) SelectNth(nth int, values playwright.SelectOptionValues) (err error) {
	t := l.T
	selectLocator, err := l.PlwLocator.Nth(nth)
	if err != nil {
		l.ScreenShot()
		err = errors.New(fmt.Sprintf("Select %s fail: %s", l.Selector, err.Error()))
		utils.PrintErrOrNot(err, t)
	}

	_, err = selectLocator.SelectOption(values)
	if err != nil {
		l.ScreenShot()
		err = errors.New(fmt.Sprintf("Select %s fail: %s", l.Selector, err.Error()))
		utils.PrintErrOrNot(err, t)
	}

	return
}
func (l MyLocator) Locator(selector string) (ret MyLocator) {
	t := l.T
	locator, err := l.PlwLocator.Locator(selector)
	c, err := locator.Count()
	if err == nil && c == 0 {
		err = errors.New("not found")
	}

	if err != nil {
		l.ScreenShot()
		err = errors.New(fmt.Sprintf("Get Locator %s fail: %s", selector, err.Error()))
		utils.PrintErrOrNot(err, t)
	}

	ret = MyLocator{
		PlwLocator: locator,
		Selector:   selector,
		T:          t,
	}
	return
}

func (l MyLocator) InnerText() string {
	t := l.T
	text, err := l.PlwLocator.InnerText()
	if err != nil {
		l.ScreenShot()
		err = errors.New(fmt.Sprintf("Get %s InnerText fail: %s", l.Selector, err.Error()))
		utils.PrintErrOrNot(err, t)
	}
	return text
}

func (l MyLocator) Count() int {
	t := l.T
	count, err := l.PlwLocator.Count()
	if err != nil {
		l.ScreenShot()
		err = errors.New(fmt.Sprintf("Get %s count fail: %s", l.Selector, err.Error()))
		utils.PrintErrOrNot(err, t)
	}
	return count
}

func (l MyLocator) Hover(options ...playwright.PageHoverOptions) {
	t := l.T
	err := l.PlwLocator.Hover(options...)
	if err != nil {
		l.ScreenShot()
		err = errors.New(fmt.Sprintf("Hover %s fail: %s", l.Selector, err.Error()))
		utils.PrintErrOrNot(err, t)
	}
	return
}

func (l *MyLocator) ScreenShot() {
	if !conf.ShowErr && !conf.ExitAllOnError {
		return
	}
	var screenshotPath = fmt.Sprintf("%stest/screenshot/%v.png", constTestHelper.RootPath, time.Now().Unix())
	l.Page.Screenshot(playwright.PageScreenshotOptions{Path: &screenshotPath})
}

func (l MyElementHandle) InnerText(nth int) string {
	t := l.T
	text, err := l.ElementHandles[nth].InnerText()
	if err != nil {
		err = errors.New(fmt.Sprintf("Get %s InnerText fail: %s", l.Selector, err.Error()))
		utils.PrintErrOrNot(err, t)
	}
	return text
}

func (l MyElementHandle) GetAttribute(nth int, name string) string {
	t := l.T
	if len(l.ElementHandles) == 0 {
		return ""
	}
	text, err := l.ElementHandles[nth].GetAttribute(name)
	if err != nil {
		err = errors.New(fmt.Sprintf("Get %s Attribute fail: %s", l.Selector, err.Error()))
		utils.PrintErrOrNot(err, t)
	}
	return text
}
