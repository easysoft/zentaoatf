package service

import (
	"errors"
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	langHelper "github.com/aaronchen2k/deeptest/internal/comm/helper/lang"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
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
	langs := langHelper.GetSupportLanguageArrSort()

	mpData := map[string]map[string]string{}
	for _, lang := range langs {
		subMap := map[string]string{
			"name":        commConsts.LangMap[lang]["name"],
			"interpreter": commConsts.LangMap[lang]["interpreter"],
			"versionCmd":  commConsts.LangMap[lang]["versionCmd"],
		}
		mpData[lang] = subMap
	}

	mp = map[string]interface{}{}
	mp["languages"] = langs
	mp["languageMap"] = mpData

	return
}
