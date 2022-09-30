package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	commonTestHelper "github.com/easysoft/zentaoatf/test/helper/common"
	ztfTestHelper "github.com/easysoft/zentaoatf/test/helper/ztf"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	playwright "github.com/playwright-community/playwright-go"
)

var pw, err = os.Getwd()
var (
	workspacePath = fmt.Sprintf("%stest%sdemo%sphp", commonTestHelper.RootPath, commonTestHelper.FilePthSep, commonTestHelper.FilePthSep)
)
var workspacePage playwright.Page

func CreateWorkspace(t provider.T) {
	t.ID("5468")
	t.AddParentSuite("管理禅道站点下工作目录")
	err = workspacePage.Click(`[title="新建工作目录"]`)
	if err != nil {
		t.Errorf("The Click create workspace fail: %v", err)
		t.FailNow()
		return
	}
	_, err = workspacePage.WaitForSelector("#workspaceFormModal")
	locator, err := workspacePage.Locator("#workspaceFormModal input")
	if err != nil {
		t.Errorf("Find create workspace input fail: %v", err)
		t.FailNow()
		return
	}
	titleInput, err := locator.Nth(0)
	if err != nil {
		t.Errorf("Find title input fail: %v", err)
		t.FailNow()
		return
	}
	err = titleInput.Fill("单元测试工作目录")
	if err != nil {
		t.Errorf("Fil title input fail: %v", err)
		t.FailNow()
		return
	}
	pathInput, err := locator.Nth(1)
	if err != nil {
		t.Errorf("Find address input fail: %v", err)
		t.FailNow()
		return
	}
	err = pathInput.Fill(workspacePath)
	if err != nil {
		t.Errorf("Fil address input fail: %v", err)
		t.FailNow()
		return
	}
	locator, err = workspacePage.Locator("#workspaceFormModal select")
	if err != nil {
		t.Errorf("Find create workspace select fail: %v", err)
		t.FailNow()
		return
	}
	typeInput, err := locator.Nth(0)
	if err != nil {
		t.Errorf("Find name input fail: %v", err)
		t.FailNow()
		return
	}
	_, err = typeInput.SelectOption(playwright.SelectOptionValues{Values: &[]string{"ztf"}})
	if err != nil {
		t.Errorf("Fil name input fail: %v", err)
		t.FailNow()
		return
	}
	langInput, err := locator.Nth(1)
	if err != nil {
		t.Errorf("Find lang input fail: %v", err)
		t.FailNow()
		return
	}
	_, err = langInput.SelectOption(playwright.SelectOptionValues{Values: &[]string{"php"}})
	if err != nil {
		t.Errorf("Fil lang input fail: %v", err)
		t.FailNow()
		return
	}
	err = workspacePage.Click("#workspaceFormModal>>.modal-action>>span:has-text(\"确定\")")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		t.FailNow()
		return
	}
	var waitTimeOut float64 = 5000
	_, err = workspacePage.WaitForSelector(".tree-node", playwright.PageWaitForSelectorOptions{Timeout: &waitTimeOut})
	if err != nil {
		t.Errorf("Wait created workspace result fail: %v", err)
		t.FailNow()
		return
	}
	workspacePage.WaitForTimeout(1000)
	locator, err = workspacePage.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	c, err := locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find created workspace fail: %v", err)
		t.FailNow()
		return
	}
}

func SyncFromZentao(t provider.T) {
	t.ID("5751")
	t.AddParentSuite("管理禅道站点下工作目录")
	var waitTimeOut float64 = 5000
	_, err = workspacePage.WaitForSelector(".tree-node", playwright.PageWaitForSelectorOptions{Timeout: &waitTimeOut})
	if err != nil {
		CreateWorkspace(t)
		SyncFromZentao(t)
		return
	}
	locator, err := workspacePage.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	c, err := locator.Count()
	if err != nil || c == 0 {
		CreateWorkspace(t)
		SyncFromZentao(t)
		return
	}
	locator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	if err != nil {
		t.Errorf("Right click node fail: %v", err)
		t.FailNow()
		return
	}
	workspacePage.Click(".tree-context-menu>>text=从禅道同步")
	if err != nil {
		t.Errorf("Click sync from zentao fail: %v", err)
		t.FailNow()
		return
	}
	_, err = workspacePage.WaitForSelector("#syncFromZentaoFormModal .z-tbody-checkbox")
	if err != nil {
		t.Errorf("Wait syncFromZentaoFormModal fail: %v", err)
		t.FailNow()
		return
	}
	err = workspacePage.Click("#syncFromZentaoFormModal>>.modal-action>>span:has-text(\"确定\")")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		t.FailNow()
		return
	}
	_, err = workspacePage.WaitForSelector("#syncFromZentaoFormModal", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateHidden})
	if err != nil {
		t.Errorf("Wait syncFromZentaoFormModal hide fail: %v", err)
		t.FailNow()
		return
	}
	locator, err = workspacePage.Locator(".toast-notification-container", playwright.PageLocatorOptions{HasText: "成功从禅道同步"})
	c, err = locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Sync from zentao fail: %v", err)
		t.FailNow()
		return
	}
}

func SyncTwoCaseFromZentao(t provider.T) {
	t.ID("5752")
	t.AddParentSuite("管理禅道站点下工作目录")
	var waitTimeOut float64 = 5000
	_, err = workspacePage.WaitForSelector(".tree-node", playwright.PageWaitForSelectorOptions{Timeout: &waitTimeOut})
	if err != nil {
		CreateWorkspace(t)
		SyncFromZentao(t)
		return
	}
	locator, err := workspacePage.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	c, err := locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find workspace fail: %v", err)
		t.FailNow()
		return
	}
	locator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	if err != nil {
		t.Errorf("Right click node fail: %v", err)
		t.FailNow()
		return
	}
	workspacePage.Click(".tree-context-menu>>text=从禅道同步")
	if err != nil {
		t.Errorf("Click sync from zentao fail: %v", err)
		t.FailNow()
		return
	}
	_, err = workspacePage.WaitForSelector("#syncFromZentaoFormModal .z-tbody-checkbox")
	if err != nil {
		t.Errorf("Wait syncFromZentaoFormModal fail: %v", err)
		t.FailNow()
		return
	}
	err = workspacePage.Click("text=编号标题类型状态结果 >> input[type=\"checkbox\"]")
	workspacePage.Click("text=1check string matches pattern功能测试正常 >> input[type=\"checkbox\"]")
	workspacePage.Click("text=2extract content from webpage功能测试 >> input[type=\"checkbox\"]")
	err = workspacePage.Click("#syncFromZentaoFormModal>>.modal-action>>span:has-text(\"确定\")")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		t.FailNow()
		return
	}
	_, err = workspacePage.WaitForSelector("#syncFromZentaoFormModal", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateHidden})
	if err != nil {
		t.Errorf("Wait syncFromZentaoFormModal hide fail: %v", err)
		t.FailNow()
		return
	}
	locator, err = workspacePage.Locator(".toast-notification-container", playwright.PageLocatorOptions{HasText: "成功从禅道同步2个用例"})

	c, err = locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Sync from zentao fail: %v", err)
		t.FailNow()
		return
	}
}

func SyncToZentao(t provider.T) {
	t.ID("5431")
	t.AddParentSuite("管理禅道站点下工作目录")
	var waitTimeOut float64 = 5000
	_, err = workspacePage.WaitForSelector(".tree-node", playwright.PageWaitForSelectorOptions{Timeout: &waitTimeOut})
	if err != nil {
		CreateWorkspace(t)
		SyncFromZentao(t)
		return
	}
	locator, err := workspacePage.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	c, err := locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find workspace fail: %v", err)
		t.FailNow()
		return
	}
	err = locator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	if err != nil {
		t.Errorf("Right click node fail: %v", err)
		t.FailNow()
		return
	}
	workspacePage.Click(".tree-context-menu>>text=同步到禅道")
	if err != nil {
		t.Errorf("Click sync to zentao fail: %v", err)
		t.FailNow()
		return
	}
	_, err = workspacePage.WaitForSelector(".toast-notification-close")
	if err != nil {
		t.Errorf("Wait toast-notification-close fail: %v", err)
		t.FailNow()
		return
	}
	locator, err = workspacePage.Locator(".toast-notification-container", playwright.PageLocatorOptions{HasText: "成功同步"})
	c, err = locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Sync to zentao fail: %v", err)
		t.FailNow()
		return
	}
}
func Copy(t provider.T) {
	t.ID("5474")
	t.AddParentSuite("管理禅道站点下工作目录")
	ztfTestHelper.SelectSite(workspacePage)
	ztfTestHelper.ExpandProduct(workspacePage)
	scriptLocator, err := workspacePage.Locator(".tree-node:has-text('单元测试工作目录')>>.tree-node-title>>text=1_string_match.php")
	if err != nil {
		t.Errorf("Find 1_string_match.php fail: %v", err)
		t.FailNow()
		return
	}
	err = scriptLocator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	if err != nil {
		t.Errorf("Right click script fail: %v", err)
		t.FailNow()
		return
	}
	err = workspacePage.Click(".tree-context-menu>>text=复制")
	if err != nil {
		t.Errorf("Click copy fail: %v", err)
		t.FailNow()
		return
	}
	productLocator, err := workspacePage.Locator(".tree-node:has-text('单元测试工作目录')>>.tree-node-item:has-text('product1')")
	err = productLocator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	if err != nil {
		t.Errorf("Right click workspace fail: %v", err)
		t.FailNow()
		return
	}
	err = workspacePage.Click(".tree-context-menu>>text=粘贴")
	if err != nil {
		t.Errorf("Click parse fail: %v", err)
		t.FailNow()
		return
	}
	workspacePage.WaitForTimeout(1000)
	scriptLocator, err = workspacePage.Locator(".tree-node:has-text('单元测试工作目录')>>.tree-node-title>>text=1_string_match.php")
	c, err := scriptLocator.Count()
	if err != nil || c < 2 {
		t.Errorf("Find 1_string_match fail: %v", err)
		t.FailNow()
		return
	}
}
func DeleteScript(t provider.T) {
	t.ID("5478")
	t.AddParentSuite("管理禅道站点下工作目录")
	locator, err := workspacePage.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	c, err := locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find workspace fail: %v", err)
		t.FailNow()
		return
	}
	scriptLocator, err := workspacePage.Locator(".tree-node-title>>text=1.php")
	if err != nil {
		t.Errorf("Find 1.php fail: %v", err)
		t.FailNow()
		return
	}
	err = scriptLocator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	if err != nil {
		t.Errorf("Right click script fail: %v", err)
		t.FailNow()
		return
	}
	err = workspacePage.Click(".tree-context-menu>>text=删除")
	if err != nil {
		t.Errorf("Click delete fail: %v", err)
		t.FailNow()
		return
	}
	err = workspacePage.Click(".modal-action>>span:has-text(\"确定\")")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		t.FailNow()
		return
	}
	workspacePage.WaitForTimeout(1000)
	scriptLocator, err = workspacePage.Locator(".tree-node-item>>div:has-text('1.php')")
	c, err = scriptLocator.Count()
	if err != nil || c > 0 {
		t.Errorf("Delete script fail: %v", err)
		t.FailNow()
		return
	}
}
func DeleteDir(t provider.T) {
	t.ID("5477")
	t.AddParentSuite("管理禅道站点下工作目录")
	_, err = workspacePage.WaitForSelector(".tree-node")
	if err != nil {
		t.Errorf("Wait tree-node fail: %v", err)
		t.FailNow()
		return
	}
	ztfTestHelper.ExpandWorspace(workspacePage)
	productLocator, err := workspacePage.Locator(".tree-node:has-text('单元测试工作目录')>>.tree-node-item:has-text('product1')")
	if err != nil {
		t.Errorf("Find product1 fail: %v", err)
		t.FailNow()
		return
	}
	err = productLocator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	if err != nil {
		t.Errorf("Right click script fail: %v", err)
		t.FailNow()
		return
	}
	err = workspacePage.Click(".tree-context-menu>>text=删除")
	if err != nil {
		t.Errorf("Click delete fail: %v", err)
		t.FailNow()
		return
	}
	err = workspacePage.Click(".modal-action>>span:has-text(\"确定\")")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		t.FailNow()
		return
	}
	workspacePage.WaitForTimeout(1000)
	scriptLocator, err := workspacePage.Locator(".tree-node:has-text('单元测试工作目录')>>.tree-node-item:has-text('product1')")
	c, err := scriptLocator.Count()
	if err != nil || c > 0 {
		t.Errorf("Delete workspace fail: %v", err)
		t.FailNow()
		return
	}
}

func DeleteWorkspace(t provider.T) {
	t.ID("5468")
	t.AddParentSuite("管理禅道站点下工作目录")
	_, err = workspacePage.WaitForSelector(".tree-node")
	if err != nil {
		t.Errorf("Wait tree-node fail: %v", err)
		t.FailNow()
		return
	}
	locator, err := workspacePage.Locator(".tree-node-item", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	c, err := locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find workspace fail: %v", err)
		t.FailNow()
		return
	}
	err = locator.Hover()
	if err != nil {
		t.Errorf("The hover workspace fail: %v", err)
		t.FailNow()
		return
	}
	err = workspacePage.Click(`[title="删除"]`)
	if err != nil {
		t.Errorf("The click delete fail: %v", err)
		t.FailNow()
		return
	}
	err = workspacePage.Click(".modal-action>>span:has-text(\"确定\")")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		t.FailNow()
		return
	}
	workspacePage.WaitForTimeout(1000)
	scriptLocator, err := workspacePage.Locator(".tree-node-title:has-text('单元测试工作目录')")
	c, err = scriptLocator.Count()
	if err != nil || c > 0 {
		t.Errorf("Delete workspace fail: %v", err)
		t.FailNow()
		return
	}
}
func Clip(t provider.T) {
	t.ID("5476")
	t.AddParentSuite("管理禅道站点下工作目录")
	_, err = workspacePage.WaitForSelector(".tree-node")
	if err != nil {
		t.Errorf("Wait tree-node fail: %v", err)
		t.FailNow()
		return
	}
	locator, err := workspacePage.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	c, err := locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find workspace fail: %v", err)
		t.FailNow()
		return
	}
	err = locator.Click()
	if err != nil {
		t.Errorf("Click node fail: %v", err)
		t.FailNow()
		return
	}
	ztfTestHelper.ExpandProduct(workspacePage)
	scriptLocator, err := locator.Locator(".tree-node-title>>text=1.php")
	if err != nil {
		t.Errorf("Find 1.php fail: %v", err)
		t.FailNow()
		return
	}
	err = scriptLocator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	if err != nil {
		t.Errorf("Right click script fail: %v", err)
		t.FailNow()
		return
	}
	err = workspacePage.Click(".tree-context-menu>>text=剪切")
	if err != nil {
		t.Errorf("Click copy fail: %v", err)
		t.FailNow()
		return
	}
	workspaceLocator, err := workspacePage.Locator(".tree-node-title", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	if err != nil {
		t.Errorf("Find workspace fail: %v", err)
		t.FailNow()
		return
	}
	err = workspaceLocator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	if err != nil {
		t.Errorf("Right click workspace fail: %v", err)
		t.FailNow()
		return
	}
	err = workspacePage.Click(".tree-context-menu>>text=粘贴")
	if err != nil {
		t.Errorf("Click parse fail: %v", err)
		t.FailNow()
		return
	}
	workspacePage.WaitForTimeout(1000)
	scriptLocator, err = locator.Locator(".tree-node-item>>div:has-text('1.php')")
	c, err = scriptLocator.Count()
	if err != nil || c < 1 {
		t.Errorf("Find workspace fail: %v", err)
		t.FailNow()
		return
	}
}

func Collapse(t provider.T) {
	t.ID("5472")
	t.AddParentSuite("管理禅道站点下工作目录")
	locator, err := workspacePage.Locator("#siteMenuToggle")
	if err != nil {
		t.Errorf("The siteMenuToggle is missing: %v", err)
		t.FailNow()
		return
	}
	err = locator.Click()
	if err != nil {
		t.Errorf("The Click is fail: %v", err)
		t.FailNow()
		return
	}
	_, err = workspacePage.WaitForSelector("#navbar .list-item")
	if err != nil {
		t.Errorf("Wait for workspace list nav fail: %v", err)
		t.FailNow()
		return
	}
	err = workspacePage.Click(".list-item-title>>text=单元测试站点")
	if err != nil {
		t.Errorf("The Click workspace nav fail: %v", err)
		t.FailNow()
		return
	}
	className, err := workspacePage.GetAttribute(".tree-node-title:has-text(\"单元测试工作目录\")", "class")
	if err != nil {
		t.Errorf("Find workspace fail: %v", err)
		t.FailNow()
		return
	}
	if strings.Contains(className, "collapsed") {
		err = workspacePage.Click(`#leftPane>>.toolbar>>[title="展开"]`)
	} else {
		err = workspacePage.Click(`#leftPane>>.toolbar>>[title="折叠"]`)
	}
	if err != nil {
		t.Errorf("Click expand workspace btn fail: %v", err)
		t.FailNow()
		return
	}
	workspacePage.WaitForTimeout(1000)
	locator, _ = workspacePage.Locator("#leftPane>>.tree-node-item>>text=1_string_match.php")
	count, _ := locator.Count()
	if strings.Contains(className, "collapsed") && count == 0 {
		t.Error("Expand workspace fail")
		t.FailNow()
		return
	} else if !strings.Contains(className, "collapsed") && count > 0 {

	}
	if strings.Contains(className, "collapsed") {
		err = workspacePage.Click(`#leftPane>>.toolbar>>[title="折叠"]`)
	} else {
		err = workspacePage.Click(`#leftPane>>.toolbar>>[title="展开"]`)
	}
	if err != nil {
		t.Errorf("Click Collapse workspace btn fail: %v", err)
		t.FailNow()
		return
	}
	workspacePage.WaitForTimeout(100)
	locator, _ = workspacePage.Locator("#leftPane>>.tree-node-item>>text=1_string_match.php")
	count, _ = locator.Count()
	if strings.Contains(className, "collapsed") && count == 0 {
		t.Error("Expand workspace fail")
		t.FailNow()
		return
	} else if !strings.Contains(className, "collapsed") && count > 0 {

	}
}
func TestUiWorkspace(t *testing.T) {
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
		return
	}
	headless := true
	var slowMo float64 = 100
	workspaceBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web workspaceBrowser: %v", err)
		t.FailNow()
		return
	}
	workspacePage, err = workspaceBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
		return
	}
	if _, err = workspacePage.Goto("http://127.0.0.1:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
		t.FailNow()
		return
	}
	ztfTestHelper.SelectSite(workspacePage)

	runner.Run(t, "客户端-同步到禅道", SyncToZentao)
	runner.Run(t, "客户端-从禅道同步选中用例", SyncTwoCaseFromZentao)
	runner.Run(t, "客户端-从禅道同步", SyncFromZentao)
	runner.Run(t, "客户端-复制粘贴树状脚本文件", Copy)
	runner.Run(t, "客户端-剪切粘贴树状脚本文件", Clip)
	runner.Run(t, "客户端-删除树状脚本文件", DeleteScript)
	runner.Run(t, "客户端-删除树状脚本文件夹", DeleteDir)
	runner.Run(t, "客户端-显示展开折叠脚本树状结构", Collapse)
	runner.Run(t, "客户端-删除禅道工作目录", DeleteWorkspace)
	runner.Run(t, "客户端-创建禅道工作目录", CreateWorkspace)
	if err = workspaceBrowser.Close(); err != nil {
		t.Errorf("The workspaceBrowser cannot be closed: %v", err)
		t.FailNow()
		return
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
		return
	}
}
