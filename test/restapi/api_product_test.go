package main

import (
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"testing"
)

func TestProductApi(t *testing.T) {
	suite.RunSuite(t, new(ProductApiSuite))
}

type ProductApiSuite struct {
	suite.Suite
}

func (s *ProductApiSuite) BeforeEach(t provider.T) {
	t.ID("1")
	t.AddSubSuite("ProductApi")
}

func (s *ProductApiSuite) TestProductApiTest(t provider.T) {
	t.WithNewStep("My Second Step!", func(sCtx provider.StepCtx) {
		sCtx.NewStep("My First SubStep!")
	})
}
