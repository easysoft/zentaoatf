package testingService

import (
	"encoding/xml"
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"io/ioutil"
	"os"
	"path"
)

func RetriveResult() model.UnitTestSuite {
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
				break
			}
		}
	}

	testsuite := model.UnitTestSuite{}
	if len(resultFile) == 0 {
		return testsuite
	}

	content := fileUtils.ReadFile(resultFile)
	err = xml.Unmarshal([]byte(content), &testsuite)
	if err != nil {
		return testsuite
	}

	return testsuite
}
