package ztfTest

import (
	"fmt"
	"strings"

	commonTestHelper "github.com/easysoft/zentaoatf/test/helper/common"
	"github.com/playwright-community/playwright-go"
)

var expandTimes = 0

func createTestWorkspace(page playwright.Page) {
	err := page.Click(`[title="新建工作目录"]`)
	if err != nil {
		return

	}
	_, err = page.WaitForSelector("#workspaceFormModal")
	locator, err := page.Locator("#workspaceFormModal input")
	if err != nil {
		return
	}
	titleInput, err := locator.Nth(0)
	if err != nil {
		return
	}
	err = titleInput.Fill("单元测试工作目录")
	if err != nil {
		return
	}
	pathInput, err := locator.Nth(1)
	if err != nil {
		return
	}
	workspacePath := fmt.Sprintf("%stest%sdemo%sphp", commonTestHelper.RootPath, commonTestHelper.FilePthSep, commonTestHelper.FilePthSep)
	err = pathInput.Fill(workspacePath)
	if err != nil {
		return
	}
	locator, err = page.Locator("#workspaceFormModal select")
	if err != nil {
		return
	}
	typeInput, err := locator.Nth(0)
	if err != nil {
		return
	}
	_, err = typeInput.SelectOption(playwright.SelectOptionValues{Values: &[]string{"ztf"}})
	if err != nil {
		return
	}
	langInput, err := locator.Nth(1)
	if err != nil {
		return
	}
	_, err = langInput.SelectOption(playwright.SelectOptionValues{Values: &[]string{"php"}})
	if err != nil {
		return
	}
	err = page.Click("#workspaceFormModal>>.modal-action>>span:has-text(\"确定\")")
	if err != nil {
		return
	}
	page.WaitForSelector("#workspaceFormModal", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})
	page.WaitForTimeout(1000)
	locator, err = page.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "testng工作目录"})
	c, err := locator.Count()
	if err != nil || c == 0 {
		return
	}
}

func RunScript(page playwright.Page, scriptName string) {
	locator, err := page.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	c, err := locator.Count()
	if err != nil || c == 0 {
		createTestWorkspace(page)
	}
	locator, err = page.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	c, err = locator.Count()
	err = locator.Click()
	if err != nil {
		return
	}
	scriptLocator, err := locator.Locator("text=" + scriptName)
	if err != nil {
		return
	}
	err = scriptLocator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	if err != nil {
		return
	}
	err = page.Click(".tree-context-menu>>text=执行")
	if err != nil {
		return
	}
	_, err = page.WaitForSelector("#log-list>>.msg-span>>:has-text('执行1个用例，耗时')")
	if err != nil {
		return
	}
	element, err := page.QuerySelector("#log-list>>.msg-span>>:has-text('执行1个用例，耗时')")
	innerText, err := element.InnerText()
	if err != nil {
		return
	}
	if !strings.Contains(innerText, "1(100.0%) 失败") {
		return
	}
	resultTitleElement, err := page.QuerySelector("#rightPane .result-list-item .list-item-title")
	if err != nil {
		return
	}
	resultTitle, err := resultTitleElement.InnerText()
	if err != nil || resultTitle != scriptName {
		return
	}
	timeElement, err := page.QuerySelector("#log-list .item .time")
	if err != nil {
		return
	}
	logTime, err := timeElement.InnerText()
	if err != nil {
		return
	}
	resultTimeElement, err := page.QuerySelector("#rightPane .result-list-item .list-item-trailing-text")
	if err != nil {
		return
	}
	resultTime, err := resultTimeElement.InnerText()
	if err != nil || logTime[:5] != resultTime {
		return
	}
}

func SelectSite(page playwright.Page) (err error) {
	locator, err := page.Locator("#siteMenuToggle")
	if err != nil {
		return
	}
	err = locator.Click()
	if err != nil {
		return
	}
	_, err = page.WaitForSelector("#navbar .list-item")
	if err != nil {
		return
	}
	err = page.Click(".list-item-title>>text=单元测试站点")
	if err != nil {
		return
	}
	var waitTimeOut float64 = 5000
	_, err = page.WaitForSelector("#siteMenuToggle:has-text('单元测试站点')", playwright.PageWaitForSelectorOptions{Timeout: &waitTimeOut})
	if err != nil {
		return
	}
	return nil
}

func ExpandWorspace(page playwright.Page) (err error) {
	var waitTimeOut float64 = 5000
	_, err = page.WaitForSelector(".tree-node:has-text('单元测试工作目录')", playwright.PageWaitForSelectorOptions{Timeout: &waitTimeOut})
	if err != nil {
		createTestWorkspace(page)
	}
	selector, _ := page.QuerySelector(".tree-node:has-text('单元测试工作目录')")
	className, _ := selector.GetAttribute("class")
	if !strings.Contains(className, "collapsed") {
		return
	}

	err = page.Click(".tree-node-title:has-text(\"单元测试工作目录\")")
	if err != nil {
		return
	}
	_, err = page.WaitForSelector(".tree-node-item>>div:has-text('1_string_match.php')", playwright.PageWaitForSelectorOptions{Timeout: &waitTimeOut})
	if err != nil {
		if expandTimes > 5 {
			expandTimes = 0
			return err
		}
		expandTimes++
		ExpandWorspace(page)
		return
	}
	return nil
}

func ExpandProduct(page playwright.Page) (err error) {
	ExpandWorspace(page)
	var waitTimeOut float64 = 5000
	_, err = page.WaitForSelector(".tree-node-item:has-text('product1')", playwright.PageWaitForSelectorOptions{Timeout: &waitTimeOut})
	if err != nil {
		return
	}
	selector, err := page.QuerySelector(".tree-node-root .tree-node:has-text('product1')")
	className, err := selector.GetAttribute("class")
	if !strings.Contains(className, "collapsed") {
		return
	}
	page.Click(".tree-node-item:has-text('product1')")
	page.WaitForTimeout(100)
	selector, err = page.QuerySelector(".tree-node-root .tree-node:has-text('product1')")
	className, err = selector.GetAttribute("class")
	if strings.Contains(className, "collapsed") {
		if expandTimes > 5 {
			expandTimes = 0
			return err
		}
		expandTimes++
		ExpandProduct(page)
		return
	}
	return nil
}

func DeleteWorkspace(page playwright.Page, workspaceName string) {
	locator, err := page.Locator(".tree-node-item", playwright.PageLocatorOptions{HasText: workspaceName})
	c, err := locator.Count()
	if err != nil || c == 0 {
		return
	}
	err = locator.Hover()
	if err != nil {
		return
	}
	err = page.Click(`[title="删除"]`)
	if err != nil {
		return
	}
	err = page.Click(".modal-action>>span:has-text(\"确定\")")
	if err != nil {
		return
	}
	page.WaitForTimeout(1000)
}
