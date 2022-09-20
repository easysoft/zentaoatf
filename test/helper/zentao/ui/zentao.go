package uiTest

import (
	"strings"

	playwright "github.com/playwright-community/playwright-go"
)

var page playwright.Page

func Login() (err error) {
	if _, err = page.Goto("http://127.0.0.1:8081/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		return
	}
	err = page.Fill(`input[name="account"]`, "admin")
	if err != nil {
		return
	}
	err = page.Fill(`input[name="password"]`, "123456.")
	if err != nil {
		return
	}
	err = page.Click(`button:has-text("登录")`)
	if err != nil {
		return
	}
	_, err = page.WaitForSelector("#login", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})
	if err != nil {
		return
	}
	title, err := page.Title()
	if err != nil {
		return
	}
	if title == "流程 - 禅道" {
		err = page.Click(`button:has-text("保存")`)
		if err != nil {
			return
		}
	}
	return
}

func CreateModule() (err error) {
	if _, err = page.Goto("http://127.0.0.1:8081/index.php?m=tree&f=browse&productID=1&view=case", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		return
	}
	err = page.Fill(`input[name="modules\[\]"]>>nth=0`, "module1")
	if err != nil {
		return
	}
	err = page.Fill(`input[name="modules\[\]"]>>nth=1`, "module2")
	if err != nil {
		return
	}
	err = page.Click(`#submit`)
	if err != nil {
		return
	}
	_, err = page.WaitForSelector("#modulesTree>>:has-text('module')")
	if err != nil {
		return
	}
	return
}

func CreateSuite() (err error) {
	if _, err = page.Goto("http://127.0.0.1:8081/index.php?m=testsuite&f=create&product=1", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		return
	}
	err = page.Fill(`#name`, "test_suite")
	if err != nil {
		return
	}
	err = page.Click(`#submit`)
	if err != nil {
		return
	}
	_, err = page.WaitForSelector("#submit", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})
	if err != nil {
		return
	}
	return
}

func goToLastUnitTestResult() {

}

func checkUnitTestResult() {

}

func InitZentaoData() (err error) {
	if _, err = page.Goto("http://127.0.0.1:8081/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		return
	}
	title, err := page.Title()
	if err != nil {
		return
	}
	if strings.Contains(title, "欢迎使用禅道") {
		err = page.Click("text=开始安装")
		if err != nil {
			return
		}
		err = page.Click("text=下一步")
		if err != nil {
			return
		}
		err = page.Click("text=下一步")
		if err != nil {
			return
		}
		err = page.Fill(`input[name="dbPassword"]`, "123456")
		if err != nil {
			return
		}
		err = page.Click(`input[name="clearDB\[\]"]`)
		if err != nil {
			return
		}
		err = page.Click("text=保存")
		if err != nil {
			return
		}
		err = page.Click("text=下一步")
		if err != nil {
			return
		}
		err = page.Fill(`input[name="company"]`, "test")
		if err != nil {
			return
		}
		err = page.Fill(`input[name="account"]`, "admin")
		if err != nil {
			return
		}
		err = page.Fill(`input[name="password"]`, "123456.")
		if err != nil {
			return
		}
		err = page.Click(`input[name="importDemoData\[\]"]`)
		if err != nil {
			return
		}
		err = page.Click("text=保存")
		if err != nil {
			return
		}
		_, err = page.WaitForSelector("text=保存", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})
		if err != nil {
			return
		}
		err = Login()
		if err != nil {
			return
		}
		err = CreateModule()
		if err != nil {
			return
		}
		err = CreateSuite()
		if err != nil {
			return
		}
	}

	return
}

func init() {
	pw, err := playwright.Run()
	if err != nil {
		return
	}
	runBrowser, err := pw.Chromium.Launch()
	if err != nil {
		return
	}
	defer func() {
		// if err = runBrowser.Close(); err != nil {
		// 	return
		// }
		// if err = pw.Stop(); err != nil {
		// 	return
		// }
	}()
	page, err = runBrowser.NewPage()
	if err != nil {
		return
	}
}
