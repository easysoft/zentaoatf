package scriptHelper

import (
	"encoding/json"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	commonUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	langUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/lang"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	stringUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/string"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/kataras/iris/v12"
	"io/ioutil"
	"path"
	"regexp"
	"strings"
)

func LoadScriptTree(dir string) (asset serverDomain.TestAsset, err error) {
	if !fileUtils.FileExist(dir) {
		logUtils.Errorf("dir %s not exist", dir)
		return
	}

	//commonUtils.ChangeScriptForDebug(&dir)

	asset = serverDomain.TestAsset{Path: dir, Title: fileUtils.GetDirName(dir), IsDir: true, Slots: iris.Map{"icon": "icon"}}
	LoadScriptNodesInDir(dir, &asset, 0)

	jsn, _ := json.Marshal(asset)
	logUtils.Infof(string(jsn))

	return
}

func LoadScriptByWorkspace(workspacePath string) (scriptFiles []string) {
	LoadScriptListInDir(workspacePath, &scriptFiles, 0)

	return
}

func GetScriptContent(pth string) (script serverDomain.TestScript, err error) {
	script.Code = fileUtils.ReadFile(pth)
	script.Lang = getScriptLang(pth)

	return
}
func getScriptLang(pth string) (lang string) {
	extName := strings.TrimLeft(fileUtils.GetExtName(pth), ".")
	for key, val := range langUtils.LangMap {
		if extName == val["extName"] {
			lang = key
			return
		}
	}

	return
}

func LoadScriptNodesInDir(childPath string, parent *serverDomain.TestAsset, level int) (err error) {
	if !fileUtils.IsDir(childPath) { // is file
		addScript(childPath, parent)
		return
	}

	childPath = fileUtils.AddPathSepIfNeeded(fileUtils.AbsolutePath(childPath))

	list, err := ioutil.ReadDir(childPath)
	if err != nil {
		return err
	}

	for _, grandson := range list {
		name := grandson.Name()
		if commonUtils.IgnoreFile(name) {
			continue
		}

		childPath := childPath + name
		if grandson.IsDir() && level < 3 { // 目录, 递归遍历
			dirNode := addDir(childPath, parent)

			LoadScriptNodesInDir(childPath, dirNode, level+1)
		} else {
			addScript(childPath, parent)
		}
	}

	return
}

func LoadScriptListInDir(path string, files *[]string, level int) error {
	regx := langUtils.GetSupportLanguageExtRegx()

	if !fileUtils.IsDir(path) { // first call, param is file
		pass, _ := regexp.MatchString(`.*\.`+regx+`$`, path)
		if pass {
			pass = CheckFileIsScript(path)
			if pass {
				*files = append(*files, path)
			}
		}

		return nil
	}

	path = fileUtils.AbsolutePath(path)

	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, fi := range dir {
		name := fi.Name()
		if commonUtils.IgnoreFile(name) {
			continue
		}

		if fi.IsDir() && level < 3 { // 目录, 递归遍历
			LoadScriptListInDir(path+name+consts.PthSep, files, level+1)
		} else {
			path := path + name
			pass, _ := regexp.MatchString("^*.\\."+regx+"$", path)

			if pass {
				pass = CheckFileIsScript(path)
				if pass {
					*files = append(*files, path)
				}
			}
		}
	}

	return nil
}

func addScript(pth string, parent *serverDomain.TestAsset) {
	regx := langUtils.GetSupportLanguageExtRegx()
	pass, _ := regexp.MatchString("^*.\\."+regx+"$", pth)

	if pass {
		pass = CheckFileIsScript(pth)
		if pass {
			childScript := &serverDomain.TestAsset{Path: pth, Title: fileUtils.GetFileName(pth),
				IsDir: false, Slots: iris.Map{"icon": "icon"}}

			parent.Children = append(parent.Children, childScript)
			parent.ScriptCount += 1
		}
	}
}
func addDir(pth string, parent *serverDomain.TestAsset) (dirNode *serverDomain.TestAsset) {
	dirNode = &serverDomain.TestAsset{Path: pth, Title: fileUtils.GetDirName(pth),
		IsDir: true, Slots: iris.Map{"icon": "icon"}}
	parent.Children = append(parent.Children, dirNode)

	return
}

func GetScriptType(scripts []string) []string {
	exts := make([]string, 0)
	for _, script := range scripts {
		ext := path.Ext(script)
		if ext != "" {
			ext = ext[1:]
			name := langUtils.ScriptExtToNameMap[ext]

			if !stringUtils.FindInArr(name, exts) {
				exts = append(exts, name)
			}
		}
	}

	return exts
}

func GetFailedCasesDirectlyFromTestResult(resultFile string) []string {
	cases := make([]string, 0)

	extName := path.Ext(resultFile)

	if extName == "."+consts.ExtNameResult {
		resultFile = strings.Replace(resultFile, extName, "."+consts.ExtNameJson, -1)
	}

	//if vari.ServerWorkspaceDir != "" {
	//	resultFile = vari.ServerWorkspaceDir + resultFile
	//}

	content := fileUtils.ReadFile(resultFile)

	var report commDomain.ZtfReport
	json.Unmarshal([]byte(content), &report)

	for _, cs := range report.FuncResult {
		if cs.Status != "pass" {
			cases = append(cases, cs.Path)
		}
	}

	return cases
}
