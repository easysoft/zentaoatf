package service

import (
	"os"
	"path/filepath"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	codeHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/code"
	configHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/config"
	scriptHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/script"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/repo"
	"github.com/easysoft/zentaoatf/pkg/domain"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	"github.com/kataras/iris/v12"
)

type TestScriptService struct {
	WorkspaceRepo     *repo.WorkspaceRepo `inject:""`
	SiteService       *SiteService        `inject:""`
	TestResultService *TestResultService  `inject:""`
}

func NewTestScriptService() *TestScriptService {
	return &TestScriptService{}
}

func (s *TestScriptService) LoadTestScriptsBySiteProduct(
	siteId, productId uint, displayBy, filerType string, filerValue int) (root serverDomain.TestAsset, err error) {

	scriptIdsFromZentao := s.getScriptIdsFromZentao(siteId, productId, filerType, filerValue)
	workspaces, _ := s.WorkspaceRepo.ListByProduct(siteId, productId)

	// load scripts from disk
	root = serverDomain.TestAsset{Path: "", Title: "all", Type: commConsts.Root, Slots: iris.Map{"icon": "icon"}}
	for _, workspace := range workspaces {

		if filerType == string(commConsts.FilterWorkspace) &&
			(filerValue > 0 && uint(filerValue) != workspace.ID) { // filter by workspace
			continue
		}

		scriptsInDir := s.loadTestScriptsInWorkspace(displayBy, filerType, workspace, scriptIdsFromZentao, productId, siteId, filerValue)

		if scriptsInDir.Title != "" {
			root.Children = append(root.Children, &scriptsInDir)
		}
	}

	return
}

func (s *TestScriptService) loadTestScriptsInWorkspace(displayBy, filerType string, workspace model.Workspace, scriptIdsFromZentao map[int]string, productId, siteId uint, filerValue int) (scriptsInDir serverDomain.TestAsset) {
	if displayBy != "module" && workspace.Type != commConsts.ZTF {
		scriptsInDir, _ = codeHelper.LoadCodeTree(workspace)
		return
	}

	if displayBy == "workspace" {
		scriptsInDir, _ = scriptHelper.LoadScriptTreeByDir(workspace, scriptIdsFromZentao)
		return
	}

	//display by module
	site, _ := s.SiteService.Get(siteId)
	config := configHelper.LoadBySite(site)

	suiteId := 0
	taskId := 0
	if filerType == string(commConsts.FilterSuite) {
		suiteId = filerValue
	} else if filerType == string(commConsts.FilterTask) {
		taskId = filerValue
	}
	scriptsInDir, _ = zentaoHelper.LoadTestCasesInModuleTree(workspace, scriptIdsFromZentao,
		int(productId), suiteId, taskId, config)

	return
}

func (s *TestScriptService) LoadCodeChildren(dir string, workspaceId int) (nodes []*serverDomain.TestAsset, err error) {
	workspace, _ := s.WorkspaceRepo.Get(uint(workspaceId))

	nodes, _ = codeHelper.LoadCodeNodesInDir(dir, workspaceId, workspace.Type)

	return
}

func (s *TestScriptService) getScriptIdsFromZentao(siteId, productId uint, filerType string, filerValue int) (
	ret map[int]string) {

	if filerType == "" || filerValue < 0 {
		return nil
	}

	// get script ids from zentao
	currSite, _ := s.SiteService.Get(uint(siteId))
	config := commDomain.WorkspaceConf{
		Url:      currSite.Url,
		Username: currSite.Username,
		Password: currSite.Password,
	}

	if filerType == string(commConsts.FilterModule) {
		ret, _ = zentaoHelper.GetCaseIdsInZentaoModule(productId, uint(filerValue), config)
	} else if filerType == string(commConsts.FilterSuite) {
		ret, _ = zentaoHelper.GetCaseIdsInZentaoSuite(productId, filerValue, config)
	} else if filerType == string(commConsts.FilterTask) {
		ret, _ = zentaoHelper.GetCaseIdsInZentaoTask(productId, filerValue, config)
	}

	return
}

func (s *TestScriptService) GetCaseIdsFromReport(workspaceId int, seq, scope string) (caseIds []string, err error) {
	report, err := s.TestResultService.Get(workspaceId, seq)

	if err != nil {
		return
	}

	for _, cs := range report.FuncResult {
		path := cs.Path
		status := cs.Status

		if path != "" && (scope == "all" || scope == status.String()) {
			caseIds = append(caseIds, path)
		}
	}

	return
}

func (s *TestScriptService) CreateNode(req serverDomain.CreateScriptReq) (pth string, err *domain.BizError) {
	name := req.Name
	extName := fileUtils.GetExtNameWithoutDot(name)
	mode := req.Mode
	typ := req.Type
	target := req.Target
	//workspaceId := req.WorkspaceId
	productId := req.ProductId

	if !fileUtils.IsDir(target) && mode == commConsts.Child {
		mode = commConsts.Brother
	}

	dir := ""
	if mode == commConsts.Child {
		dir = target
	} else if mode == commConsts.Brother {
		dir = filepath.Dir(target)
	}

	pth = filepath.Join(dir, name)
	if fileUtils.FileExist(pth) {
		err = &domain.BizError{Code: commConsts.ErrFileOrDirExist.Code}
		return
	}

	if typ == commConsts.CreateDir {
		fileUtils.MkDirIfNeeded(pth)
	} else {
		if extName == "" {
			fileUtils.WriteFile(pth, "")
			return
		}

		lang := commConsts.ScriptExtToNameMap[extName]
		if lang == "" {
			fileUtils.WriteFile(pth, "")
			return
		}

		scriptHelper.GenEmptyScript(name, lang, pth, productId)
	}

	return
}

func (s *TestScriptService) UpdateCode(script serverDomain.TestScript) (err error) {
	fileUtils.WriteFile(script.Path, script.Code)

	return
}

func (s *TestScriptService) UpdateName(script serverDomain.TestScript) (err error) {
	dir := filepath.Dir(script.Path)
	newPath := filepath.Join(dir, script.Name)
	os.Rename(script.Path, newPath)

	return
}

func (s *TestScriptService) Delete(pth string) (bizErr *domain.BizError) {
	err := os.RemoveAll(pth)

	_, ok := err.(*os.PathError)
	if ok {
		bizErr = &domain.BizError{Code: commConsts.ErrDirNotEmpty.Code}
	} else {
		bizErr = nil
	}

	return
}

func (s *TestScriptService) Rename(pth string, name string) (bizErr *domain.BizError) {
	dir, _ := filepath.Split(pth)
	err := os.Rename(pth, dir+name)

	_, ok := err.(*os.PathError)
	if ok {
		bizErr = &domain.BizError{Code: commConsts.ErrDirNotEmpty.Code}
	} else {
		bizErr = nil
	}

	return
}

func (s *TestScriptService) Move(req serverDomain.MoveScriptReq) (err error) {
	src := req.DragKey
	srcName := fileUtils.GetFileName(src)
	dist := req.DropKey

	distDir := ""
	if req.DropPosition == commConsts.Inner && fileUtils.IsDir(dist) {
		distDir = dist
	} else {
		distDir = filepath.Dir(dist)
	}

	pth := filepath.Join(distDir, srcName)
	err = os.Rename(src, pth)

	return
}

func (s *TestScriptService) Paste(req serverDomain.PasteScriptReq) (err error) {
	distPath := filepath.Join(req.DistKey, filepath.Base(req.SrcKey))

	if distPath == req.SrcKey {
		return
	}

	if req.Action == "copy" {
		if req.SrcType == commConsts.File {
			fileUtils.CopyFile(req.SrcKey, distPath)
		} else {
			fileUtils.CopyDir(req.SrcKey, distPath)
		}
	} else if req.Action == "cut" {
		os.Rename(req.SrcKey, distPath)
	}

	return
}
