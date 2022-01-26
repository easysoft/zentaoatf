package scriptUtils

import (
	"encoding/json"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	commonUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	langUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/lang"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/kataras/iris/v12"
	"io/ioutil"
	"regexp"
)

func LoadScriptTree(dir string) (asset serverDomain.TestAsset, err error) {
	if !fileUtils.FileExist(dir) {
		logUtils.Errorf("dir %s not exist", dir)
		return
	}

	commonUtils.ChangeScriptForDebug(&dir)

	asset = serverDomain.TestAsset{Path: dir, Title: fileUtils.GetDirName(dir), IsDir: true, Slots: iris.Map{"icon": "icon"}}
	LoadScriptNodesInDir(dir, &asset, 0)

	jsn, _ := json.Marshal(asset)
	logUtils.Infof(string(jsn))

	return
}

func LoadScriptByProject(projectPath string) (scriptFiles []string) {
	LoadScriptListInDir(projectPath, &scriptFiles, 0)

	return
}

func GetScriptContent(pth string) (script model.TestScript, err error) {
	script.Code = fileUtils.ReadFile(pth)

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
