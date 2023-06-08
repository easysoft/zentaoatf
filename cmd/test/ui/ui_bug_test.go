package main

import (
	"fmt"
	"testing"

	constTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/conf"
	ztfTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/ztf"
	"github.com/easysoft/zentaoatf/cmd/test/ui/conf"
	plwHelper "github.com/easysoft/zentaoatf/cmd/test/ui/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	playwright "github.com/playwright-community/playwright-go"
)

func ScriptBug(t provider.T) {
	t.ID("5747")
	t.AddParentSuite("查看bug列表")

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage, "")
	ztfTestHelper.ExpandWorspace(webpage)

	ztfTestHelper.RunScript(webpage, "1_string_match.php")
	ztfTestHelper.SubmitResult(webpage)

	webpage.Click(".tree-node-title:has-text('1_string_match.php')")

	webpage.WaitForResponse("**/bugs*")
	webpage.WaitForSelectorTimeout(".statistic>>span>>nth=3", 3000)

	webpage.Click(".statistic>>span>>nth=3")
	webpage.WaitForSelectorTimeout("#bugsModal>>tr", 3000)

	elements := webpage.QuerySelectorAll("#bugsModal>>tr")
	bugTimesInt := len(elements.ElementHandles)
	if bugTimesInt < 2 {
		webpage.ScreenShot()
		t.Error("View script bug error")
		t.FailNow()
	}
}

func ScriptsBug(t provider.T) {
	t.ID("5748")
	t.AddParentSuite("查看bug列表")

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage, "")
	webpage.WaitForSelector(fmt.Sprintf("#siteMenuToggle:has-text('%s')", constTestHelper.SiteName), playwright.PageWaitForSelectorOptions{Timeout: &conf.Timeout})
	ztfTestHelper.ExpandWorspace(webpage)

	webpage.Click(`[title="批量选择"]`)
	webpage.Click(".tree-node-item:has-text('1_string_match.php')>>.tree-node-check")
	webpage.Click(".tree-node-item:has-text('2_webpage_extract.php')>>.tree-node-check")
	webpage.Click(`[title="禅道BUG"]`)
	elements := webpage.QuerySelectorAll("#bugsModal>>tr")

	bugTimesInt := len(elements.ElementHandles)
	if bugTimesInt < 2 {
		t.Error("statistic error")
		t.FailNow()
	}
}
func TestUiBug(t *testing.T) {
	runner.Run(t, "客户端-查看单个脚本bug列表", ScriptBug)
	runner.Run(t, "客户端-查看选中脚本bug列表", ScriptsBug)
}
