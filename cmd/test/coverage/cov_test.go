package main

import (
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"testing"
)

var ()

func testCov() (status string) {
	return "pass"
}

// exec entrypoint
func TestCov(t *testing.T) {
	suite.RunSuite(t, new(CovSuite))
}

// suite runner
func (s *CovSuite) BeforeEach(t provider.T) {
	t.ID("1")
	t.AddSubSuite("")
}
func (s *CovSuite) TestRun(t provider.T) {
	t.Require().Equal("pass", testCov())
}

// suite def
type CovSuite struct {
	runner runner.TestRunner
}

func (s *CovSuite) GetRunner() runner.TestRunner {
	return s.runner
}

func (s *CovSuite) SetRunner(runner runner.TestRunner) {
	s.runner = runner
}
