package main

import (
	"fmt"
	"strings"
	"testing"

	commonTestHelper "github.com/easysoft/zentaoatf/test/helper/common"
	ztfTestHelper "github.com/easysoft/zentaoatf/test/helper/ztf"
	plwConf "github.com/easysoft/zentaoatf/test/ui/conf"
	plwHelper "github.com/easysoft/zentaoatf/test/ui/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	playwright "github.com/playwright-community/playwright-go"
)

var (
	workspacePath = fmt.Sprintf("%stest%sdemo%sphp", commonTestHelper.RootPath, commonTestHelper.FilePthSep, commonTestHelper.FilePthSep)
)

func CreateWorkspace(t provider.T) {
	t.ID("5468")
	t.AddParentSuite("管理禅道站点下工作目录")
	webpage, _ := plwHelper.OpenUrl("http://127.0.0.1:8000/", t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage)
	webpage.Click(`[title="新建工作目录"]`)
	webpage.WaitForSelector("#workspaceFormModal")
	locator := webpage.Locator("#workspaceFormModal input")
	locator.FillNth(0, "单元测试工作目录")
	locator.FillNth(1, workspacePath)
	locator = webpage.Locator("#workspaceFormModal select")
	locator.SelectNth(0, playwright.SelectOptionValues{Values: &[]string{"ztf"}})
	locator.SelectNth(1, playwright.SelectOptionValues{Values: &[]string{"php"}})
	webpage.Click("#workspaceFormModal>>.modal-action>>span:has-text(\"确定\")")
	var waitTimeOut float64 = 5000
	webpage.WaitForSelector(".tree-node", playwright.PageWaitForSelectorOptions{Timeout: &waitTimeOut})
	webpage.WaitForTimeout(1000)
	webpage.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
}

func SyncFromZentao(t provider.T) {
	t.ID("5751")
	t.AddParentSuite("管理禅道站点下工作目录")
	webpage, _ := plwHelper.OpenUrl("http://127.0.0.1:8000/", t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage)
	ztfTestHelper.ExpandWorspace(webpage)
	locator := webpage.Locator(".tree-node-root>>has-text('单元测试工作目录')")
	plwConf.EnableErr()
	locator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	webpage.Click(".tree-context-menu>>text=从禅道同步")
	webpage.WaitForSelector("#syncFromZentaoFormModal .z-tbody-checkbox")
	webpage.Click("#syncFromZentaoFormModal>>.modal-action>>span:has-text(\"确定\")")
	webpage.WaitForSelector("#syncFromZentaoFormModal", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateHidden})
	locator = webpage.Locator(".toast-notification-container", playwright.PageLocatorOptions{HasText: "成功从禅道同步"})
}

func SyncTwoCaseFromZentao(t provider.T) {
	t.ID("5752")
	t.AddParentSuite("管理禅道站点下工作目录")
	webpage, _ := plwHelper.OpenUrl("http://127.0.0.1:8000/", t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage)
	var waitTimeOut float64 = 5000
	webpage.WaitForSelector(".tree-node", playwright.PageWaitForSelectorOptions{Timeout: &waitTimeOut})
	locator := webpage.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	locator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	webpage.Click(".tree-context-menu>>text=从禅道同步")
	webpage.WaitForSelector("#syncFromZentaoFormModal .z-tbody-checkbox")
	webpage.Click("text=编号标题类型状态结果 >> input[type=\"checkbox\"]")
	webpage.Click("text=1check string matches pattern功能测试正常 >> input[type=\"checkbox\"]")
	webpage.Click("text=2extract content from webpage功能测试 >> input[type=\"checkbox\"]")
	webpage.Click("#syncFromZentaoFormModal>>.modal-action>>span:has-text(\"确定\")")
	webpage.WaitForSelector("#syncFromZentaoFormModal", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateHidden})
	locator = webpage.Locator(".toast-notification-container", playwright.PageLocatorOptions{HasText: "成功从禅道同步2个用例"})
}

func SyncToZentao(t provider.T) {
	t.ID("5431")
	t.AddParentSuite("管理禅道站点下工作目录")
	webpage, _ := plwHelper.OpenUrl("http://127.0.0.1:8000/", t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage)
	var waitTimeOut float64 = 5000
	webpage.WaitForSelector(".tree-node", playwright.PageWaitForSelectorOptions{Timeout: &waitTimeOut})
	locator := webpage.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	locator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	webpage.Click(".tree-context-menu>>text=同步到禅道")
	webpage.WaitForSelector(".toast-notification-close")
	locator = webpage.Locator(".toast-notification-container", playwright.PageLocatorOptions{HasText: "成功同步"})
}
func Copy(t provider.T) {
	t.ID("5474")
	t.AddParentSuite("管理禅道站点下工作目录")
	webpage, _ := plwHelper.OpenUrl("http://127.0.0.1:8000/", t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage)
	ztfTestHelper.ExpandProduct(webpage)
	scriptLocator := webpage.Locator(".tree-node:has-text('单元测试工作目录')>>.tree-node-title>>text=1_string_match.php")
	scriptLocator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	webpage.Click(".tree-context-menu>>text=复制")
	productLocator := webpage.Locator(".tree-node:has-text('单元测试工作目录')>>.tree-node-item:has-text('product1')")
	productLocator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	webpage.Click(".tree-context-menu>>text=粘贴")
	webpage.WaitForTimeout(1000)
	plwConf.DisableErr()
	defer plwConf.EnableErr()
	scriptLocator = webpage.Locator(".tree-node:has-text('单元测试工作目录')>>.tree-node-title>>text=1_string_match.php")
	c := scriptLocator.Count()
	if c < 2 {
		t.Errorf("Find 1_string_match fail")
		t.FailNow()
		return
	}
}
func DeleteScript(t provider.T) {
	t.ID("5478")
	t.AddParentSuite("管理禅道站点下工作目录")
	webpage, _ := plwHelper.OpenUrl("http://127.0.0.1:8000/", t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage)
	ztfTestHelper.ExpandProduct(webpage)
	scriptLocator := webpage.Locator(".tree-node-title>>text=1.php")
	scriptLocator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	webpage.Click(".tree-context-menu>>text=删除")
	webpage.WaitForSelector(".modal-action>>span:has-text(\"确定\")")
	webpage.Click(".modal-action>>span:has-text(\"确定\")")
	webpage.WaitForSelector(".tree-node-title>>text=1.php", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})
	plwConf.DisableErr()
	defer plwConf.EnableErr()
	scriptLocator = webpage.Locator(".tree-node-item>>div:has-text('1.php')")
	c := scriptLocator.Count()
	if c > 0 {
		t.Errorf("Delete script fail")
		t.FailNow()
		return
	}
}
func DeleteDir(t provider.T) {
	t.ID("5477")
	t.AddParentSuite("管理禅道站点下工作目录")
	webpage, _ := plwHelper.OpenUrl("http://127.0.0.1:8000/", t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage)
	ztfTestHelper.ExpandWorspace(webpage)
	productLocator := webpage.Locator(".tree-node:has-text('单元测试工作目录')>>.tree-node-item:has-text('product1')")
	productLocator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	webpage.Click(".tree-context-menu>>text=删除")
	webpage.Click(".modal-action>>span:has-text(\"确定\")")
	webpage.WaitForSelector(".tree-node:has-text('单元测试工作目录')>>.tree-node-item:has-text('product1')", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})
	plwConf.DisableErr()
	defer plwConf.EnableErr()
	scriptLocator := webpage.Locator(".tree-node:has-text('单元测试工作目录')>>.tree-node-item:has-text('product1')")
	c := scriptLocator.Count()
	if c > 0 {
		t.Errorf("Delete workspace fail")
		t.FailNow()
		return
	}
}

func DeleteWorkspace(t provider.T) {
	t.ID("5468")
	t.AddParentSuite("管理禅道站点下工作目录")
	webpage, _ := plwHelper.OpenUrl("http://127.0.0.1:8000/", t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage)
	webpage.WaitForSelector(".tree-node")
	locator := webpage.Locator(".tree-node-item", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	locator.Hover()
	webpage.Click(`[title="删除"]`)
	webpage.Click(".modal-action>>span:has-text(\"确定\")")
	webpage.WaitForTimeout(1000)
	plwConf.DisableErr()
	defer plwConf.EnableErr()
	scriptLocator := webpage.Locator(".tree-node-title:has-text('单元测试工作目录')")
	c := scriptLocator.Count()
	if c > 0 {
		t.Errorf("Delete workspace fail")
		t.FailNow()
		return
	}
}
func Clip(t provider.T) {
	t.ID("5476")
	t.AddParentSuite("管理禅道站点下工作目录")
	webpage, _ := plwHelper.OpenUrl("http://127.0.0.1:8000/", t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage)
	webpage.WaitForSelector(".tree-node")
	locator := webpage.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	locator.Click()
	ztfTestHelper.ExpandProduct(webpage)
	scriptLocator := locator.Locator(".tree-node-title>>text=1.php")
	scriptLocator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	webpage.Click(".tree-context-menu>>text=剪切")
	workspaceLocator := webpage.Locator(".tree-node-title", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	workspaceLocator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	webpage.Click(".tree-context-menu>>text=粘贴")
	webpage.WaitForTimeout(1000)
	locator.Locator(".tree-node-item>>div:has-text('1.php')")
}

func Collapse(t provider.T) {
	t.ID("5472")
	t.AddParentSuite("管理禅道站点下工作目录")
	webpage, _ := plwHelper.OpenUrl("http://127.0.0.1:8000/", t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage)
	ztfTestHelper.ExpandWorspace(webpage)
	className := webpage.GetAttribute(".tree-node:has-text(\"单元测试工作目录\")", "class")
	if strings.Contains(className, "collapsed") {
		webpage.Click(`#leftPane>>.toolbar>>[title="展开"]`)
	} else {
		webpage.Click(`#leftPane>>.toolbar>>[title="折叠"]`)
	}
	if strings.Contains(className, "collapsed") {
		webpage.WaitForSelector("#leftPane>>.tree-node-item>>text=1_string_match.php")
	} else if !strings.Contains(className, "collapsed") {
		webpage.WaitForSelector("#leftPane>>.tree-node-item>>text=1_string_match.php", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})
	}
	if strings.Contains(className, "collapsed") {
		webpage.Click(`#leftPane>>.toolbar>>[title="折叠"]`)
	} else {
		webpage.Click(`#leftPane>>.toolbar>>[title="展开"]`)
	}
	if strings.Contains(className, "collapsed") {
		webpage.WaitForSelector("#leftPane>>.tree-node-item>>text=1_string_match.php", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})
	} else if !strings.Contains(className, "collapsed") {
		webpage.WaitForSelector("#leftPane>>.tree-node-item>>text=1_string_match.php")
	}
}
func TestUiWorkspace(t *testing.T) {
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
}
