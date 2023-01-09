package service

import (
	"errors"
	"fmt"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"

	langHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/lang"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/repo"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
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

func (s *InterpreterService) Get(id uint) (interpreter model.Interpreter, err error) {
	return s.InterpreterRepo.Get(id)
}

func (s *InterpreterService) Create(interpreter model.Interpreter) (id uint, err error) {
	if !fileUtils.FileExist(interpreter.Path) {
		err = errors.New(i118Utils.Sprintf("wrong_interpreter_format", interpreter.Path))
		return
	}

	id, err = s.InterpreterRepo.Create(interpreter)
	return
}

func (s *InterpreterService) Update(interpreter model.Interpreter) (err error) {
	if !fileUtils.FileExist(interpreter.Path) {
		err = errors.New(fmt.Sprintf("可执行文件%s不存在", interpreter.Path))
		return
	}

	err = s.InterpreterRepo.Update(interpreter)
	return
}

func (s *InterpreterService) Delete(id uint) error {
	return s.InterpreterRepo.Delete(id)
}

func (s *InterpreterService) GetMap(pos []model.Interpreter) (mp map[string]string, err error) {
	mp = map[string]string{}

	for _, item := range pos {
		mp[item.Lang] = item.Path
	}

	return
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
