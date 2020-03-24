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
)

func RetriveResult() []model.UnitTestSuite {
	sep := string(os.PathSeparator)

	resultDir := ""
	resultFiles := make([]string, 0)

	if vari.UnitTestType == "junit" && vari.UnitTestTool == "mvn" {
		resultDir = fmt.Sprintf("target%ssurefire-reports%s", sep, sep)
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

	suites := make([]model.UnitTestSuite, 0)
	for _, file := range resultFiles {
		content := fileUtils.ReadFile(file)

		testsuite := model.UnitTestSuite{}
		err = xml.Unmarshal([]byte(content), &testsuite)
		if err == nil {
			suites = append(suites, testsuite)
		}
	}

	return suites
}

func ParserUnitTestResult(testSuites []model.UnitTestSuite) ([]model.UnitTestCase, int) {
	cases := make([]model.UnitTestCase, 0)
	classNameMaxWidth := 0
	idx := 1
	for _, suite := range testSuites {
		for _, cs := range suite.Testcase {
			cs.Id = idx

			if cs.Failure != nil {
				cs.Status = "fail"

				cs.Failure.Desc = strings.Replace(cs.Failure.Desc, "<![CDATA[", "", -1)
				cs.Failure.Desc = strings.Replace(cs.Failure.Desc, "]]>", "", -1)
				logUtils.Screen(cs.Failure.Desc)
			} else {
				cs.Status = "pass"
			}

			lent2 := runewidth.StringWidth(cs.Classname)
			if lent2 > classNameMaxWidth {
				classNameMaxWidth = lent2
			}

			cases = append(cases, cs)
			idx++
		}
	}

	return cases, classNameMaxWidth
}
