package main

import (
	"fmt"
	httpUtils "github.com/easysoft/zentaoatf/pkg/lib/http"
	constTestHelper "github.com/easysoft/zentaoatf/test/helper/conf"
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

func (s *ProductApiSuite) TestProductApi(t provider.T) {
	url := fmt.Sprintf("http://127.0.0.1:%d/api.php/v1/products", constTestHelper.ZentaoPort)
	httpUtils.Get(url)

	t.WithNewStep("My Second Step!", func(sCtx provider.StepCtx) {
		sCtx.NewStep("My First SubStep!")
	})
}
