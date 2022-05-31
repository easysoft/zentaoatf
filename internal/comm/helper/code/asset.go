package codeHelper

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	commonUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/common"
	fileUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/file"
	logUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/log"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	"github.com/kataras/iris/v12"
)

func LoadCodeTree(workspace model.Workspace) (asset serverDomain.TestAsset, err error) {
	workspaceId := int(workspace.ID)
	workspaceDir := workspace.Path

	if !fileUtils.FileExist(workspaceDir) {
		logUtils.Errorf("workspaceDir %s not exist", workspaceDir)
		return
	}

	asset = serverDomain.TestAsset{
		Type:          commConsts.Workspace,
		WorkspaceId:   workspaceId,
		WorkspaceType: workspace.Type,
		Path:          workspaceDir,
		Title:         fileUtils.GetDirName(workspaceDir),
		Slots:         iris.Map{"icon": "icon"},

		Checkable: false,
		IsLeaf:    false,
	}

	nodes, err := LoadCodeNodesInDir(workspaceDir, workspaceId, workspace.Type)
	asset.Children = append(asset.Children, nodes...)

	jsn, _ := json.Marshal(asset)
	logUtils.Infof(string(jsn))

	return
}

func LoadCodeNodesInDir(dir string, workspaceId int, workspaceType commConsts.TestTool) (
	nodes []*serverDomain.TestAsset, err error) {
	list, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	for _, child := range list {
		name := child.Name()
		if commonUtils.IgnoreCodeFile(name) {
			continue
		}

		childPath := filepath.Join(dir, name)
		if child.IsDir() { // 目录
			dirNode := getDir(childPath, workspaceId, workspaceType)
			dirNode.Children, _ = LoadCodeNodesInDir(dirNode.Path, workspaceId, workspaceType)
			nodes = append(nodes, &dirNode)
		} else {
			fileNode := getFile(childPath, workspaceId, workspaceType)
			nodes = append(nodes, &fileNode)
		}
	}

	return
}

func getFile(pth string, workspaceId int, workspaceType commConsts.TestTool) (fileNode serverDomain.TestAsset) {
	fileNode = serverDomain.TestAsset{
		Type:          commConsts.File,
		WorkspaceId:   workspaceId,
		WorkspaceType: workspaceType,
		Path:          pth,
		Title:         fileUtils.GetFileName(pth),
		Slots:         iris.Map{"icon": "icon"},

		Checkable: false,
		IsLeaf:    true,
	}

	return
}

func getDir(pth string, workspaceId int, workspaceType commConsts.TestTool) (dirNode serverDomain.TestAsset) {
	dirNode = serverDomain.TestAsset{
		Type:          commConsts.Dir,
		WorkspaceId:   workspaceId,
		WorkspaceType: workspaceType,
		Path:          pth,
		Title:         fileUtils.GetDirName(pth),

		Checkable: false,
		Slots:     iris.Map{"icon": "icon"}}

	return
}
