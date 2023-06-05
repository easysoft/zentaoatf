package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	commonTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/common"
	constTestHelper "github.com/easysoft/zentaoatf/test/helper/conf"
	apiTest "github.com/easysoft/zentaoatf/test/helper/zentao/api"
	ztfTestHelper "github.com/easysoft/zentaoatf/test/helper/ztf"
	plwConf "github.com/easysoft/zentaoatf/cmd/test/ui/conf"
	plwHelper "github.com/easysoft/zentaoatf/test/ui/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	playwright "github.com/playwright-community/playwright-go"
)

var (
	workspacePath = fmt.Sprintf("%stest%sdemo%sphp", constTestHelper.RootPath, constTestHelper.FilePthSep, constTestHelper.FilePthSep)
)

func CreateWorkspace(t provider.T) {
	t.ID("5468")
	t.AddParentSuite("管理禅道站点下工作目录")

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
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

	webpage.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
}

func SyncFromZentao(t provider.T) {
	syncCaseFromZentaoTask(t)
	syncCaseFromZentaoModule(t)
	syncCaseFromZentaoSuite(t)
	syncAllCaseFromZentao(t)
}
func syncAllCaseFromZentao(t provider.T) {
	syncDir := filepath.Join(workspacePath, "product1")
	os.RemoveAll(syncDir)

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	ztfTestHelper.SelectSite(webpage)
	ztfTestHelper.ExpandWorspace(webpage)

	locator := webpage.Locator(".tree-node-title:has-text('单元测试工作目录')")
	plwConf.EnableErr()
	locator.RightClick()

	webpage.Click(".tree-context-menu>>text=从禅道同步")

	webpage.WaitForSelector("#syncFromZentaoFormModal .z-tbody-checkbox")
	webpage.Click("#syncFromZentaoFormModal>>.modal-action>>span:has-text(\"确定\")")

	webpage.WaitForSelector("#syncFromZentaoFormModal", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateHidden})
	webpage.Locator(".toast-notification-container", playwright.PageLocatorOptions{HasText: "成功从禅道同步"})

	//check file info
	scriptPath := filepath.Join(workspacePath, "product1", "2.php")
	content := fileUtils.ReadFile(scriptPath)
	t.Require().Contains(content, "extract content from webpage-synctozentao")
}

func syncCaseFromZentaoTask(t provider.T) {
	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	syncDir := filepath.Join(workspacePath, "product1")
	os.RemoveAll(syncDir)

	ztfTestHelper.SelectSite(webpage)

	webpage.WaitForSelector(".tree-node", playwright.PageWaitForSelectorOptions{Timeout: &plwConf.Timeout})
	locator := webpage.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	locator.RightClick()

	webpage.Click(".tree-context-menu>>text=从禅道同步")
	locator = webpage.Locator("#syncFromZentaoFormModal select")
	locator.SelectNth(2, playwright.SelectOptionValues{Labels: &[]string{"企业网站第一期测试任务"}})

	webpage.WaitForSelector("#syncFromZentaoFormModal .z-tbody-checkbox")
	webpage.Click("text=编号标题类型状态结果 >> input[type=\"checkbox\"]")
	webpage.Click("text=1check string matches pattern功能测试 >> input[type=\"checkbox\"]")
	webpage.Click("text=2extract content from webpage-synctozentao功能测试 >> input[type=\"checkbox\"]")
	webpage.Click("#syncFromZentaoFormModal>>.modal-action>>span:has-text(\"确定\")")

	webpage.WaitForSelector("#syncFromZentaoFormModal", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateHidden})
	locator = webpage.Locator(".toast-notification-container", playwright.PageLocatorOptions{HasText: "成功从禅道同步2个用例"})

	//check file info
	scriptPath := filepath.Join(workspacePath, "product1", "2.php")
	content := fileUtils.ReadFile(scriptPath)
	t.Require().Contains(content, "extract content from webpage-synctozentao")
}

func syncCaseFromZentaoModule(t provider.T) {
	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	syncDir := filepath.Join(workspacePath, "product1")
	os.RemoveAll(syncDir)

	ztfTestHelper.SelectSite(webpage)

	webpage.WaitForSelector(".tree-node", playwright.PageWaitForSelectorOptions{Timeout: &plwConf.Timeout})
	locator := webpage.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	locator.RightClick()

	webpage.Click(".tree-context-menu>>text=从禅道同步")
	locator = webpage.Locator("#syncFromZentaoFormModal select")
	locator.SelectNth(0, playwright.SelectOptionValues{Labels: &[]string{"/module1"}})

	webpage.WaitForSelector("#syncFromZentaoFormModal .z-tbody-checkbox")
	webpage.Click("text=编号标题类型状态结果 >> input[type=\"checkbox\"]")
	webpage.Click("text=6module1-case2功能测试 >> input[type=\"checkbox\"]")
	webpage.Click("text=5module1-case1功能测试 >> input[type=\"checkbox\"]")
	webpage.Click("#syncFromZentaoFormModal>>.modal-action>>span:has-text(\"确定\")")

	webpage.WaitForSelector("#syncFromZentaoFormModal", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateHidden})
	locator = webpage.Locator(".toast-notification-container", playwright.PageLocatorOptions{HasText: "成功从禅道同步2个用例"})

	//check file info
	scriptPath := filepath.Join(workspacePath, "product1", "6.php")
	content := fileUtils.ReadFile(scriptPath)
	t.Require().Contains(content, "module1-case2")
}

func syncCaseFromZentaoSuite(t provider.T) {
	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	syncDir := filepath.Join(workspacePath, "product1")
	os.RemoveAll(syncDir)

	ztfTestHelper.SelectSite(webpage)

	webpage.WaitForSelector(".tree-node", playwright.PageWaitForSelectorOptions{Timeout: &plwConf.Timeout})
	locator := webpage.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	locator.RightClick()

	webpage.Click(".tree-context-menu>>text=从禅道同步")
	locator = webpage.Locator("#syncFromZentaoFormModal select")
	locator.SelectNth(1, playwright.SelectOptionValues{Labels: &[]string{"test_suite"}})

	webpage.WaitForSelector("#syncFromZentaoFormModal .z-tbody-checkbox")
	webpage.Click("text=编号标题类型状态结果 >> input[type=\"checkbox\"]")
	webpage.Click("text=1check string matches pattern功能测试 >> input[type=\"checkbox\"]")
	webpage.Click("#syncFromZentaoFormModal>>.modal-action>>span:has-text(\"确定\")")

	webpage.WaitForSelector("#syncFromZentaoFormModal", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateHidden})
	locator = webpage.Locator(".toast-notification-container", playwright.PageLocatorOptions{HasText: "成功从禅道同步1个用例"})

	//check file info
	scriptPath := filepath.Join(workspacePath, "product1", "1.php")
	content := fileUtils.ReadFile(scriptPath)
	t.Require().Contains(content, "check string matches pattern")
}

func SyncTwoCaseFromZentao(t provider.T) {
	t.ID("5752")
	t.AddParentSuite("管理禅道站点下工作目录")

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	syncDir := filepath.Join(workspacePath, "product1")
	os.RemoveAll(syncDir)

	ztfTestHelper.SelectSite(webpage)

	webpage.WaitForSelector(".tree-node", playwright.PageWaitForSelectorOptions{Timeout: &plwConf.Timeout})
	locator := webpage.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	locator.RightClick()

	webpage.Click(".tree-context-menu>>text=从禅道同步")

	webpage.WaitForSelector("#syncFromZentaoFormModal .z-tbody-checkbox")
	webpage.Click("text=编号标题类型状态结果 >> input[type=\"checkbox\"]")
	webpage.Click("text=1check string matches pattern功能测试 >> input[type=\"checkbox\"]")
	webpage.Click("text=2extract content from webpage-synctozentao功能测试 >> input[type=\"checkbox\"]")
	webpage.Click("#syncFromZentaoFormModal>>.modal-action>>span:has-text(\"确定\")")

	webpage.WaitForSelector("#syncFromZentaoFormModal", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateHidden})
	locator = webpage.Locator(".toast-notification-container", playwright.PageLocatorOptions{HasText: "成功从禅道同步2个用例"})

	//check file info
	scriptPath := filepath.Join(workspacePath, "product1", "2.php")
	content := fileUtils.ReadFile(scriptPath)
	t.Require().Contains(content, "extract content from webpage-synctozentao")
}

func SyncToZentao(t provider.T) {
	t.ID("5431")
	t.AddParentSuite("管理禅道站点下工作目录")

	//update script title
	scriptPath := filepath.Join(workspacePath, "2_webpage_extract.php")
	content := fileUtils.ReadFile(scriptPath)
	newContent := strings.Replace(content, "title=extract content from webpage", "title=extract content from webpage-synctozentao", 1)
	fileUtils.WriteFile(scriptPath, newContent)
	defer fileUtils.WriteFile(scriptPath, content)

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	ztfTestHelper.SelectSite(webpage)

	webpage.WaitForSelector(".tree-node", playwright.PageWaitForSelectorOptions{Timeout: &plwConf.Timeout})
	locator := webpage.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	locator.RightClick()

	webpage.Click(".tree-context-menu>>text=同步到禅道")

	webpage.WaitForSelector(".toast-notification-close")
	locator = webpage.Locator(".toast-notification-container", playwright.PageLocatorOptions{HasText: "成功同步"})

	//check zentao info
	title := apiTest.GetCaseTitleById(2)
	t.Require().Equal("extract content from webpage-synctozentao", title)
}

func Copy(t provider.T) {
	t.ID("5474")
	t.AddParentSuite("管理禅道站点下工作目录")

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	ztfTestHelper.SelectSite(webpage)
	ztfTestHelper.ExpandProduct(webpage)

	scriptLocator := webpage.Locator(".tree-node:has-text('单元测试工作目录')>>.tree-node-title>>text=1_string_match.php")
	scriptLocator.RightClick()

	webpage.Click(".tree-context-menu>>text=复制")
	productLocator := webpage.Locator(".tree-node:has-text('单元测试工作目录')>>.tree-node-item:has-text('product1')")
	productLocator.RightClick()
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

func ClipAndCopyDir(t provider.T) {
	CopyDir(t)
	ClipDir(t)
}

func CopyDir(t provider.T) {
	t.ID("7593")
	t.AddParentSuite("管理禅道站点下工作目录")

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	ztfTestHelper.SelectSite(webpage)
	ztfTestHelper.ExpandProduct(webpage)
	CreateDir(t)

	scriptLocator := webpage.Locator(".tree-node:has-text('单元测试工作目录')>>.tree-node-title>>text=oldDir")
	scriptLocator.RightClick()

	webpage.Click(".tree-context-menu>>text=复制")
	productLocator := webpage.Locator(".tree-node:has-text('单元测试工作目录')>>.tree-node-item:has-text('product1')")
	productLocator.RightClick()
	webpage.Click(".tree-context-menu>>text=粘贴")

	if webpage.ElementExist(".tree-node:has-text('单元测试工作目录')>>.tree-node-item:has-text('product1')>>.tree-node-title>>text=oldDir") {
		t.Errorf("Copy dir fail")
		t.FailNow()
		return
	}

	os.Remove(commonTestHelper.GetPhpWorkspacePath() + "oldDir")
	os.Remove(filepath.Join(commonTestHelper.GetPhpWorkspacePath(), "product1", "oldDir"))
}

func ClipDir(t provider.T) {
	t.ID("7593")
	t.AddParentSuite("管理禅道站点下工作目录")

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	ztfTestHelper.SelectSite(webpage)
	ztfTestHelper.ExpandProduct(webpage)
	CreateDir(t)

	scriptLocator := webpage.Locator(".tree-node:has-text('单元测试工作目录')>>.tree-node-title>>text=oldDir")
	scriptLocator.RightClick()

	webpage.Click(".tree-context-menu>>text=剪切")
	productLocator := webpage.Locator(".tree-node:has-text('单元测试工作目录')>>.tree-node-item:has-text('product1')")
	productLocator.RightClick()
	webpage.Click(".tree-context-menu>>text=粘贴")

	if webpage.ElementExist(".tree-node:has-text('单元测试工作目录')>>.tree-node-item:has-text('product1')>>.tree-node-title>>text=oldDir") {
		t.Errorf("Copy dir fail")
		t.FailNow()
		return
	}

	os.Remove(commonTestHelper.GetPhpWorkspacePath() + "oldDir")
	os.Remove(filepath.Join(commonTestHelper.GetPhpWorkspacePath(), "product1", "oldDir"))
}

func CreateScript(t provider.T) {
	// t.ID("7596")
	// t.AddParentSuite("管理禅道站点下工作目录")
	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	ztfTestHelper.SelectSite(webpage)
	ztfTestHelper.ExpandProduct(webpage)

	if webpage.ElementExist(".tree-node:has-text('单元测试工作目录')>>.tree-node-title>>text=old.php") {
		return
	}

	workspaceLocator := webpage.Locator(".tree-node-title:has-text('单元测试工作目录')")
	workspaceLocator.Hover()
	webpage.Click(`[title="新建脚本"]`)
	productLocator := webpage.Locator("#scriptFormModal>>input")
	productLocator.Fill("old.php")
	webpage.Click("#scriptFormModal>>.modal-action>>span:has-text(\"确定\")")

	webpage.Locator(".tree-node:has-text('单元测试工作目录')>>.tree-node-title>>text=old.php")
}

func RenameScript(t provider.T) {
	t.ID("7596")
	t.AddParentSuite("管理禅道站点下工作目录")

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	ztfTestHelper.SelectSite(webpage)
	ztfTestHelper.ExpandProduct(webpage)
	CreateScript(t)

	scriptLocator := webpage.Locator(".tree-node:has-text('单元测试工作目录')>>.tree-node-title>>text=old.php")
	scriptLocator.RightClick()

	webpage.Click(".tree-context-menu>>text=重命名")
	productLocator := webpage.Locator("#scriptFormModal>>input")
	productLocator.Fill("new.php")
	webpage.Click("#scriptFormModal>>.modal-action>>span:has-text(\"确定\")")

	scriptLocator = webpage.Locator(".tree-node:has-text('单元测试工作目录')>>.tree-node-title>>text=new.php")

	os.Remove(commonTestHelper.GetPhpWorkspacePath() + "old.php")
	os.Remove(commonTestHelper.GetPhpWorkspacePath() + "new.php")
}

func CreateDir(t provider.T) {
	// t.ID("7596")
	// t.AddParentSuite("管理禅道站点下工作目录")
	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	ztfTestHelper.SelectSite(webpage)
	ztfTestHelper.ExpandProduct(webpage)

	if webpage.ElementExist(".tree-node:has-text('单元测试工作目录')>>.tree-node-title>>text=oldDir") {
		return
	}

	workspaceLocator := webpage.Locator(".tree-node-title:has-text('单元测试工作目录')")
	workspaceLocator.Hover()
	webpage.Click(`.tree-node-item:has-text('单元测试工作目录')>>[title="新建工作目录"]`)
	productLocator := webpage.Locator("#scriptFormModal>>input")
	productLocator.Fill("oldDir")
	webpage.Click("#scriptFormModal>>.modal-action>>span:has-text(\"确定\")")

	webpage.Locator(".tree-node:has-text('单元测试工作目录')>>.tree-node-title>>text=oldDir")
}

func RenameDir(t provider.T) {
	t.ID("7595")
	t.AddParentSuite("管理禅道站点下工作目录")

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	ztfTestHelper.SelectSite(webpage)
	ztfTestHelper.ExpandProduct(webpage)
	CreateDir(t)

	dirLocator := webpage.Locator(".tree-node:has-text('单元测试工作目录')>>.tree-node-title>>text=oldDir")
	dirLocator.RightClick()
	webpage.Click(".tree-context-menu>>text=重命名")
	productLocator := webpage.Locator("#scriptFormModal>>input")
	productLocator.Fill("newDir")
	webpage.Click("#scriptFormModal>>.modal-action>>span:has-text(\"确定\")")

	webpage.Locator(".tree-node:has-text('单元测试工作目录')>>.tree-node-title>>text=newDir")

	os.Remove(commonTestHelper.GetPhpWorkspacePath() + "oldDir")
	os.Remove(commonTestHelper.GetPhpWorkspacePath() + "newDir")
}

func DeleteScript(t provider.T) {
	t.ID("5478")
	t.AddParentSuite("管理禅道站点下工作目录")

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	ztfTestHelper.SelectSite(webpage)
	ztfTestHelper.ExpandProduct(webpage)

	scriptLocator := webpage.Locator(".tree-node-title>>text=1.php")
	scriptLocator.RightClick()
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

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	ztfTestHelper.SelectSite(webpage)
	ztfTestHelper.ExpandWorspace(webpage)

	productLocator := webpage.Locator(".tree-node:has-text('单元测试工作目录')>>.tree-node-item:has-text('product1')")
	productLocator.RightClick()
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

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
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

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	ztfTestHelper.SelectSite(webpage)

	webpage.WaitForSelector(".tree-node")
	locator := webpage.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	locator.Click()
	ztfTestHelper.ExpandProduct(webpage)

	scriptLocator := locator.Locator(".tree-node-title>>text=1.php")
	scriptLocator.RightClick()
	webpage.Click(".tree-context-menu>>text=剪切")
	workspaceLocator := webpage.Locator(".tree-node-title", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	workspaceLocator.RightClick()
	webpage.Click(".tree-context-menu>>text=粘贴")

	webpage.WaitForTimeout(1000)
	locator.Locator(".tree-node-item>>div:has-text('1.php')")
}

func Collapse(t provider.T) {
	t.ID("5472")
	t.AddParentSuite("管理禅道站点下工作目录")

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	ztfTestHelper.SelectSite(webpage)
	ztfTestHelper.ExpandWorspace(webpage)

	webpage.WaitForSelectorTimeout(".tree-node:has-text(\"单元测试工作目录\")", 5000)
	className := webpage.GetAttribute(".tree-node:has-text(\"单元测试工作目录\")", "class")

	if strings.Contains(className, "collapsed") {
		webpage.Click(`#leftPane>>.toolbar>>[title="展开"]`)
	} else {
		webpage.Click(`#leftPane>>.toolbar>>[title="折叠"]`)
	}

	if strings.Contains(className, "collapsed") {
		webpage.WaitForSelectorTimeout("#leftPane>>.tree-node-item>>text=1_string_match.php", 5000)
	} else if !strings.Contains(className, "collapsed") {
		webpage.WaitForSelectorTimeout("#leftPane>>.tree-node-item>>text=1_string_match.php", 5000, playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})
	}

	if strings.Contains(className, "collapsed") {
		webpage.Click(`#leftPane>>.toolbar>>[title="折叠"]`)
	} else {
		webpage.Click(`#leftPane>>.toolbar>>[title="展开"]`)
	}

	if strings.Contains(className, "collapsed") {
		webpage.WaitForSelectorTimeout("#leftPane>>.tree-node-item>>text=1_string_match.php", 5000, playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})
	} else if !strings.Contains(className, "collapsed") {
		webpage.WaitForSelectorTimeout("#leftPane>>.tree-node-item>>text=1_string_match.php", 5000)
	}
}
func TestUiWorkspace(t *testing.T) {
	runner.Run(t, "客户端-同步到禅道", SyncToZentao)
	runner.Run(t, "客户端-从禅道同步选中用例", SyncTwoCaseFromZentao)
	runner.Run(t, "客户端-从禅道同步", SyncFromZentao)
	runner.Run(t, "客户端-复制粘贴树状脚本文件", Copy)
	runner.Run(t, "客户端-复制粘贴目录", ClipAndCopyDir)
	runner.Run(t, "客户端-剪切粘贴树状脚本文件", Clip)
	runner.Run(t, "客户端-重命名脚本文件", RenameScript)
	runner.Run(t, "客户端-重命名目录", RenameDir)
	runner.Run(t, "客户端-删除树状脚本文件", DeleteScript)
	runner.Run(t, "客户端-删除树状脚本文件夹", DeleteDir)
	runner.Run(t, "客户端-显示展开折叠脚本树状结构", Collapse)
	runner.Run(t, "客户端-删除禅道工作目录", DeleteWorkspace)
	runner.Run(t, "客户端-创建禅道工作目录", CreateWorkspace)
}
