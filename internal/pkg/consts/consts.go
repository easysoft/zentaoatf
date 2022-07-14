package commConsts

import "os"

const (
	App        = "ztf"
	AppServer  = "server"
	AppAgent   = "agent"
	AppCommand = "cmd"

	ConfigVersion      = "3.0"
	ConfigDir          = "conf"
	ConfigFile         = "ztf.conf"
	LogDirName         = "log"
	LogBakDirName      = "log-bak"
	ExtNameSuite       = "cs"
	LogText            = "log.txt"
	ResultText         = "result.txt"
	ResultJson         = "result.json"
	ResultZip          = "result.zip"
	ExecZip            = "exec.zip"
	ExecZipPath        = "uploadTmp"
	DownloadServerPath = "serverTmp"
	DownloadPath       = "downloadTmp"
	ExecProxyPath      = "execTmp"

	ExpectResultPass = "pass"

	PathInfo = "PATH_INFO"
	Get      = "GET"
	PthSep   = string(os.PathSeparator)
)

var (
	UnitBuildToolMap = map[string]BuildTool{
		"mvn": Maven,
	}

	AutoTestTypeSelenium = "selenium"
	AutoTestTypeAppium   = "appium"
	AutoTestTypes        = []string{AutoTestTypeSelenium, AutoTestTypeAppium}

	UnitTestTypeJunit      = "junit"
	UnitTestTypeTestNG     = "testng"
	UnitTestTypeRobot      = "robot"
	UnitTestTypeCypress    = "cypress"
	UnitTestTypePlaywright = "playwright"
	UnitTestTypePuppeteer  = "puppeteer"
	UnitTestPhpUnit        = "phpunit"
	UnitTestTypePyTest     = "pytest"
	UnitTestTypeJest       = "jest"
	UnitTestTypeCppUnit    = "cppunit"
	UnitTestTypeGTest      = "gtest"
	UnitTestTypeQTest      = "qtest"
	UnitTestTypes          = []string{
		UnitTestTypeJunit, UnitTestTypeTestNG,
		UnitTestTypeRobot, UnitTestTypeCypress, UnitTestTypePlaywright, UnitTestTypePuppeteer,
		UnitTestPhpUnit, UnitTestTypePyTest, UnitTestTypeJest, UnitTestTypeCppUnit, UnitTestTypeGTest, UnitTestTypeQTest,
	}

	UnitTestToolMvn   = "mvn"
	UnitTestToolMocha = "mocha"
	UnitTestToolRobot = "robot"
)
