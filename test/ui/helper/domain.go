package plw

import (
	"github.com/playwright-community/playwright-go"
)

type Webpage struct {
	Page playwright.Page
}

type MyLocator struct {
	Selector string
	Locator  playwright.Locator
}
