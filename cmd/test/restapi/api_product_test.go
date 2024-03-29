package main

import (
	"fmt"
	"testing"

	commonTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/common"
	constTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/conf"
	httpHelper "github.com/easysoft/zentaoatf/cmd/test/helper/http"
	"github.com/easysoft/zentaoatf/cmd/test/restapi/config"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/tidwall/gjson"
)

func TestProductApi(t *testing.T) {
	suite.RunSuite(t, new(ProductApiSuite))
}

type ProductApiSuite struct {
	suite.Suite
}

func (s *ProductApiSuite) BeforeEach(t provider.T) {
	commonTestHelper.ReplaceLabel(t, "ProductApi")
}

func (s *ProductApiSuite) TestProductListApi(t provider.T) {
	t.ID("7639")
	//t.ID("1")

	token := httpHelper.Login()

	url := zentaoHelper.GenApiUrl("products", nil, constTestHelper.ZentaoSiteUrl)

	bodyBytes, _ := httpHelper.Get(url, token)

	firstProductId := gjson.Get(string(bodyBytes), "products.0.id").Int()

	t.Require().Greater(firstProductId, int64(0), "list product")
}

func (s *ProductApiSuite) TestProductDetailApi(t provider.T) {
	t.ID("7640")
	//t.ID("2")

	token := httpHelper.Login()

	url := zentaoHelper.GenApiUrl(fmt.Sprintf("/products/%d", config.ProductId), nil, constTestHelper.ZentaoSiteUrl)

	bodyBytes, _ := httpHelper.Get(url, token)

	name := gjson.Get(string(bodyBytes), "name").String()

	t.Require().Greater(len(name), 0, "get product")
}
