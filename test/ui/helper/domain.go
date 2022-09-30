package helper

import (
	"github.com/playwright-community/playwright-go"
)

type Webpage struct {
	Page playwright.Page
}

type Element struct {
	Locator playwright.Locator
}

func (p *Webpage) GetElement(selector string, options ...playwright.PageLocatorOptions) (elem Element) {
	locator, err := p.Page.Locator(selector, options...)
	c, err := locator.Count()
	if err != nil || c == 0 {
		panic(err)
	}

	elem = Element{
		Locator: locator,
	}

	return
}
