package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	commonTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/common"
	constTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/conf"
	apiTest "github.com/easysoft/zentaoatf/cmd/test/helper/zentao/api"
	ztfTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/ztf"
	plwConf "github.com/easysoft/zentaoatf/cmd/test/ui/conf"
	plwHelper "github.com/easysoft/zentaoatf/cmd/test/ui/helper"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	playwright "github.com/playwright-community/playwright-go"
)

var (
	workspacePath = fileUtils.AddFilePathSepIfNeeded(commonTestHelper.GetPhpWorkspacePath())

	syncDir    = filepath.Join(workspacePath, "product1")
	oldDirPath = workspacePath + "oldDir"
	newDirPath = filepath.Join(workspacePath, "product1", "oldDir")
)

func MangeWorkspace(t provider.T) {
	t.ID("5468")
	commonTestHelper.ReplaceLabel(t, "客户端-管理禅道站点下工作目录")

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage, "")

	if webpage.ElementExist(fmt.Sprintf(".tree-node-title>>text=%s", constTestHelper.WorkspaceName)) {
		DeleteWorkspace(t, webpage)
	}

	CreateWorkspace(t, webpage)
	DeleteWorkspace(t, webpage)
}

func MangeLocalWorkspace(t provider.T) {
	t.ID("5467")
	commonTestHelper.ReplaceLabel(t, "客户端-管理本地站点下工作目录")

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage, "本地")

	if webpage.ElementExist(fmt.Sprintf(".tree-node-title>>text=%s", constTestHelper.WorkspaceName)) {
		DeleteWorkspace(t, webpage)
	}

	CreateWorkspace(t, webpage)
	DeleteWorkspace(t, webpage)
}

func CreateWorkspace(t provider.T, webpage plwHelper.Webpage) {
	webpage.Click(`[title="新建工作目录"]`)
	webpage.WaitForSelector("#workspaceFormModal")

	locator := webpage.Locator("#workspaceFormModal input")
	locator.FillNth(0, constTestHelper.WorkspaceName)
	locator.FillNth(1, workspacePath)
	locator = webpage.Locator("#workspaceFormModal select")
	locator.SelectNth(0, playwright.SelectOptionValues{Values: &[]string{"ztf"}})
	locator.SelectNth(1, playwright.SelectOptionValues{Values: &[]string{"php"}})
	webpage.Click("#workspaceFormModal>>.modal-action>>span:has-text(\"确定\")")
	var waitTimeOut float64 = 5000

	webpage.WaitForSelector(".tree-node", playwright.PageWaitForSelectorOptions{Timeout: &waitTimeOut})

	webpage.Locator(".tree-node", playwright.PageLocatorOptions{HasText: constTestHelper.WorkspaceName})
}

func SyncFromZentao(t provider.T) {
	syncCaseFromZentaoTask(t)
	syncCaseFromZentaoModule(t)
	syncCaseFromZentaoSuite(t)
	syncAllCaseFromZentao(t)
}
func syncAllCaseFromZentao(t provider.T) {
	os.RemoveAll(syncDir)

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	ztfTestHelper.SelectSite(webpage, "")
	ztfTestHelper.ExpandWorspace(webpage)

	locator := webpage.Locator(fmt.Sprintf(".tree-node-title:has-text('%s')", constTestHelper.WorkspaceName))
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
	t.Require().Contains(content, "extract content from webpage")
}

func syncCaseFromZentaoTask(t provider.T) {
	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	os.RemoveAll(syncDir)

	ztfTestHelper.SelectSite(webpage, "")

	webpage.WaitForSelector(".tree-node", playwright.PageWaitForSelectorOptions{Timeout: &plwConf.Timeout})
	locator := webpage.Locator(".tree-node-title", playwright.PageLocatorOptions{HasText: constTestHelper.WorkspaceName})
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

	os.RemoveAll(syncDir)

	ztfTestHelper.SelectSite(webpage, "")

	webpage.WaitForSelector(".tree-node", playwright.PageWaitForSelectorOptions{Timeout: &plwConf.Timeout})
	locator := webpage.Locator(".tree-node-title", playwright.PageLocatorOptions{HasText: constTestHelper.WorkspaceName})
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

	os.RemoveAll(syncDir)

	ztfTestHelper.SelectSite(webpage, "")

	webpage.WaitForSelector(".tree-node", playwright.PageWaitForSelectorOptions{Timeout: &plwConf.Timeout})
	locator := webpage.Locator(".tree-node-title", playwright.PageLocatorOptions{HasText: constTestHelper.WorkspaceName})
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
	commonTestHelper.ReplaceLabel(t, "客户端-管理禅道站点下工作目录")

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	os.RemoveAll(syncDir)

	ztfTestHelper.SelectSite(webpage, "")

	webpage.WaitForSelector(".tree-node", playwright.PageWaitForSelectorOptions{Timeout: &plwConf.Timeout})
	locator := webpage.Locator(".tree-node-title", playwright.PageLocatorOptions{HasText: constTestHelper.WorkspaceName})
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
	commonTestHelper.ReplaceLabel(t, "客户端-管理禅道站点下工作目录")

	//update script title
	scriptPath := filepath.Join(workspacePath, "2_webpage_extract.php")
	content := fileUtils.ReadFile(scriptPath)
	newContent := strings.Replace(content, "title=extract content from webpage", "title=extract content from webpage-synctozentao", 1)
	fileUtils.WriteFile(scriptPath, newContent)
	defer fileUtils.WriteFile(scriptPath, content)

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	ztfTestHelper.SelectSite(webpage, "")
	ztfTestHelper.ExpandWorspace(webpage)

	locator := webpage.Locator(".tree-node-title", playwright.PageLocatorOptions{HasText: constTestHelper.WorkspaceName})
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
	commonTestHelper.ReplaceLabel(t, "客户端-管理禅道站点下工作目录")

	fileUtils.MkDirIfNeeded(syncDir)
	defer os.RemoveAll(syncDir)

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	ztfTestHelper.SelectSite(webpage, "")
	ztfTestHelper.ExpandProduct(webpage)

	scriptLocator := webpage.Locator(fmt.Sprintf(".tree-node:has-text('%s')>>.tree-node-title>>text=1_string_match.php", constTestHelper.WorkspaceName))
	scriptLocator.RightClick()

	webpage.Click(".tree-context-menu>>text=复制")
	productLocator := webpage.Locator(fmt.Sprintf(".tree-node:has-text('%s')>>.tree-node-item:has-text('product1')", constTestHelper.WorkspaceName))
	productLocator.RightClick()
	webpage.Click(".tree-context-menu>>text=粘贴")
	webpage.WaitForTimeout(1000)

	plwConf.DisableErr()
	defer plwConf.EnableErr()
	scriptLocator = webpage.Locator(fmt.Sprintf(".tree-node:has-text('%s')>>.tree-node-title>>text=1_string_match.php", constTestHelper.WorkspaceName))
	c := scriptLocator.Count()
	if c < 2 {
		t.Errorf("Find 1_string_match fail")
		t.FailNow()
		return
	}
}

func ClipAndCopyDir(t provider.T) {
	t.ID("5473,7593")
	commonTestHelper.ReplaceLabel(t, "客户端-管理禅道站点下工作目录")

	CopyDir(t)
	ClipDir(t)
}

func CopyDir(t provider.T) {
	fileUtils.MkDirIfNeeded(syncDir)
	defer os.RemoveAll(syncDir)
	if !fileUtils.FileExist(oldDirPath) {
		fileUtils.MkDirIfNeeded(oldDirPath)
	}
	if fileUtils.FileExist(newDirPath) {
		os.RemoveAll(newDirPath)
	}

	defer os.Remove(oldDirPath)
	defer os.Remove(newDirPath)

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	ztfTestHelper.SelectSite(webpage, "")
	ztfTestHelper.ExpandWorspace(webpage)
	ztfTestHelper.ExpandProduct(webpage)

	scriptLocator := webpage.Locator(".tree-node-title>>text=oldDir")
	scriptLocator.RightClick()

	webpage.Click(".tree-context-menu>>text=复制")
	productLocator := webpage.Locator(fmt.Sprintf(".tree-node:has-text('%s')>>.tree-node-item:has-text('product1')", constTestHelper.WorkspaceName))
	productLocator.RightClick()
	webpage.Click(".tree-context-menu>>text=粘贴")

	if webpage.ElementExist(fmt.Sprintf(".tree-node:has-text('%s')>>.tree-node-item:has-text('product1')>>.tree-node-title>>text=oldDir", constTestHelper.WorkspaceName)) {
		t.Errorf("Copy dir fail")
		t.FailNow()
		return
	}
}

func ClipDir(t provider.T) {
	fileUtils.MkDirIfNeeded(syncDir)
	defer os.RemoveAll(syncDir)
	if !fileUtils.FileExist(oldDirPath) {
		fileUtils.MkDirIfNeeded(oldDirPath)
	}
	if fileUtils.FileExist(newDirPath) {
		os.RemoveAll(newDirPath)
	}

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	ztfTestHelper.SelectSite(webpage, "")
	ztfTestHelper.ExpandWorspace(webpage)
	ztfTestHelper.ExpandProduct(webpage)

	scriptLocator := webpage.Locator(fmt.Sprintf(".tree-node:has-text('%s')>>.tree-node-title>>text=oldDir", constTestHelper.WorkspaceName))
	scriptLocator.RightClick()

	webpage.Click(".tree-context-menu>>text=剪切")
	productLocator := webpage.Locator(fmt.Sprintf(".tree-node:has-text('%s')>>.tree-node-item:has-text('product1')", constTestHelper.WorkspaceName))
	productLocator.RightClick()
	webpage.Click(".tree-context-menu>>text=粘贴")

	if webpage.ElementExist(fmt.Sprintf(".tree-node:has-text('%s')>>.tree-node-item:has-text('product1')>>.tree-node-title>>text=oldDir", constTestHelper.WorkspaceName)) {
		t.Errorf("Copy dir fail")
		t.FailNow()
		return
	}

	os.Remove(oldDirPath)
	os.Remove(newDirPath)
}

func CreateScript(t provider.T) {
	// t.ID("7596")
	// commonTestHelper.ReplaceLabel(t, "客户端-管理禅道站点下工作目录")
	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	ztfTestHelper.SelectSite(webpage, "")
	ztfTestHelper.ExpandProduct(webpage)

	if webpage.ElementExist(fmt.Sprintf(".tree-node:has-text('%s')>>.tree-node-title>>text=old.php", constTestHelper.WorkspaceName)) {
		return
	}

	workspaceLocator := webpage.Locator(fmt.Sprintf(".tree-node-title:has-text('%s')", constTestHelper.WorkspaceName))
	workspaceLocator.Hover()
	webpage.Click(`[title="新建脚本"]`)
	productLocator := webpage.Locator("#scriptFormModal>>input")
	productLocator.Fill("old.php")
	webpage.Click("#scriptFormModal>>.modal-action>>span:has-text(\"确定\")")

	webpage.Locator(fmt.Sprintf(".tree-node:has-text('%s')>>.tree-node-title>>text=old.php", constTestHelper.WorkspaceName))
}

func RenameScript(t provider.T) {
	t.ID("7596")
	commonTestHelper.ReplaceLabel(t, "客户端-管理禅道站点下工作目录")

	fileUtils.MkDirIfNeeded(syncDir)

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)

	defer func() {
		webpage.Close()
		os.RemoveAll(syncDir)
		os.Remove(workspacePath + "old.php")
		os.Remove(workspacePath + "new.php")
	}()

	ztfTestHelper.SelectSite(webpage, "")
	ztfTestHelper.ExpandWorspace(webpage)
	ztfTestHelper.ExpandProduct(webpage)
	CreateScript(t)

	scriptLocator := webpage.Locator(fmt.Sprintf(".tree-node:has-text('%s')>>.tree-node-title>>text=old.php", constTestHelper.WorkspaceName))
	scriptLocator.RightClick()

	webpage.Click(".tree-context-menu>>text=重命名")
	productLocator := webpage.Locator("#scriptFormModal>>input")
	productLocator.Fill("new.php")
	webpage.Click("#scriptFormModal>>.modal-action>>span:has-text(\"确定\")")

	scriptLocator = webpage.Locator(fmt.Sprintf(".tree-node:has-text('%s')>>.tree-node-title>>text=new.php", constTestHelper.WorkspaceName))
}

func CreateDir(t provider.T) {
	// t.ID("7596")
	// commonTestHelper.ReplaceLabel(t, "客户端-管理禅道站点下工作目录")
	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	ztfTestHelper.SelectSite(webpage, "")
	ztfTestHelper.ExpandProduct(webpage)

	if webpage.ElementExist(".tree-node-title>>text=oldDir") {
		return
	}

	workspaceLocator := webpage.Locator(fmt.Sprintf(".tree-node-title:has-text('%s')", constTestHelper.WorkspaceName))
	workspaceLocator.Hover()
	webpage.Click(fmt.Sprintf(".tree-node-item:has-text('%s')>>[title=\"新建工作目录\"]", constTestHelper.WorkspaceName))
	productLocator := webpage.Locator("#scriptFormModal>>input")
	productLocator.Fill("oldDir")
	webpage.Click("#scriptFormModal>>.modal-action>>span:has-text(\"确定\")")

	webpage.Locator(fmt.Sprintf(".tree-node:has-text('%s')>>.tree-node-title>>text=oldDir", constTestHelper.WorkspaceName))
}

func RenameDir(t provider.T) {
	t.ID("7595")
	commonTestHelper.ReplaceLabel(t, "客户端-管理禅道站点下工作目录")

	fileUtils.MkDirIfNeeded(syncDir)
	defer os.RemoveAll(syncDir)

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()
	defer os.RemoveAll(oldDirPath)
	defer os.RemoveAll(workspacePath + "newDir")

	ztfTestHelper.SelectSite(webpage, "")
	ztfTestHelper.ExpandWorspace(webpage)
	ztfTestHelper.ExpandProduct(webpage)
	CreateDir(t)

	dirLocator := webpage.Locator(fmt.Sprintf(".tree-node:has-text('%s')>>.tree-node-title>>text=oldDir", constTestHelper.WorkspaceName))
	dirLocator.RightClick()
	webpage.Click(".tree-context-menu>>text=重命名")
	productLocator := webpage.Locator("#scriptFormModal>>input")
	productLocator.Fill("newDir")
	webpage.Click("#scriptFormModal>>.modal-action>>span:has-text(\"确定\")")

	webpage.Locator(fmt.Sprintf(".tree-node:has-text('%s')>>.tree-node-title>>text=newDir", constTestHelper.WorkspaceName))
}

func DeleteScript(t provider.T) {
	t.ID("5478")
	commonTestHelper.ReplaceLabel(t, "客户端-管理禅道站点下工作目录")

	fileUtils.MkDirIfNeeded(syncDir)
	scriptPath := filepath.Join(workspacePath, "product1", "1.php")
	fileUtils.WriteFile(scriptPath, "")
	defer os.RemoveAll(syncDir)
	defer os.RemoveAll(scriptPath)

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	ztfTestHelper.SelectSite(webpage, "")
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
	commonTestHelper.ReplaceLabel(t, "客户端-管理禅道站点下工作目录")

	fileUtils.MkDirIfNeeded(syncDir)
	defer os.RemoveAll(syncDir)

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	ztfTestHelper.SelectSite(webpage, "")
	ztfTestHelper.ExpandWorspace(webpage)

	productLocator := webpage.Locator(fmt.Sprintf(".tree-node:has-text('%s')>>.tree-node-item:has-text('product1')", constTestHelper.WorkspaceName))
	productLocator.RightClick()
	webpage.Click(".tree-context-menu>>text=删除")
	webpage.Click(".modal-action>>span:has-text(\"确定\")")
	webpage.WaitForSelector(fmt.Sprintf(".tree-node:has-text('%s')>>.tree-node-item:has-text('product1')", constTestHelper.WorkspaceName), playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})

	plwConf.DisableErr()
	defer plwConf.EnableErr()
	scriptLocator := webpage.Locator(fmt.Sprintf(".tree-node:has-text('%s')>>.tree-node-item:has-text('product1')", constTestHelper.WorkspaceName))
	c := scriptLocator.Count()
	if c > 0 {
		t.Errorf("Delete workspace fail")
		t.FailNow()
		return
	}
}

func DeleteWorkspace(t provider.T, webpage plwHelper.Webpage) {
	locator := webpage.Locator(".tree-node-item", playwright.PageLocatorOptions{HasText: constTestHelper.WorkspaceName})
	locator.Hover()
	webpage.Click(`[title="删除"]`)
	webpage.Click(".modal-action>>span>>text=确定")

	webpage.WaitForTimeout(1000)
	plwConf.DisableErr()
	defer plwConf.EnableErr()

	scriptLocator := webpage.Locator(fmt.Sprintf(".tree-node-title>>text=%s", constTestHelper.WorkspaceName))
	c := scriptLocator.Count()
	if c > 0 {
		t.Errorf("Delete workspace fail")
		t.FailNow()
		return
	}
}
func Clip(t provider.T) {
	t.ID("5476")
	commonTestHelper.ReplaceLabel(t, "客户端-管理禅道站点下工作目录")

	fileUtils.MkDirIfNeeded(syncDir)
	scriptPath := filepath.Join(workspacePath, "product1", "1.php")
	scriptNewPath := filepath.Join(workspacePath, "1.php")
	fileUtils.WriteFile(scriptPath, "")

	defer func() {
		os.RemoveAll(syncDir)
		os.RemoveAll(scriptPath)
		os.RemoveAll(scriptNewPath)
	}()

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	ztfTestHelper.SelectSite(webpage, "")
	ztfTestHelper.ExpandWorspace(webpage)
	ztfTestHelper.ExpandProduct(webpage)

	scriptLocator := webpage.Locator(".tree-node-title>>text=1.php")
	scriptLocator.RightClick()
	webpage.Click(".tree-context-menu>>text=剪切")
	workspaceLocator := webpage.Locator(".tree-node-title", playwright.PageLocatorOptions{HasText: constTestHelper.WorkspaceName})
	workspaceLocator.RightClick()
	webpage.Click(".tree-context-menu>>text=粘贴")

	webpage.WaitForTimeout(1000)
	webpage.Locator(".tree-node-item>>div:has-text('1.php')")
}

func Collapse(t provider.T) {
	t.ID("5472")
	commonTestHelper.ReplaceLabel(t, "客户端-管理禅道站点下工作目录")

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()

	ztfTestHelper.SelectSite(webpage, "")
	ztfTestHelper.ExpandWorspace(webpage)

	webpage.WaitForSelectorTimeout(fmt.Sprintf(".tree-node:has-text(\"%s\")", constTestHelper.WorkspaceName), 5000)
	className := webpage.GetAttribute(fmt.Sprintf(".tree-node:has-text(\"%s\")", constTestHelper.WorkspaceName), "class")

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
	runner.Run(t, "客户端-从禅道同步", SyncFromZentao)
	runner.Run(t, "客户端-从禅道同步选中用例", SyncTwoCaseFromZentao)
	runner.Run(t, "客户端-复制粘贴树状脚本文件", Copy)
	runner.Run(t, "客户端-复制粘贴树状脚本文件夹", ClipAndCopyDir)
	runner.Run(t, "客户端-剪切粘贴树状脚本文件", Clip)
	runner.Run(t, "客户端-重命名脚本文件", RenameScript)
	runner.Run(t, "客户端-重命名目录", RenameDir)
	runner.Run(t, "客户端-删除树状脚本文件", DeleteScript)
	runner.Run(t, "客户端-删除树状脚本文件夹", DeleteDir)
	runner.Run(t, "客户端-显示展开折叠脚本树状结构", Collapse)
	runner.Run(t, "客户端-管理禅道工作目录", MangeWorkspace)
	runner.Run(t, "客户端-管理本地工作目录", MangeLocalWorkspace)
}
