package service

import (
	"encoding/json"
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	commonUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	httpUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/http"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
	configUtils "github.com/aaronchen2k/deeptest/internal/server/modules/v1/utils/config"
	zentaoUtils "github.com/aaronchen2k/deeptest/internal/server/modules/v1/utils/zentao"
	"github.com/emirpasic/gods/maps"
	"sort"
	"strconv"
	"strings"
)

type ZtfCaseService struct {
	ProjectRepo *repo.ProjectRepo `inject:""`
}

func (s *ZtfCaseService) NewZtfCaseService() *ZtfCaseService {
	return &ZtfCaseService{}
}

func (s *ZtfCaseService) LoadTestCases(productId, moduleId, suiteId, taskId int, projectPath string) (
	cases []commDomain.ZtfCase, loginFail bool) {

	config := configUtils.LoadByProjectPath(projectPath)

	ok := zentaoUtils.Login(config)
	if !ok {
		loginFail = true
		return
	}

	if moduleId != 0 {
		cases = s.ListCaseByModule(config.Url, productId, moduleId)
	} else if suiteId != 0 {
		cases = s.ListCaseBySuite(config.Url, suiteId)
	} else if taskId != 0 {
		cases = s.ListCaseByTask(config.Url, taskId)
	} else if productId != 0 {
		cases = s.ListCaseByProduct(config.Url, productId)
	}

	return
}

func (s *ZtfCaseService) ListCaseByProduct(baseUrl string, productId int) []commDomain.ZtfCase {
	// $productID=productId, $branch = '', $browseType = 'byModule', $param=moduleId,
	// $orderBy='id_desc', $recTotal=0, $recPerPage=10000, $pageID=1)

	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%d--byModule-all-id_asc-0-10000-1", productId)
	} else {
		params = fmt.Sprintf("productID=%d&branch=&browseType=byModule&param=0&orderBy=id_desc&recTotal=0&recPerPage=10000", productId)
	}

	url := baseUrl + zentaoUtils.GenApiUri("testcase", "browse", params)
	dataStr, ok := httpUtils.Get(url)

	if ok {
		var product commDomain.ZtfProduct
		json.Unmarshal(dataStr, &product)

		caseArr := make([]commDomain.ZtfCase, 0)
		for _, cs := range product.Cases {
			caseId := cs.Id

			csWithSteps := s.GetCaseById(baseUrl, caseId)
			stepArr := s.genCaseSteps(csWithSteps)
			caseArr = append(caseArr, commDomain.ZtfCase{Id: caseId, Product: cs.Product, Module: cs.Module,
				Title: cs.Title, StepArr: stepArr})
		}

		return caseArr
	}

	return nil
}

func (s *ZtfCaseService) ListCaseByModule(baseUrl string, productId, moduleId int) []commDomain.ZtfCase {
	// $productID=productId, $branch = '', $browseType = 'byModule', $param=moduleId,
	// $orderBy='id_desc', $recTotal=0, $recPerPage=10000, $pageID=1)

	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%d--byModule-%d-id_asc-0-10000-1", productId, moduleId)
	} else {
		params = fmt.Sprintf("productID=%d&branch=&browseType=byModule&param=%d&orderBy=id_desc&recTotal=0&recPerPage=10000", productId, moduleId)
	}

	url := baseUrl + zentaoUtils.GenApiUri("testcase", "browse", params)
	dataStr, ok := httpUtils.Get(url)

	if ok {
		var module commDomain.ZtfModule
		json.Unmarshal([]byte(dataStr), &module)

		caseArr := make([]commDomain.ZtfCase, 0)
		for _, cs := range module.Cases {
			caseId := cs.Id

			csWithSteps := s.GetCaseById(baseUrl, caseId)
			stepArr := s.genCaseSteps(csWithSteps)

			caseArr = append(caseArr, commDomain.ZtfCase{Id: caseId, Product: cs.Product, Module: cs.Module,
				Title: cs.Title, StepArr: stepArr})
		}

		return caseArr
	}

	return nil
}

func (s *ZtfCaseService) ListCaseBySuite(baseUrl string, suiteId int) []commDomain.ZtfCase {
	// $suiteID, $orderBy = 'id_desc', $recTotal = 0, $recPerPage = 20, $pageID = 1

	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%d-id_asc-0-10000-1", suiteId)
	} else {
		params = fmt.Sprintf("suiteID=%d&orderBy=id_desc&recTotal=0&recPerPage=10000", suiteId)
	}

	url := baseUrl + zentaoUtils.GenApiUri("testsuite", "view", params)
	dataStr, ok := httpUtils.Get(url)

	if ok {
		var suite commDomain.ZtfSuite
		json.Unmarshal([]byte(dataStr), &suite)

		caseArr := make([]commDomain.ZtfCase, 0)
		for _, cs := range suite.Cases {
			caseId := cs.Id

			csWithSteps := s.GetCaseById(baseUrl, caseId)
			stepArr := s.genCaseSteps(csWithSteps)

			// check cs.Suite -> module
			caseArr = append(caseArr, commDomain.ZtfCase{Id: caseId, Product: cs.Product, Module: cs.Suite,
				Title: cs.Title, StepArr: stepArr})
		}

		return caseArr
	}

	return nil
}

func (s *ZtfCaseService) ListCaseByTask(baseUrl string, taskId int) []commDomain.ZtfCase {
	// $taskID, $browseType = 'all', $param = 0,
	// $orderBy = 'id_desc', $recTotal = 0, $recPerPage = 20, $pageID = 1

	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%d-all-0-id_asc-0-10000-1", taskId)
	} else {
		params = fmt.Sprintf("taskID=%d&browseType=all&param=0&orderBy=id_desc&recTotal=0&recPerPage=10000", taskId)
	}

	url := baseUrl + zentaoUtils.GenApiUri("testtask", "cases", params)
	dataStr, ok := httpUtils.Get(url)

	if ok {
		var task commDomain.ZtfTask
		json.Unmarshal([]byte(dataStr), &task)

		caseArr := make([]commDomain.ZtfCase, 0)
		for _, item := range task.Runs {
			caseId := item.Case

			csWithSteps := s.GetCaseById(baseUrl, caseId)
			stepArr := s.genCaseSteps(csWithSteps)

			caseArr = append(caseArr, commDomain.ZtfCase{Id: caseId, Product: item.Product, Module: item.Module,
				Title: item.Title, StepArr: stepArr})
		}

		return caseArr
	}

	return nil
}

func (s *ZtfCaseService) genCaseSteps(csWithSteps commDomain.ZtfCase) (ret []commDomain.ZtfStep) {
	// get order keys
	keys := make([]int, 0, len(csWithSteps.Steps))
	for k := range csWithSteps.Steps {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, key := range keys {
		step := csWithSteps.Steps[key]
		ret = append(ret, step)
	}

	return
}

func (s *ZtfCaseService) GetCaseById(baseUrl string, caseId string) commDomain.ZtfCase {
	// $caseID, $version = 0, $from = 'testcase', $taskID = 0

	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%s-0-testcase-0", caseId)
	} else {
		params = fmt.Sprintf("caseID=%s&version=0&$from=testcase&taskID=0", caseId)
	}

	url := baseUrl + zentaoUtils.GenApiUri("testcase", "view", params)
	dataStr, ok := httpUtils.Get(url)

	if ok {
		var csw commDomain.ZtfCaseWrapper
		json.Unmarshal([]byte(dataStr), &csw)

		cs := csw.Case
		return cs
	}

	return commDomain.ZtfCase{}
}

func (s *ZtfCaseService) GetCaseIdsBySuite(suiteId int, idMap *map[int]string, projectPath string) {
	config := configUtils.LoadByProjectPath(projectPath)

	ok := zentaoUtils.Login(config)
	if !ok {
		return
	}

	testcases := s.ListCaseBySuite(config.Url, suiteId)

	for _, tc := range testcases {
		id, _ := strconv.Atoi(tc.Id)
		(*idMap)[id] = ""
	}
}

func (s *ZtfCaseService) GetCaseIdsByTask(taskId int, idMap *map[int]string, projectPath string) {
	config := configUtils.LoadByProjectPath(projectPath)

	ok := zentaoUtils.Login(config)
	if !ok {
		return
	}

	testcases := s.ListCaseByTask(config.Url, taskId)

	for _, tc := range testcases {
		id, _ := strconv.Atoi(tc.Id)
		(*idMap)[id] = ""
	}
}

func (s *ZtfCaseService) CommitCase(caseId int, title string,
	stepMap maps.Map, stepTypeMap maps.Map, expectMap maps.Map, projectPath string) {
	config := configUtils.LoadByProjectPath(projectPath)

	ok := zentaoUtils.Login(config)
	if !ok {
		return
	}

	// $caseID, $comment = false
	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%d-0", caseId)
	} else {
		params = fmt.Sprintf("caseID=%d&comment=0", caseId)
	}

	url := config.Url + zentaoUtils.GenApiUri("testcase", "edit", params)

	requestObj := map[string]interface{}{"title": title,
		"steps":    commonUtils.LinkedMapToMap(stepMap),
		"stepType": commonUtils.LinkedMapToMap(stepTypeMap),
		"expects":  commonUtils.LinkedMapToMap(expectMap)}

	json, _ := json.Marshal(requestObj)

	if commConsts.Verbose {
		logUtils.Infof(string(json))
	}

	_, ok = httpUtils.Post(url, requestObj, true)
	if ok {
		logUtils.Infof(i118Utils.Sprintf("success_to_commit_case", caseId) + "\n")
	}
}

func (s *ZtfCaseService) addPrefixSpace(str string, numb int) string {
	arr := strings.Split(str, "\r\n")

	ret := make([]string, 0)
	for _, line := range arr {
		line = fmt.Sprintf("%*s", numb, " ") + line

		ret = append(ret, line)
	}

	return strings.Join(ret, "\n")
}
