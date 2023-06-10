package commonTestHelper

import (
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func ReplaceLabel(t provider.T, value string) {
	label := allure.Label{Name: "suite", Value: value}
	t.ReplaceLabel(&label)
}
