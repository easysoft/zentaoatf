package testingService

import (
	"encoding/xml"
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/mattn/go-runewidth"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

func RetrieveUnitResult() []model.UnitTestSuite {
	sep := string(os.PathSeparator)

	resultDir := ""
	resultFiles := make([]string, 0)

	if vari.UnitTestType == "junit" && vari.UnitTestTool == "mvn" {
		resultDir = fmt.Sprintf("target%ssurefire-reports%s", sep, sep)
	} else if vari.UnitTestType == "testng" && vari.UnitTestTool == "mvn" {
		resultDir = fmt.Sprintf("target%ssurefire-reports%sjunitreports", sep, sep)
	} else if vari.UnitTestType == "robot" {
		resultDir = vari.UnitTestResults
	} else {
		resultDir = vari.UnitTestResult
	}

	if fileUtils.IsDir(resultDir) {
		if resultDir[len(resultDir)-1:] != sep {
			resultDir = resultDir + sep
		}

		dir, err := ioutil.ReadDir(resultDir)
		if err == nil {
			for _, fi := range dir {
				name := fi.Name()
				ext := path.Ext(name)
				if ext == ".xml" {
					resultFiles = append(resultFiles, resultDir+name)
				}
			}
		}
	} else {
		resultFiles = append(resultFiles, resultDir)
	}

	suites := make([]model.UnitTestSuite, 0)
	for _, file := range resultFiles {
		content := fileUtils.ReadFile(file)

		var err error
		var testSuite model.UnitTestSuite

		if vari.UnitTestType == "gtest" {
			gTestSuite := model.GTestSuites{}
			err = xml.Unmarshal([]byte(content), &gTestSuite)
			if err == nil {
				testSuite = ConvertGTestResult(gTestSuite)
			}
		} else if vari.UnitTestType == "jest" {
			jestSuite := model.JestSuites{}
			err = xml.Unmarshal([]byte(content), &jestSuite)
			if err == nil {
				testSuite = ConvertJestResult(jestSuite)
			}
		} else if vari.UnitTestType == "phpunit" {
			phpTestSuite := model.PhpUnitSuites{}
			err = xml.Unmarshal([]byte(content), &phpTestSuite)
			if err == nil {
				testSuite = ConvertPhpUnitResult(phpTestSuite)
			}
		} else if vari.UnitTestType == "pytest" {
			pyTestSuite := model.PyTestSuites{}
			err = xml.Unmarshal([]byte(content), &pyTestSuite)
			if err == nil {
				testSuite = ConvertPyTestResult(pyTestSuite)
			}
		} else if vari.UnitTestType == "cppunit" {
			content = strings.Replace(content, "ISO-8859-1", "UTF-8", -1)

			cppUnitSuites := model.CppUnitSuites{}
			err = xml.Unmarshal([]byte(content), &cppUnitSuites)
			if err == nil {
				testSuite = ConvertCppUnitResult(cppUnitSuites)
			}
		} else if vari.UnitTestType == "qtest" {
			qTestSuite := model.QTestSuites{}
			err = xml.Unmarshal([]byte(content), &qTestSuite)
			if err == nil {
				testSuite = ConvertQTestResult(qTestSuite)
			}
		} else if vari.UnitTestType == "robot" {
			robotResult := model.RobotResult{}
			err = xml.Unmarshal([]byte(content), &robotResult)
			if err == nil {
				testSuite = ConvertRobotResult(robotResult)
			}
		} else { // junit, testng
			testSuite = model.UnitTestSuite{}
			err = xml.Unmarshal([]byte(content), &testSuite)
		}

		if err == nil {
			suites = append(suites, testSuite)
		}
	}

	//suites[0].Time = suites[0].Time + 10000 // TODO: delay for testing

	return suites
}

func ParserUnitTestResult(testSuites []model.UnitTestSuite) (cases []model.UnitResult, classNameMaxWidth int, dur float32) {
	idx := 1
	for _, suite := range testSuites {
		if suite.Time != 0 { // for junit, there is a time on suite level
			dur += suite.Time
		}

		for _, cs := range suite.TestCases {
			cs.Id = idx

			if cs.Failure != nil {
				cs.Status = "fail"

				cs.Failure.Desc = strings.Replace(cs.Failure.Desc, "<![CDATA[", "", -1)
				cs.Failure.Desc = strings.Replace(cs.Failure.Desc, "]]>", "", -1)
				logUtils.Screen(cs.Failure.Desc)
			} else {
				cs.Status = "pass"
			}

			lent2 := runewidth.StringWidth(cs.TestSuite)
			if lent2 > classNameMaxWidth {
				classNameMaxWidth = lent2
			}

			cases = append(cases, cs)
			idx++
		}
	}

	return
}

func ConvertJestResult(jestSuite model.JestSuites) model.UnitTestSuite {
	testSuite := model.UnitTestSuite{}
	testSuite.Time = jestSuite.Time

	for _, suite := range jestSuite.TestSuites {
		for _, cs := range suite.TestCases {
			caseResult := model.UnitResult{}
			caseResult.Title = cs.Title
			caseResult.Duration = cs.Duration

			if suite.Title != "" && suite.Title != "undefined" {
				caseResult.TestSuite = suite.Title
			} else {
				caseResult.TestSuite = jestSuite.Title
			}

			caseResult.Failure = cs.Failure

			testSuite.TestCases = append(testSuite.TestCases, caseResult)
		}
	}

	return testSuite
}

func ConvertPhpUnitResult(phpUnitSuite model.PhpUnitSuites) model.UnitTestSuite {
	testSuite := model.UnitTestSuite{}

	var total float32 = 0
	for _, cs := range phpUnitSuite.TestCases {
		caseResult := model.UnitResult{}
		caseResult.Title = cs.Title
		caseResult.Duration = cs.Time

		total += cs.Time

		if cs.Groups != "" && cs.Groups != "default" {
			caseResult.TestSuite = cs.Groups
		} else {
			caseResult.TestSuite = cs.TestSuite
		}

		if cs.Status != 0 {
			fail := model.Failure{}
			fail.Desc = cs.Fail
			caseResult.Failure = &fail
		}

		testSuite.TestCases = append(testSuite.TestCases, caseResult)
	}
	testSuite.Duration = int64(total)
	testSuite.Time = total

	return testSuite
}

func ConvertPyTestResult(pytestSuites model.PyTestSuites) model.UnitTestSuite {
	testSuite := model.UnitTestSuite{}

	var total float32 = 0
	for _, suite := range pytestSuites.TestSuites {
		total += suite.Time

		for _, cs := range suite.TestCases {
			caseResult := model.UnitResult{}
			caseResult.Title = cs.Title
			caseResult.Duration = cs.Duration

			if suite.Title != "" && suite.Title != "pytest" {
				caseResult.TestSuite = suite.Title
			} else {
				caseResult.TestSuite = cs.TestSuite
			}

			if cs.Failure != nil {
				fail := model.Failure{}
				fail.Type = cs.Failure.Type
				fail.Desc = cs.Failure.Desc
				caseResult.Failure = &fail
			} else if cs.Error != nil {
				fail := model.Failure{}
				fail.Type = cs.Error.Message
				fail.Desc = cs.Error.Text
				caseResult.Failure = &fail
			}

			testSuite.TestCases = append(testSuite.TestCases, caseResult)

		}
	}

	testSuite.Duration = int64(total)
	testSuite.Time = total

	return testSuite
}

func ConvertGTestResult(gTestSuite model.GTestSuites) model.UnitTestSuite {
	testSuite := model.UnitTestSuite{}
	testSuite.Time = gTestSuite.Time

	for _, suite := range gTestSuite.TestSuites {
		for _, cs := range suite.TestCases {
			caseResult := model.UnitResult{}
			caseResult.Title = cs.Title
			caseResult.Duration = cs.Duration
			caseResult.Status = cs.Status

			if suite.Title != "" && suite.Title != "pytest" {
				caseResult.TestSuite = suite.Title
			}

			if cs.Failure != nil {
				fail := model.Failure{}
				fail.Type = cs.Failure.Type
				fail.Desc = cs.Failure.Desc
				caseResult.Failure = &fail
			}

			testSuite.TestCases = append(testSuite.TestCases, caseResult)

		}
	}

	return testSuite
}

func ConvertCppUnitResult(cppunitSuite model.CppUnitSuites) model.UnitTestSuite {
	testSuite := model.UnitTestSuite{}

	for _, cs := range cppunitSuite.FailedTests.TestCases {
		caseResult := model.UnitResult{}
		caseResult.Id = cs.Id
		caseResult.Title = cs.Title

		fail := model.Failure{}
		fail.Type = cs.FailureType
		fail.Desc = cs.Message
		caseResult.Failure = &fail

		testSuite.TestCases = append(testSuite.TestCases, caseResult)
	}

	for _, cs := range cppunitSuite.SuccessfulTests.TestCases {
		caseResult := model.UnitResult{}
		caseResult.Id = cs.Id
		caseResult.Title = cs.Title

		testSuite.TestCases = append(testSuite.TestCases, caseResult)
	}

	return testSuite
}

func ConvertQTestResult(qTestSuite model.QTestSuites) model.UnitTestSuite {
	testSuite := model.UnitTestSuite{}

	for _, cs := range qTestSuite.TestCases {
		caseResult := model.UnitResult{}
		caseResult.TestSuite = qTestSuite.Name
		caseResult.Title = cs.Title
		caseResult.Status = cs.Result

		if cs.Failure != nil {
			fail := model.Failure{}
			fail.Type = cs.Failure.Type
			fail.Desc = cs.Failure.Desc
			caseResult.Failure = &fail
		}

		testSuite.TestCases = append(testSuite.TestCases, caseResult)
	}

	return testSuite
}

func ConvertRobotResult(result model.RobotResult) model.UnitTestSuite {
	testSuite := model.UnitTestSuite{}

	suiteMap := map[string]string{}
	for _, state := range result.Statistics.Suite.States {
		suiteMap[state.ID] = state.Text
	}

	tests := make([]model.RobotTest, 0)
	for _, suite := range result.Suites {
		RetrieveRobotTests(suite, &tests)
	}

	for _, cs := range tests {
		caseResult := model.UnitResult{}
		caseResult.Title = cs.Name
		caseResult.Status = strings.ToLower(cs.Status.Status)

		suiteId := cs.ID[0:strings.LastIndex(cs.ID, "-")]
		caseResult.TestSuite = suiteMap[suiteId]

		templ := "20060102 15:04:05.000"
		startTime, _ := time.ParseInLocation(templ, cs.Status.StartTime, time.Local)
		endTime, _ := time.ParseInLocation(templ, cs.Status.EndTime, time.Local)

		caseResult.StartTime = startTime.Unix()
		caseResult.EndTime = endTime.Unix()
		caseResult.Duration = float32(caseResult.EndTime - caseResult.StartTime)

		if caseResult.Status != "pass" {
			fail := model.Failure{}
			fail.Type = ""
			fail.Desc = cs.Status.Text
			caseResult.Failure = &fail
		}

		testSuite.TestCases = append(testSuite.TestCases, caseResult)
	}

	return testSuite
}

func RetrieveRobotTests(suite model.RobotSuite, tests *[]model.RobotTest) {
	for _, suite := range suite.Suites {
		RetrieveRobotTests(suite, tests)
	}

	for _, test := range suite.Tests {
		*tests = append(*tests, test)
	}
}
