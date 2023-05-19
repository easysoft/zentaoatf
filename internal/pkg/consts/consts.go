package commConsts

import "os"

const (
	App        = "ztf"
	AppServer  = "server"
	AppAgent   = "agent"
	AppCommand = "cmd"

	Ip   = "127.0.0.1"
	Port = 8085

	JobTimeoutTime = 60 * 30
	JobRetryTime   = 3

	ConfigVersion      = "3.0"
	ConfigDir          = "conf"
	ConfigFile         = "ztf.conf"
	LogDirName         = "log"
	ExtNameSuite       = "cs"
	LogText            = "log.txt"
	ResultText         = "result.txt"
	ResultJson         = "result.json"
	ResultZip          = "result.zip"
	ExecZip            = "exec.zip"
	ExecZipPath        = "uploadTmp"
	DownloadServerPath = "serverTmp"
	DownloadPath       = "downloadTmp"
	ExecProxyPath      = "proxyExecDir"

	ExpectResultPass = "pass"

	PathInfo = "PATH_INFO"
	Get      = "GET"
	PthSep   = string(os.PathSeparator)
)

var (
	UnitBuildToolMap = map[string]BuildTool{
		"mvn": Maven,
	}

	SpaceQuote = " "

	AutoTestTypes = []string{Selenium.String(), Appium.String(), AutoIt.String()}
	UnitTestTypes = []string{
		Allure.String(),
		JUnit.String(), TestNG.String(), PHPUnit.String(), PyTest.String(), Jest.String(), CppUnit.String(), GTest.String(), QTest.String(),
		RobotFramework.String(), Cypress.String(), Playwright.String(), Puppeteer.String(), K6.String(), Zap.String(),
	}
)
