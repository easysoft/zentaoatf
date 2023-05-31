package main

import (
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	constTestHelper "github.com/easysoft/zentaoatf/test/helper/conf"
	httpHelper "github.com/easysoft/zentaoatf/test/helper/http"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/tidwall/gjson"
	"testing"
)

func TestProductApi(t *testing.T) {
	suite.RunSuite(t, new(ProductApiSuite))
}

type ProductApiSuite struct {
	suite.Suite
}

func (s *ProductApiSuite) BeforeEach(t provider.T) {
	t.ID("1,2")
	t.AddSubSuite("ProductApi")
}

func (s *ProductApiSuite) TestProductListApi(t provider.T) {
	token := httpHelper.Login()

	url := zentaoHelper.GenApiUrl("products", nil, constTestHelper.ZentaoSiteUrl)

	bodyBytes, _ := httpHelper.Get(url, token)

	firstProductId := gjson.Get(string(bodyBytes), "products.0.id").Int()

	t.Require().Greater(firstProductId, int64(0), "list product")
}
