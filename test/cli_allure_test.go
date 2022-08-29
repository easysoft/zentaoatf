package main

import (
	"fmt"
	expect "github.com/easysoft/zentaoatf/pkg/lib/expect"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type AsyncSuite struct {
	suite.Suite
}

func (s *AsyncSuite) BeforeEach(t provider.T) {
	t.ID("1")
	t.AddSubSuite("子套件1")
}

func (s *AsyncSuite) TestAsyncStep(t provider.T) {
	t.Title("Async Testing")

	t.Parallel()

	t.WithNewAsyncStep("Async Step 1", func(ctx provider.StepCtx) {
		t.Require().Equal("pass", testList("ztf list ./demo", regexp.MustCompile("Found 5 test cases|发现5个用例")), "Assertion Success")
	})

	t.WithNewAsyncStep("Async Step 2", func(ctx provider.StepCtx) {
		t.Require().Equal("pass", testList("ztf ls ./demo -k 1", regexp.MustCompile("Found 2 test cases|发现2个用例")), "Assertion Success")
	})

	t.WithNewAsyncStep("Async Step 3", func(ctx provider.StepCtx) {
		t.Require().Equal("pass", testList("ztf ls demo -k match", regexp.MustCompile("Found 3 test cases|发现3个用例")), "Assertion Success")
	})
}

func testList(cmd string, successRe *regexp.Regexp) string {
	if runtime.GOOS == "windows" {
		cmd = strings.ReplaceAll(cmd, "/", "\\")
	}
	child, err := expect.Spawn(cmd, -1)
	if err != nil {
		return err.Error()
	}
	defer child.Close()

	if _, err = child.Expect(successRe, 10*time.Second); err != nil {
		return fmt.Sprintf("expect %s, actual %s", successRe, err.Error())
	}

	return "pass"
}

func TestSuite(t *testing.T) {
	suite.RunSuite(t, new(AsyncSuite))
}
