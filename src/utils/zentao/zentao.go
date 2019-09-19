package zentaoUtils

import (
	"fmt"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	dateUtils "github.com/easysoft/zentaoatf/src/utils/date"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GenSuperApiUri(model string, methd string, params [][]string) string {
	var sep string
	if vari.RequestType == constant.RequestTypePathInfo {
		sep = ","
	} else {
		sep = "&"
	}

	paramStr := ""
	i := 0
	for _, p := range params {
		if i > 0 {
			paramStr += sep
		}
		paramStr += p[0] + "=" + p[1]
		i++
	}

	var uri string
	if vari.RequestType == constant.RequestTypePathInfo {
		uri = fmt.Sprintf("api-getmodel-%s-%s-%s.json", model, methd, paramStr)
	} else {
		uri = fmt.Sprintf("?m=api&f=getmodel&model=%s&methodName=%s&params=%s", model, methd, paramStr)
	}
	return uri
}

func GenApiUri(module string, methd string, param string) string {
	if vari.RequestType == constant.RequestTypePathInfo {
		return fmt.Sprintf("%s-%s-%s.json", module, methd, param)
	}

	return ""
}

func ScriptToExpectName(file string) string {
	fileSuffix := path.Ext(file)
	expectName := strings.TrimSuffix(file, fileSuffix) + ".exp"

	return expectName
}

func RunDateFolder() string {
	runName := dateUtils.DateTimeStrFmt(time.Now(), "2006-01-02T150405") + string(os.PathSeparator)

	return runName
}

func GetCaseInfo(file string) (bool, int, int, string) {
	var caseId int
	var productId int
	var title string

	content := fileUtils.ReadFile(file)

	pass := CheckFileContentIsScript(content)
	if !pass {
		return false, caseId, productId, title
	}

	caseInfo := ""
	myExp := regexp.MustCompile(`(?s)\[case\](.*)\[esac\]`)
	arr := myExp.FindStringSubmatch(content)
	if len(arr) > 1 {
		caseInfo = arr[1]
	}

	myExp = regexp.MustCompile(`[\S\s]*cid=\s*([^\n]*?)\s*\n`)
	arr = myExp.FindStringSubmatch(caseInfo)
	if len(arr) > 1 {
		caseId, _ = strconv.Atoi(arr[1])
	}

	myExp = regexp.MustCompile(`[\S\s]*pid=\s*([^\n]*?)\s*\n`)
	arr = myExp.FindStringSubmatch(caseInfo)
	if len(arr) > 1 {
		productId, _ = strconv.Atoi(arr[1])
	}

	myExp = regexp.MustCompile(`[\S\s]*title=\s*([^\n]*?)\s*\n`)
	arr = myExp.FindStringSubmatch(caseInfo)
	if len(arr) > 1 {
		title = arr[1]
	}

	return pass, caseId, productId, title
}

func ReadCheckpoints(file string) ([]string, [][]string) {
	expectIndependentFile := strings.Replace(file, path.Ext(file), ".exp", -1)
	expectIndependentContent := ""

	if fileUtils.FileExist(expectIndependentFile) {
		expectIndependentContent = fileUtils.ReadFile(expectIndependentFile)
	}

	content := fileUtils.ReadFile(file)

	myExp := regexp.MustCompile(`(?s)\[case\](?U:.*)(\[.*)\[esac\]`)
	arr := myExp.FindStringSubmatch(content)

	if len(arr) > 1 {
		checkpoints := arr[1]
		content = commonUtils.RemoveBlankLine(checkpoints)
	}

	cpStepArr, expectArr := genCheckpointStepArr(content, expectIndependentContent)

	return cpStepArr, expectArr
}
func genCheckpointStepArr(content string, expectIndependentContent string) ([]string, [][]string) {
	cpStepArr := make([]string, 0)
	expectArr := make([][]string, 0)

	independentExpect := expectIndependentContent != ""

	lines := strings.Split(content, "\n")
	i := 0
	for i < len(lines) {
		step := ""
		expects := make([]string, 0)

		line := strings.TrimSpace(lines[i])

		regx := regexp.MustCompile(`([\d\.]+).*>>(.*)`)
		arr := regx.FindStringSubmatch(line)
		if len(arr) > 1 {
			step = arr[1]
			if !independentExpect {
				expects = append(expects, arr[2])
			}
		} else {
			regx = regexp.MustCompile(`\[([\d\.]*).*expects\]`)
			arr = regx.FindStringSubmatch(line)
			if len(arr) > 1 {
				step = arr[1]

				if !independentExpect {
					for i+1 < len(lines) {
						ln := strings.TrimSpace(lines[i+1])

						if strings.Index(ln, "[") == 0 || strings.Index(ln, ">>") > 0 || ln == "" {
							break
						} else {
							expects = append(expects, ln)
							i++
						}
					}
				}
			}
		}

		if step != "" {
			cpStepArr = append(cpStepArr, step)
			if !independentExpect {
				expectArr = append(expectArr, expects)
			}
		}
		i++
	}

	if independentExpect {
		expectIndependentArr := strings.Split(expectIndependentContent, "\n")
		expectArr = expectIndependentArr
	}

	return cpStepArr, expectArr
}

func GenLogArr(str string) (bool, [][]string) {
	skip, arr := GenArr(str, true)
	return skip, arr
}
func GenArr(str string, checkSkip bool) (bool, [][]string) {
	ret := make([][]string, 0)
	indx := -1
	for _, line := range strings.Split(str, "\n") {
		line := strings.TrimSpace(line)

		if checkSkip && strings.ToLower(line) == "skip" {
			return true, nil
		}

		if strings.Index(line, "#") == 0 {
			ret = append(ret, make([]string, 0))
			indx++
		} else if indx > -1 {
			if len(line) > 0 && indx < len(ret) {
				ret[indx] = append(ret[indx], line)
			}
		}
	}

	return false, ret
}

func ReadLog(logFile string) (bool, [][]string) {
	str := fileUtils.ReadFile(logFile)

	skip, ret := GenLogArr(str)
	return skip, ret
}

func CheckFileIsScript(path string) bool {
	content := fileUtils.ReadFile(path)

	pass := CheckFileContentIsScript(content)

	return pass
}

func CheckFileContentIsScript(content string) bool {
	pass, _ := regexp.MatchString(`\[case\]`, content)

	return pass
}
