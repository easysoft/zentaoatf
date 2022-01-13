package rpc

import (
	"encoding/json"
	"github.com/aaronchen2k/deeptest/internal/comm/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/core/zap"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	httpUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/http"
	"testing"
)

func TestUpload(t *testing.T) {
	myZap.Init()
	result := commDomain.TestResult{Name: "RasaResult Title"}

	zipFile := "/Users/aaron/testResult.zip"

	result.Payload = nil
	uploadResultUrl := httpUtils.GenUrl("http://localhost:8085/", "client/build/uploadResult")

	files := []string{zipFile}
	extraParams := map[string]string{}
	json, _ := json.Marshal(result)
	extraParams["result"] = string(json)

	fileUtils.Upload(uploadResultUrl, files, extraParams)
}
