package execHelper

import (
	"encoding/json"
	"fmt"
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	stringUtils "github.com/easysoft/zentaoatf/pkg/lib/string"
	"log"
	"regexp"
	"strings"
)

func ConvertK6Result(results []interface{}, failedCaseIdToThresholdMap map[string]string) commDomain.UnitTestSuite {
	caseResultMap := map[string]commDomain.UnitResult{}

	// parse checks
	for _, result := range results {
		point, ok := result.(commDomain.K6Point)
		if !ok || point.Type != commConsts.Point || point.Metric != "checks" {
			continue
		}

		caseId := point.Data.Tags.Id
		caseName := strings.TrimLeft(point.Data.Tags.Name, ":")
		caseResult, ok := caseResultMap[caseName]
		if !ok { // create if not exist
			if point.Data.Tags.Group == "" { // not a case
				continue
			}

			suite := strings.Trim(strings.TrimSuffix(point.Data.Tags.Group, "CASE"), ":")
			suite = regexp.MustCompile(":+").ReplaceAllString(suite, "-")

			caseResultMap[caseName] = commDomain.UnitResult{
				Cid:       stringUtils.ParseInt(caseId),
				Title:     point.Data.Tags.Name,
				TestSuite: suite,
				Status:    commConsts.PASS,
			}
			caseResult = caseResultMap[caseName]
		}

		if caseResult.StartTime == 0 || caseResult.StartTime > point.Data.Time.Unix() {
			caseResult.StartTime = point.Data.Time.Unix()
		}
		if caseResult.EndTime == 0 || caseResult.EndTime < point.Data.Time.Unix() {
			caseResult.EndTime = point.Data.Time.Unix()
		}

		if point.Metric == "checks" && point.Data.Value == 0 {
			caseResult.Status = commConsts.FAIL
			caseResult.Failure = &commDomain.Failure{Type: "CheckFailed", Desc: point.Data.Tags.Checkpoint}

		} else if failedCaseIdToThresholdMap[caseId] != "" {
			caseResult.Status = commConsts.FAIL
			caseResult.Failure = &commDomain.Failure{Type: "ThresholdFailed", Desc: failedCaseIdToThresholdMap[caseId]}

		}

		caseResultMap[caseName] = caseResult
	}

	var startTime, endTime int64
	testSuite := commDomain.UnitTestSuite{}
	for _, cs := range caseResultMap {
		cs.Duration = float32(cs.EndTime - cs.StartTime)

		testSuite.Cases = append(testSuite.Cases, cs)

		if startTime == 0 || startTime > cs.StartTime {
			startTime = cs.StartTime
		}
		if endTime == 0 || endTime < cs.EndTime {
			endTime = cs.EndTime
		}
	}

	testSuite.Time = float32(startTime)
	testSuite.Duration = endTime - startTime

	return testSuite
}

func GetK6FailCaseInSummary(content string) (ret map[string]string) {
	ret = map[string]string{}

	k6Summary := commDomain.K6Summary{}
	errInner := json.Unmarshal([]byte(content), &k6Summary)
	if errInner != nil {
		return
	}

	for key, val := range k6Summary.Metrics {
		regx := regexp.MustCompile(`.+\{id:(\d+)\}`)
		arr := regx.FindAllStringSubmatch(key, -1)
		if len(arr) > 0 { // http_req_duration{id:1}
			id := arr[0][1]
			log.Print(id)

			thresholdsMap := val.(map[string]interface{})["thresholds"].(map[string]interface{})
			for key2, val2 := range thresholdsMap {
				ok := val2.(map[string]interface{})["ok"]
				if stringUtils.ItoStr(ok) == "false" {
					ret[id] = fmt.Sprintf("'%s': ['%s']", key, key2)
				}
			}
		}
	}

	return
}
