package zentaoUtils

import (
	"fmt"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	langUtils "github.com/easysoft/zentaoatf/src/utils/lang"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

//func GenSuperApiUri(model string, methd string, params [][]string) string {
//	var sep string
//	if vari.RequestType == constant.RequestTypePathInfo {
//		sep = ","
//	} else {
//		sep = "&"
//	}
//
//	paramStr := "" // format is key1=val1&key2=val2
//	i := 0
//	for _, p := range params {
//		if i > 0 {
//			paramStr += sep
//		}
//		paramStr += p[0] + "=" + p[1]
//		i++
//	}
//
//	var uri string
//	if vari.RequestType == constant.RequestTypePathInfo {
//		uri = fmt.Sprintf("api-getmodel-%s-%s-%s.json", model, methd, paramStr)
//	} else {
//		uri = fmt.Sprintf("index.php?m=api&f=getmodel&model=%s&methodName=%s&params=%s&t=json", model, methd, paramStr)
//	}
//	return uri
//}

func GenApiUri(module string, methd string, param string) string {
	var uri string

	if vari.RequestType == constant.RequestTypePathInfo {
		uri = fmt.Sprintf("%s-%s-%s.json", module, methd, param)
	} else {
		uri = fmt.Sprintf("index.php?m=%s&f=%s&%s&t=json", module, methd, param)
	}

	return uri
}

func ScriptToExpectName(file string) string {
	fileSuffix := path.Ext(file)
	expectName := strings.TrimSuffix(file, fileSuffix) + ".exp"

	return expectName
}

//func RunDateFolder() string {
//	runName := dateUtils.DateTimeStrFmt(time.Now(), "2006-01-02T150405") + string(os.PathSeparator)
//
//	return runName
//}

func GetCaseInfo(file string) (bool, int, int, string) {
	var caseId int
	var productId int
	var title string

	content := fileUtils.ReadFile(file)
	isOldFormat := strings.Index(content, "[esac]") > -1
	pass := CheckFileContentIsScript(content)
	if !pass {
		return false, caseId, productId, title
	}

	caseInfo := ""
	lang := langUtils.GetLangByFile(file)
	regStr := ""
	if isOldFormat {
		regStr = `(?s)\[case\](.*)\[esac\]`
	} else {
		regStr = fmt.Sprintf(`(?Us)%s(.*)%s`,
			constant.LangCommentsMap[lang][0], constant.LangCommentsMap[lang][1])
	}
	myExp := regexp.MustCompile(regStr)
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

//func ReadScriptCheckpoints(file string) ([]string, [][]string) {
//	_, expectIndependentContent := GetDependentExpect(file)
//
//	content := fileUtils.ReadFile(file)
//	_, checkpoints := ReadCaseInfo(content)
//
//	cpStepArr, expectArr := getCheckpointStepArr(checkpoints, expectIndependentContent)
//
//	return cpStepArr, expectArr
//}
//func getCheckpointStepArr(content string, expectIndependentContent string) ([]string, [][]string) {
//	cpStepArr := make([]string, 0)
//	expectArr := make([][]string, 0)
//
//	independentExpect := expectIndependentContent != ""
//
//	lines := strings.Split(content, "\n")
//	i := 0
//	for i < len(lines) {
//		step := ""
//		expects := make([]string, 0)
//
//		line := strings.TrimSpace(lines[i])
//
//		regx := regexp.MustCompile(`(?U:[\d\.]*)(.+)>>(.*)`)
//		arr := regx.FindStringSubmatch(line)
//		if len(arr) > 2 {
//			step = arr[1]
//			if !independentExpect {
//				expects = append(expects, strings.TrimSpace(arr[2]))
//			}
//		} else {
//			regx = regexp.MustCompile(`\[([\d\.]*).*expects\]`)
//			arr = regx.FindStringSubmatch(line)
//			if len(arr) > 1 {
//				step = arr[1]
//
//				if !independentExpect {
//					for i+1 < len(lines) {
//						ln := strings.TrimSpace(lines[i+1])
//
//						if strings.Index(ln, "[") == 0 || strings.Index(ln, ">>") > 0 || ln == "" {
//							break
//						} else {
//							expects = append(expects, ln)
//							i++
//						}
//					}
//				}
//			}
//		}
//
//		if step != "" && len(expects) > 0 {
//			cpStepArr = append(cpStepArr, step)
//			if !independentExpect {
//				expectArr = append(expectArr, expects)
//			}
//		}
//		i++
//	}
//
//	if independentExpect {
//		expectArr = ReadExpectIndependentArr(expectIndependentContent)
//	}
//
//	return cpStepArr, expectArr
//}

func ReadExpectIndependentArr(content string) [][]string {
	lines := strings.Split(content, "\n")

	ret := make([][]string, 0)
	var cpArr []string

	for idx, line := range lines {
		line = strings.TrimSpace(line)

		if line == ">>" { // more than one line
			cpArr = make([]string, 0)
		} else if strings.Index(line, ">>") == 0 { // single line
			line = strings.Replace(line, ">>", "", -1)
			line = strings.TrimSpace(line)

			cpArr = append(cpArr, line)
			ret = append(ret, cpArr)
			cpArr = make([]string, 0)
		} else { // under >>
			cpArr = append(cpArr, line)

			if idx == len(lines)-1 || strings.Index(lines[idx+1], ">>") > -1 {
				ret = append(ret, cpArr)
				cpArr = make([]string, 0)
			}
		}
	}

	return ret
}

func ReadExpectIndependentArrNew(content string) [][]string {
	lines := strings.Split(content, "\n")

	ret := make([][]string, 0)
	var cpArr []string

	model := ""
	for idx, line := range lines {
		line = strings.TrimSpace(line)

		if line == ">>" { // more than one line
			model = "multi"
			cpArr = make([]string, 0)
		} else if model == "multi" { // in >> and << in multi line mode
			cpArr = append(cpArr, line)

			if idx == len(lines)-1 || strings.Index(lines[idx+1], "<<") > -1 {
				temp := make([]string, 0)
				temp = append(temp, strings.Join(cpArr, " | "))

				ret = append(ret, temp)
				cpArr = make([]string, 0)
				model = ""
			}
		} else if line == "<<" {
			continue
		} else {
			model = "single"

			line = strings.TrimSpace(line)

			cpArr = append(cpArr, line)
			ret = append(ret, cpArr)
			cpArr = make([]string, 0)
		}
	}

	return ret
}

func ReadLogArr(content string) (isSkip bool, ret [][]string) {
	lines := strings.Split(content, "\n")

	ret = make([][]string, 0)
	var cpArr []string

	model := ""
	for idx, line := range lines {
		line = strings.TrimSpace(line)

		if line == "skip" {
			isSkip = true
			return
		}

		if line == ">>" { // more than one line
			model = "multi"
			cpArr = make([]string, 0)
		} else if strings.Index(line, ">>") == 0 { // single line
			model = "single"

			line = strings.Replace(line, ">>", "", -1)
			line = strings.TrimSpace(line)

			cpArr = append(cpArr, line)
			ret = append(ret, cpArr)
			cpArr = make([]string, 0)
		} else {
			if model == "" || model == "single" {
				continue
			}

			// under >>
			cpArr = append(cpArr, line)

			if idx == len(lines)-1 || strings.Index(lines[idx+1], ">>") > -1 {
				ret = append(ret, cpArr)
				cpArr = make([]string, 0)
			}
		}
	}

	return
}

func ReadLogArrNew(content string) (isSkip bool, ret [][]string) {
	lines := strings.Split(content, "\n")

	ret = make([][]string, 0)
	var cpArr []string

	model := ""
	for idx, line := range lines {
		line = strings.TrimSpace(line)

		if line == "skip" {
			isSkip = true
			return
		}

		if line == ">>" { // more than one line
			model = "multi"
			cpArr = make([]string, 0)
		} else if model == "multi" { // in >> and << in multi line mode
			cpArr = append(cpArr, line)

			if idx == len(lines)-1 || strings.Index(lines[idx+1], "<<") > -1 {
				temp := make([]string, 0)
				temp = append(temp, strings.Join(cpArr, " | "))

				ret = append(ret, temp)
				cpArr = make([]string, 0)
				model = ""
			}
		} else if line == "<<" {
			continue
		} else {
			model = "single"

			line = strings.TrimSpace(line)

			cpArr = append(cpArr, line)
			ret = append(ret, cpArr)
			cpArr = make([]string, 0)
		}
	}

	return
}

func CheckFileIsScript(path string) bool {
	content := fileUtils.ReadFile(path)

	pass := CheckFileContentIsScript(content)

	return pass
}

func CheckFileContentIsScript(content string) bool {
	pass, _ := regexp.MatchString(`cid\b*=`, content)

	return pass
}

func ReadCaseInfo(content, lang string, isOldFormat bool) (info, checkpoints string) {
	regStr := ""
	if isOldFormat {
		regStr = `(?s)\[case\]((?U:.*pid.*))\n(.*)\[esac\]`
	} else {
		regStr = fmt.Sprintf(`(?s)%s((?U:.*pid.*))\n(.*)%s`,
			constant.LangCommentsMap[lang][0], constant.LangCommentsMap[lang][1])
	}
	myExp := regexp.MustCompile(regStr)
	arr := myExp.FindStringSubmatch(content)

	if len(arr) > 2 {
		info = strings.TrimSpace(arr[1])
		checkpoints = strings.TrimSpace(arr[2])

		return
	}

	return
}
func ReadCaseId(content string) string {
	myExp := regexp.MustCompile(`(?s).*\ncid=((?U:.*))\n.*`)
	arr := myExp.FindStringSubmatch(content)

	if len(arr) > 1 {
		id := strings.TrimSpace(arr[1])
		return id
	}

	return ""
}

func GetDependentExpect(file string) (bool, string) {
	dir := fileUtils.AddPathSepIfNeeded(filepath.Dir(file))
	name := strings.Replace(filepath.Base(file), path.Ext(file), ".exp", -1)
	expectIndependentFile := dir + name

	if !fileUtils.FileExist(expectIndependentFile) {
		expectIndependentFile = dir + "." + name
	}

	if fileUtils.FileExist(expectIndependentFile) {
		expectIndependentContent := fileUtils.ReadFile(expectIndependentFile)
		return true, expectIndependentContent
	}

	return false, ""
}
