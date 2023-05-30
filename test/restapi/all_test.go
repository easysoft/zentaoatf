package main

import (
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"testing"
)

type AllSuite struct {
	suite.Suite
}

func TestRestApiSet(t *testing.T) {
	suite.RunSuite(t, new(ProductApiSuite))
}

func (s *AllSuite) BeforeEach(t provider.T) {
	t.ID("0")
	t.AddSubSuite("TestProductApi")
}
