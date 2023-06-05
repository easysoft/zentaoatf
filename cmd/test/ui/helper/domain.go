package plw

import (
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/playwright-community/playwright-go"
)

type Webpage struct {
	Browser *playwright.Browser
	Pw      *playwright.Playwright
	Page    playwright.Page
	T       provider.T
}

type MyLocator struct {
	Selector   string
	PlwLocator playwright.Locator
	T          provider.T
	Page       playwright.Page
}

type MyElementHandle struct {
	Selector       string
	ElementHandles []playwright.ElementHandle
	T              provider.T
}
