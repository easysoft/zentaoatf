package service

import (
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
)

type SyncService struct {
}

func NewSyncService() *SyncService {
	return &SyncService{}
}

func (s *SyncService) SyncFromZentao(settings commDomain.SyncSettings) (err error) {

	return
}

func (s *SyncService) SyncToZentao() (err error) {

	return
}
