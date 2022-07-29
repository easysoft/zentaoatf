package service

import (
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/repo"
)

type StatisticService struct {
	StatisticRepo *repo.StatisticRepo `inject:""`
}

func NewStatisticService() *StatisticService {
	return &StatisticService{}
}

func (s *StatisticService) Get(id uint) (statistics model.Statistic, err error) {
	return s.StatisticRepo.Get(id)
}

func (s *StatisticService) GetByPath(path string) (statistics model.Statistic, err error) {
	return s.StatisticRepo.GetByPath(path)
}

func (s *StatisticService) Create(statistics model.Statistic) (id uint, isDuplicate bool, err error) {
	id, isDuplicate, err = s.StatisticRepo.Create(&statistics)

	return
}

func (s *StatisticService) Update(statistics model.Statistic) (isDuplicate bool, err error) {
	isDuplicate, err = s.StatisticRepo.Update(statistics)
	if isDuplicate || err != nil {
		return
	}
	return
}

func (s *StatisticService) Delete(id uint) error {
	return s.StatisticRepo.Delete(id)
}
