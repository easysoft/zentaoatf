package main

import (
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func (s *ProductApiSuite) BeforeEach(t provider.T) {
	t.ID("1")
	t.AddSubSuite("ProductApi")
}

func (s *ProductApiSuite) TestMe(t provider.T) {
	t.Require().Equal("Success", "Success")
}

type ProductApiSuite struct {
	suite.Suite
}
