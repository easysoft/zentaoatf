package service

import (
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/comm/domain"
	codeHelper "github.com/easysoft/zentaoatf/internal/comm/helper/code"
	configHelper "github.com/easysoft/zentaoatf/internal/comm/helper/config"
	scriptHelper "github.com/easysoft/zentaoatf/internal/comm/helper/script"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/comm/helper/zentao"
	"github.com/easysoft/zentaoatf/internal/pkg/domain"
	fileUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/file"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/repo"
	"github.com/kataras/iris/v12"
	"os"
	"path/filepath"
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

		var scriptsInDir serverDomain.TestAsset

		if workspace.Type == commConsts.ZTF {
			if displayBy == "workspace" {
				scriptsInDir, _ = scriptHelper.LoadScriptTreeByDir(workspace, scriptIdsFromZentao)

			} else if displayBy == "module" {
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
			}

		} else if displayBy != "module" && workspace.Type != commConsts.ZTF {
			scriptsInDir, _ = codeHelper.LoadCodeTree(workspace)
		}

		if scriptsInDir.Title != "" {
			root.Children = append(root.Children, &scriptsInDir)
		}
	}

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

func (s *TestScriptService) CreateNode(req serverDomain.CreateScriptReq) (pth string, err error) {
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
	if typ == commConsts.CreateDir {
		fileUtils.MkDirIfNeeded(pth)
	} else {
		if extName == "" {
			fileUtils.WriteFile(pth, "")
			return
		}

		lang := commConsts.ScriptExtToNameMap[fileUtils.GetExtNameWithoutDot(pth)]
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
	err := os.Remove(pth)

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
