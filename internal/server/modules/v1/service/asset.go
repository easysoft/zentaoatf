package service

type AssetService struct {
}

func NewAssetService() *AssetService {
	return &AssetService{}
}

func (s *AssetService) LoadScripts(dir string) (err error) {
	return // s.AssetRepo.DeleteById(id)
}
