package service

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
)

type ProductService struct {
	ProductRepo *repo.ProductRepo `inject:""`
}

func NewProductService() *ProductService {
	return &ProductService{}
}

func (s *ProductService) Paginate(req serverDomain.ReqPaginate) (ret domain.PageData, err error) {

	ret, err = s.ProductRepo.Paginate(req)

	if err != nil {
		return
	}

	return
}

func (s *ProductService) FindById(id uint) (model.Product, error) {
	return s.ProductRepo.FindById(id)
}

func (s *ProductService) Create(product model.Product) (uint, error) {
	return s.ProductRepo.Create(product)
}

func (s *ProductService) Update(id uint, product model.Product) error {
	return s.ProductRepo.Update(id, product)
}

func (s *ProductService) DeleteById(id uint) error {
	return s.ProductRepo.BatchDelete(id)
}
