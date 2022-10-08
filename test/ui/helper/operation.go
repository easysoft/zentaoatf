package plw

import (
	"errors"
	"fmt"
	"github.com/easysoft/zentaoatf/test/ui/utils"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (l MyLocator) Click(t provider.T) (err error) {
	err = l.Locator.Click()

	err = errors.New(fmt.Sprintf("%s: %s", l.Selector, err.Error()))

	utils.PrintErrOrNot(err, t)

	return
}
