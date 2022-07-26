package service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/easysoft/zentaoatf/pkg/domain"
	commonUtils "github.com/easysoft/zentaoatf/pkg/lib/common"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/helper/dir"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	configHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/config"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/repo"
)

type WorkspaceService struct {
	WorkspaceRepo      *repo.WorkspaceRepo `inject:""`
	SiteRepo           *repo.SiteRepo      `inject:""`
	InterpreterService *InterpreterService `inject:""`
	ProxyService       *ProxyService       `inject:""`
}

func NewWorkspaceService() *WorkspaceService {
	return &WorkspaceService{}
}

func (s *WorkspaceService) Paginate(req serverDomain.WorkspaceReqPaginate) (ret domain.PageData, err error) {
	ret, err = s.WorkspaceRepo.Paginate(req)
	return
}

func (s *WorkspaceService) ListWorkspacesByProduct(siteId, productId uint) (pos []model.Workspace, err error) {
	return s.WorkspaceRepo.ListByProduct(siteId, productId)
}

func (s *WorkspaceService) Get(id uint) (model.Workspace, error) {
	return s.WorkspaceRepo.Get(id)
}
func (s *WorkspaceService) GetByPath(workspacePath string) (po model.Workspace, err error) {
	return s.WorkspaceRepo.FindByPath(workspacePath)
}

func (s *WorkspaceService) Create(workspace model.Workspace) (id uint, err error) {
	if !fileUtils.IsDir(workspace.Path) {
		err = errors.New(fmt.Sprintf("目录%s不存在", workspace.Path))
		return
	}

	id, err = s.WorkspaceRepo.Create(workspace)
	s.UpdateConfig(workspace, "all")

	return
}
func (s *WorkspaceService) Update(workspace model.Workspace) (err error) {
	if !fileUtils.IsDir(workspace.Path) {
		err = errors.New(fmt.Sprintf("目录%s不存在", workspace.Path))
		return
	}

	err = s.WorkspaceRepo.Update(workspace)
	s.UpdateConfig(workspace, "all")

	return
}

func (s *WorkspaceService) Delete(id uint) (err error) {
	err = s.WorkspaceRepo.Delete(id)
	if err != nil {
		return
	}

	err = s.WorkspaceRepo.SetCurrWorkspace("")

	return
}

func (s *WorkspaceService) DeleteByPath(path string, productId uint) (err error) {
	err = s.WorkspaceRepo.DeleteByPath(path, productId)
	if err != nil {
		return
	}

	err = s.WorkspaceRepo.SetCurrWorkspace("")

	return
}

func (s *WorkspaceService) ListByProduct(siteId, productId uint) (pos []model.Workspace, err error) {
	return s.WorkspaceRepo.ListByProduct(siteId, productId)
}

func (s *WorkspaceService) UpdateAllConfig() {
	workspaces, _ := s.WorkspaceRepo.ListWorkspace()

	for _, item := range workspaces {
		if item.Type != commConsts.ZTF {
			continue
		}

		s.UpdateConfig(item, "interpreter")
	}
}

func (s *WorkspaceService) UpdateConfig(workspace model.Workspace, by string) (err error) {
	site, _ := s.SiteRepo.Get(workspace.SiteId)
	interps, _ := s.InterpreterService.List()
	mp, _ := s.InterpreterService.GetMap(interps)

	conf := configHelper.ReadFromFile(workspace.Path)
	if conf.Language == "" {
		conf.Language = commConsts.LanguageZh
	}

	if by == "all" || by == "site" {
		conf.Url = site.Url
		conf.Username = site.Username
		conf.Password = site.Password
	}

	if by == "all" || by == "interpreter" {
		conf.Javascript = mp["javascript"]
		conf.Lua = mp["lua"]
		conf.Perl = mp["perl"]
		conf.Php = mp["php"]
		conf.Python = mp["python"]
		conf.Ruby = mp["ruby"]
		conf.Tcl = mp["tcl"]
		conf.Autoit = mp["autoit"]
	}

	configHelper.SaveToFile(conf, workspace.Path)

	return
}

func (s *WorkspaceService) UploadScriptsToProxy(testSets []serverDomain.TestSet) (pathMap map[string]string, err error) {
	pathMap = make(map[string]string)
	unitResultPath := filepath.Join(commConsts.WorkDir, commConsts.ExecZip)
	uploadDir := filepath.Join(commConsts.WorkDir, commConsts.ExecZipPath)
	os.RemoveAll(uploadDir)
	os.Remove(unitResultPath)
	var workspaceInfo model.Workspace
	var workspacePathArray = []string{}
	var workspaceIsZtfMap = make(map[string]bool)
	for index, testSet := range testSets {
		po, _ := s.Get(uint(testSet.WorkspaceId))
		testSets[index].WorkspacePath = po.Path
		if workspaceInfo.ID == 0 && po.ID > 0 {
			workspaceInfo = po
		}
		workspacePathArray = append(workspacePathArray, po.Path)
		workspaceIsZtfMap[po.Path] = po.Type == "ztf"

		if po.Type == "ztf" {
			for _, casePath := range testSet.Cases {
				if fileUtils.IsDir(casePath) {
					continue
				}
				_, err = fileUtils.CopyFileAll(casePath, strings.Replace(casePath, po.Path, uploadDir, 1))
				if err != nil {
					return
				}
			}
		} else {
			err = fileUtils.CopyDir(po.Path, uploadDir)
			if err != nil {
				return
			}
		}
	}
	var uploadUrl = ""
	if workspaceInfo.ProxyId > 0 {
		proxyInfo, _ := s.ProxyService.Get(workspaceInfo.ProxyId)
		if proxyInfo.Path != "" {
			uploadUrl = proxyInfo.Path
		}
	}
	if uploadUrl == "" {
		return
	}
	fileUtils.ZipDir(unitResultPath, uploadDir)
	resp := fileUtils.UploadWithResp(uploadUrl+"api/v1/workspaces/uploadScripts", []string{unitResultPath}, nil)
	dataMap := resp["data"].(map[string]interface{})
	proxyWorkDir := dataMap["workDir"].(string)
	proxySep := dataMap["sep"].(string)
	realScriptDir := proxyWorkDir + commConsts.ExecProxyPath + proxySep + commConsts.ExecZipPath

	for _, testSet := range testSets {
		if workspaceIsZtfMap[testSet.WorkspacePath] {
			workspacePath := testSet.WorkspacePath
			for _, casePath := range testSet.Cases {
				oldCasePath := casePath
				if commConsts.PthSep != proxySep {
					workspacePath = strings.Replace(workspacePath, commConsts.PthSep, proxySep, -1)
					casePath = strings.Replace(casePath, commConsts.PthSep, proxySep, -1)
				}
				pathMap[strings.Replace(casePath, workspacePath, realScriptDir, 1)] = oldCasePath
			}
		} else {
			if commConsts.PthSep != proxySep {
				pathMap[realScriptDir] = testSet.WorkspacePath
			}
		}
	}
	return
}

func (s *WorkspaceService) UploadScripts(fh *multipart.FileHeader, ctx iris.Context) (err error) {
	path := filepath.Join(commConsts.WorkDir, commConsts.ExecProxyPath)
	err = dir.InsureDir(path)
	if err != nil {
		logUtils.Infof(color.RedString("file upload failed, error: %s.", err.Error()))
		return
	}
	zipPath := filepath.Join(path, commConsts.ExecZip)
	_, err = ctx.SaveFormFile(fh, zipPath)
	if err != nil {
		logUtils.Infof(color.RedString("file upload failed, error: %s.", err.Error()))
		return
	}

	fileUtils.Unzip(zipPath, path)
	return
}

func (s *WorkspaceService) GetValidProxy(id uint) (proxyInfo model.Proxy, err error) {
	workspaceInfo, err := s.WorkspaceRepo.Get(id)
	if err != nil {
		return proxyInfo, err
	}
	proxies := strings.Split(workspaceInfo.Proxies, ",")
	rwMap := commonUtils.NewRWMap(len(proxies))
	rwMap.Set(0, 1)
	wg := &sync.WaitGroup{}
	proxyIds := []int{}
	proxyMap := make(map[int]model.Proxy)
	proxyMap[0] = model.Proxy{
		Path: "local",
	}
	for _, proxy := range proxies {
		proxyId, err := strconv.Atoi(proxy)
		if err != nil {
			continue
		}
		proxyIds = append(proxyIds, proxyId)
		if proxyId == 0 {
			continue
		}
		proxyInfo, err = s.ProxyService.Get(uint(proxyId))
		if err != nil {
			continue
		}
		proxyMap[proxyId] = proxyInfo
		wg.Add(1)
		go func(group *sync.WaitGroup, proxyInfo model.Proxy) {
			err := s.ProxyService.CheckServer(proxyInfo.Path)
			if err != nil {
				rwMap.Set(int(proxyInfo.ID), 0)
			} else {
				rwMap.Set(int(proxyInfo.ID), 1)
			}
			group.Done()
		}(wg, proxyInfo)
	}
	wg.Wait()
	for _, proxyId := range proxyIds {
		if status, ok := rwMap.Get(proxyId); ok && status == 1 {
			if proxyId > 0 {
				_ = s.WorkspaceRepo.SetProxyId(id, uint(proxyId))
			}
			return proxyMap[proxyId], nil
		}
	}
	return proxyMap[0], err
}
