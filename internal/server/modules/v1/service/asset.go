package service

import serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"

type AssetService struct {
}

func NewAssetService() *AssetService {
	return &AssetService{}
}

func (s *AssetService) LoadScripts(dir string) (asset serverDomain.TestAsset, err error) {
	return
}
