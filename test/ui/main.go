package main

import (
	"flag"
	"fmt"
	"testing"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	serverConfig "github.com/easysoft/zentaoatf/internal/server/config"
	commonTestHelper "github.com/easysoft/zentaoatf/test/helper/common"
	ztfTest "github.com/easysoft/zentaoatf/test/helper/ztf"
	"github.com/playwright-community/playwright-go"
)

func main() {
	serverConfig.InitLog()
	commConsts.Verbose = true
	var version = flag.String("zentaoVersion", "", "")
	testing.Init()
	flag.Parse()
	fmt.Println(*version)
	err := commonTestHelper.InitZentao(*version)
	if err != nil {
		fmt.Println("Init zentao data fail ", err)
	}
	err = commonTestHelper.Pull()
	if err != nil {
		fmt.Println("Git pull code fail ", err)
	}
	err = commonTestHelper.BuildCli()
	if err != nil {
		fmt.Println("Build cli fail ", err)
	}
	err = commonTestHelper.RunServer()
	if err != nil {
		fmt.Println("Build server fail ")
	}
	err = commonTestHelper.RunUi()
	if err != nil {
		fmt.Println("Build server fail ")
	}
	commonTestHelper.TestUi()
}

func init() {
	pw, err := playwright.Run()
	if err != nil {
		return
	}
	headless := true
	var slowMo float64 = 100
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		return
	}
	page, err := browser.NewPage()
	if err != nil {
		return
	}
	if _, err = page.Goto("http://127.0.0.1:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		return
	}
	ztfTest.SelectSite(page)
	ztfTest.ExpandWorspace(page)
	if err = browser.Close(); err != nil {
		return
	}
	if err = pw.Stop(); err != nil {
		return
	}
}
