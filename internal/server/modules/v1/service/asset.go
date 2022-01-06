package service

import (
	"encoding/json"
	commonUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	langUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/lang"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	zentaoUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/zentao"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/kataras/iris/v12"
	"io/ioutil"
	"path/filepath"
	"regexp"
)

type AssetService struct {
}

func NewAssetService() *AssetService {
	return &AssetService{}
}

func (s *AssetService) LoadScripts(dir string) (asset serverDomain.TestAsset, err error) {
	if !fileUtils.FileExist(dir) {
		logUtils.Errorf("dir %s not exist", dir)
		return
	}

	if !commonUtils.IsRelease() { // debug in ide
		dir = filepath.Join(dir, "demo")
	}

	asset = serverDomain.TestAsset{Path: dir, Title: fileUtils.GetDirName(dir), IsDir: true, Slots: iris.Map{"icon": "icon"}}
	s.GetAllScriptsInDir(dir, &asset)

	jsn, _ := json.Marshal(asset)
	logUtils.Infof(string(jsn))

	return
}

func (s *AssetService) GetAllScriptsInDir(childPath string, parent *serverDomain.TestAsset) (err error) {
	if !fileUtils.IsDir(childPath) { // is file
		s.addScript(childPath, parent)
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
		if grandson.IsDir() { // 目录, 递归遍历
			dirNode := s.addDir(childPath, parent)

			s.GetAllScriptsInDir(childPath, dirNode)
		} else {
			s.addScript(childPath, parent)
		}
	}

	return
}

func (s *AssetService) addScript(pth string, parent *serverDomain.TestAsset) {
	regx := langUtils.GetSupportLanguageExtRegx()
	pass, _ := regexp.MatchString("^*.\\."+regx+"$", pth)

	if pass {
		pass = zentaoUtils.CheckFileIsScript(pth)
		if pass {
			childScript := &serverDomain.TestAsset{Path: pth, Title: fileUtils.GetFileName(pth),
				IsDir: false, Slots: iris.Map{"icon": "icon"}}

			parent.Children = append(parent.Children, childScript)
			parent.ScriptCount += 1
		}
	}
}
func (s *AssetService) addDir(pth string, parent *serverDomain.TestAsset) (dirNode *serverDomain.TestAsset) {
	dirNode = &serverDomain.TestAsset{Path: pth, Title: fileUtils.GetDirName(pth),
		IsDir: true, Slots: iris.Map{"icon": "icon"}}
	parent.Children = append(parent.Children, dirNode)

	return
}
