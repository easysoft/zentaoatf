package ztfTest

import (
	"fmt"
	"strings"

	constTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/conf"
	plwConf "github.com/easysoft/zentaoatf/cmd/test/ui/conf"
	plwHelper "github.com/easysoft/zentaoatf/cmd/test/ui/helper"
	"github.com/playwright-community/playwright-go"
)

var expandTimes = 0
var AddSiteTimes = 0

func CreateTestWorkspace(webpage plwHelper.Webpage, name, workspacePath string) {
	if workspacePath == "" {
		workspacePath = fmt.Sprintf("%stest%sdemo%sphp", constTestHelper.RootPath, constTestHelper.FilePthSep, constTestHelper.FilePthSep)
	}

	webpage.Click(`[title="新建工作目录"]`)
	webpage.WaitForSelector("#workspaceFormModal")
	locator := webpage.Locator("#workspaceFormModal input")
	locator.FillNth(0, name)
	locator.FillNth(1, workspacePath)
	locator = webpage.Locator("#workspaceFormModal select")
	locator.SelectNth(0, playwright.SelectOptionValues{Values: &[]string{"ztf"}})
	locator.SelectNth(1, playwright.SelectOptionValues{Values: &[]string{"php"}})
	webpage.Click("#workspaceFormModal>>.modal-action>>span:has-text(\"确定\")")
	webpage.WaitForSelector("#workspaceFormModal", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})
	webpage.WaitForTimeout(1000)
	webpage.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "testng工作目录"})
}

func RunScript(webpage plwHelper.Webpage, scriptName string) {
	locator := webpage.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	c := locator.Count()
	if c == 0 {
		CreateTestWorkspace(webpage, "单元测试工作目录", "")
	}
	locator = webpage.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	locator.Click()
	scriptLocator := locator.Locator("text=" + scriptName)
	scriptLocator.RightClick()
	webpage.Click(".tree-context-menu>>text=执行")
	webpage.WaitForSelector("#log-list>>.msg-span>>:has-text('执行1个用例，耗时')")
	element := webpage.QuerySelectorAll("#log-list>>.msg-span>>:has-text('执行1个用例，耗时')")
	innerText := element.InnerText(0)
	if !strings.Contains(innerText, "1(100.0%) 失败") {
		return
	}
	resultTitleElement := webpage.QuerySelectorAll("#rightPane .result-list-item .list-item-title")
	resultTitle := resultTitleElement.InnerText(0)
	if resultTitle != scriptName {
		return
	}
	timeElement := webpage.QuerySelectorAll("#log-list .item .time")
	logTime := timeElement.InnerText(0)
	resultTimeElement := webpage.QuerySelectorAll("#rightPane .result-list-item .list-item-trailing-text")
	resultTime := resultTimeElement.InnerText(0)
	if logTime[:5] != resultTime {
		return
	}
}

func SelectSite(webpage plwHelper.Webpage) (err error) {
	plwConf.DisableErr()
	defer plwConf.EnableErr()
	webpage.Click("#siteMenuToggle")
	webpage.WaitForSelectorTimeout("#navbar>>.list-item-title>>text=单元测试站点", 3000)
	locator := webpage.Locator(".list-item-title>>text=单元测试站点")
	if locator.Count() == 0 {
		AddSiteTimes++
		if AddSiteTimes > 2 {
			return
		}
		CreateSite(webpage)
		SelectSite(webpage)
		return
	}
	webpage.Click(".list-item-title>>text=单元测试站点")
	var waitTimeOut float64 = 5000
	webpage.WaitForSelector("#siteMenuToggle:has-text('单元测试站点')", playwright.PageWaitForSelectorOptions{Timeout: &waitTimeOut})
	return nil
}

func CreateSite(webpage plwHelper.Webpage) {
	webpage.WaitForSelector("#siteMenuToggle")
	webpage.Click("#siteMenuToggle")
	webpage.WaitForSelector("#navbar .list-item")
	webpage.Click("text=禅道站点管理")
	webpage.Click("text=新建站点")
	locator := webpage.Locator("#siteFormModal input")
	locator.FillNth(0, "单元测试站点")
	locator.FillNth(1, constTestHelper.ZentaoSiteUrl)
	locator.FillNth(2, constTestHelper.ZentaoUsername)
	locator.FillNth(3, constTestHelper.ZentaoPassword)
	webpage.Click("text=确定")
	webpage.WaitForSelector(".list-item-content span:has-text('单元测试站点')")
	locator = webpage.Locator(".list-item-content span", playwright.PageLocatorOptions{HasText: "单元测试站点"})
	webpage.Click("#siteModal>>.modal-close")
}

func ExpandWorspace(webpage plwHelper.Webpage) (err error) {
	if !webpage.ElementExist(".tree-node-title:has-text('单元测试工作目录')") {
		CreateTestWorkspace(webpage, "单元测试工作目录", "")
	}

	selector := webpage.QuerySelectorAll(".tree-node-root:has-text('单元测试工作目录')")
	className := selector.GetAttribute(0, "class")
	if className != "" && !strings.Contains(className, "collapsed") {
		return
	}

	webpage.Click(".tree-node-title:has-text(\"单元测试工作目录\")")
	err = webpage.WaitForSelectorTimeout(".tree-node-item>>div:has-text('1_string_match.php')", 3000)
	if err != nil {
		if expandTimes > 3 {
			expandTimes = 0
			return err
		}
		expandTimes++
		ExpandWorspace(webpage)
		return
	}
	return nil
}

func ExpandProduct(webpage plwHelper.Webpage) (err error) {
	plwConf.DisableErr()
	defer plwConf.EnableErr()
	ExpandWorspace(webpage)
	var waitTimeOut float64 = 5000
	webpage.WaitForSelector(".tree-node-item:has-text('product1')", playwright.PageWaitForSelectorOptions{Timeout: &waitTimeOut})
	selector := webpage.QuerySelectorAll(".tree-node-root .tree-node:has-text('product1')")
	className := selector.GetAttribute(0, "class")
	if className != "" && !strings.Contains(className, "collapsed") {
		return
	}
	webpage.Click(".tree-node-item:has-text('product1')")
	webpage.WaitForTimeout(100)
	selector = webpage.QuerySelectorAll(".tree-node-root .tree-node:has-text('product1')")
	className = selector.GetAttribute(0, "class")
	if className != "" && strings.Contains(className, "collapsed") {
		if expandTimes > 5 {
			expandTimes = 0
			return err
		}
		expandTimes++
		ExpandProduct(webpage)
		return
	}
	return nil
}

func DeleteWorkspace(webpage plwHelper.Webpage, workspaceName string) {
	locator := webpage.Locator(".tree-node-item", playwright.PageLocatorOptions{HasText: workspaceName})
	c := locator.Count()
	if c == 0 {
		return
	}
	locator.Hover()
	locator = locator.Locator(`[title="删除"]`)
	locator.Click()
	webpage.Click(".modal-action>>span:has-text(\"确定\")")
	webpage.WaitForTimeout(1000)
}

func SubmitResult(webpage plwHelper.Webpage) {
	webpage.Click("#rightPane .result-list-item .list-item-title>>nth=0")
	webpage.Click(".page-result .btn:has-text('提交缺陷到禅道')")
	webpage.WaitForSelector("#submitBugModal")
	webpage.Click("#submitBugModal>>text=确定")
	webpage.WaitForSelector("#submitBugModal", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateHidden})
	locator := webpage.Locator(".toast-notification-container", playwright.PageLocatorOptions{HasText: "提交成功"})
	locator.Click()
}
