package main

import (
	"log"
	"testing"

	playwright "github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/suite"
)

type UiSiteSuit struct {
	suite.Suite
	testCount uint32
}

func (s *UiSiteSuit) TestUiSiteSuite() {
	runTest()
}
func runTest() {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not launch playwright: %v", err)
	}
	browser, err := pw.Chromium.Launch()
	if err != nil {
		log.Fatalf("Fail to launch the web browser: %v", err)
	}
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("Error to create the new page: %v", err)
	}
	if _, err = page.Goto("https://ithelp.ithome.com.tw/2020-12th-ironman?utm_source=iThelp&utm_medium=nav&utm_campaign=ironman", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	}); err != nil {
		log.Fatalf("The specific URL is missing: %v", err)
	}
	if _, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String("ironmenGo2020-web.png"),
	}); err != nil {
		log.Fatalf("screenshot cannot be created: %v", err)
	}
	if err = browser.Close(); err != nil {
		log.Fatalf("The browser cannot be closed: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Fatalf("The playwright cannot be stopped: %v", err)
	}
}

func TestUiSite(t *testing.T) {
	suite.Run(t, new(UiSiteSuit))
}
