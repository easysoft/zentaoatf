package helper

import (
	"github.com/easysoft/zentaoatf/test/ui/utils"
)

func (e *Element) Click() {
	err := e.Locator.Click()

	utils.PrintErr(err)
}
