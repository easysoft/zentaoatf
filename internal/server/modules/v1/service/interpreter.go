package service

import (
	"errors"
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	langHelper "github.com/aaronchen2k/deeptest/internal/comm/helper/lang"
	commonUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	shellUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/shell"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
	"strings"
)

type InterpreterService struct {
	InterpreterRepo *repo.InterpreterRepo `inject:""`
}

func NewInterpreterService() *InterpreterService {
	return &InterpreterService{}
}

func (s *InterpreterService) List() (ret []model.Interpreter, err error) {
	ret, err = s.InterpreterRepo.List()
	return
}

func (s *InterpreterService) Get(id uint) (site model.Interpreter, err error) {
	return s.InterpreterRepo.Get(id)
}

func (s *InterpreterService) Create(site model.Interpreter) (id uint, err error) {
	if !fileUtils.FileExist(site.Path) {
		err = errors.New(fmt.Sprintf("可执行文件%s不存在", site.Path))
		return
	}

	id, err = s.InterpreterRepo.Create(site)
	return
}

func (s *InterpreterService) Update(site model.Interpreter) (err error) {
	if !fileUtils.FileExist(site.Path) {
		err = errors.New(fmt.Sprintf("可执行文件%s不存在", site.Path))
		return
	}

	err = s.InterpreterRepo.Update(site)
	return
}

func (s *InterpreterService) Delete(id uint) error {
	return s.InterpreterRepo.Delete(id)
}

func (s *InterpreterService) GetLangSettings() (mp map[string]interface{}, err error) {
	allLangs := langHelper.GetSupportLanguageArrSort()

	langs := []string{}
	mpData := map[string]map[string]string{}
	for _, lang := range allLangs {
		mp := commConsts.LangMap[lang]
		if mp["interpreter"] == "" {
			continue
		}

		subMap := map[string]string{
			"name":        mp["name"],
			"interpreter": mp["interpreter"],
			"versionCmd":  mp["versionCmd"],
		}
		mpData[lang] = subMap
		langs = append(langs, lang)
	}

	mp = map[string]interface{}{}
	mp["languages"] = langs
	mp["languageMap"] = mpData

	return
}

func (s *InterpreterService) GetLangInterpreter(language string) (mp map[string]interface{}, err error) {
	mp = map[string]interface{}{}

	langSettings := commConsts.LangMap[language]
	whereCmd := strings.TrimSpace(langSettings["whereCmd"])
	versionCmd := strings.TrimSpace(langSettings["versionCmd"])

	path := langSettings["interpreter"]
	info := ""

	if !commonUtils.IsWin() || whereCmd == "" {
		return
	}

	output, err := shellUtils.ExeSysCmd(whereCmd)
	if err != nil {
		return
	}

	path = strings.TrimSpace(output)

	info, err = shellUtils.ExeSysCmd(fmt.Sprintf(versionCmd, "\""+path+"\""))
	if err != nil {
		return
	}

	mp["path"] = path
	mp["info"] = info

	return
}
