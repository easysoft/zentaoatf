package httpUtils

import (
	"encoding/json"
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"io/ioutil"
	"net/http"
	"strings"
)

func Get(url string) (interface{}, bool) {
	return GetObj(url, "farm")
}

func GetObj(url string, requestTo string) (interface{}, bool) {
	client := &http.Client{}

	if consts.Verbose {
		logUtils.Info(url)
	}

	req, reqErr := http.NewRequest("GET", url, nil)
	if reqErr != nil {
		logUtils.Error(reqErr.Error())
		return nil, false
	}

	resp, respErr := client.Do(req)

	if respErr != nil {
		logUtils.Error(respErr.Error())
		return nil, false
	}

	bodyStr, _ := ioutil.ReadAll(resp.Body)
	if consts.Verbose {
		logUtils.PrintUnicode(bodyStr)
	}
	defer resp.Body.Close()

	if requestTo == "farm" {
		var bodyJson domain.RpcResp
		jsonErr := json.Unmarshal(bodyStr, &bodyJson)
		if jsonErr != nil {
			if strings.Index(string(bodyStr), "<html>") > -1 {
				logUtils.Error(i118Utils.Sprintf("wrong_response_format", "html"))
				return nil, false
			} else {
				logUtils.Error(jsonErr.Error())
				return nil, false
			}
		}
		code := bodyJson.Code
		return bodyJson.Payload, code == consts.ResultCodeSuccess
	} else {
		var bodyJson map[string]interface{}
		jsonErr := json.Unmarshal(bodyStr, &bodyJson)
		if jsonErr != nil {
			logUtils.Error(jsonErr.Error())
			return nil, false
		} else {
			return bodyJson, true
		}
	}
}

func Post(url string, params interface{}) (interface{}, bool) {
	if consts.Verbose {
		logUtils.Info(url)
	}
	client := &http.Client{}

	paramStr, err := json.Marshal(params)
	if err != nil {
		logUtils.Error(err.Error())
		return nil, false
	}

	req, reqErr := http.NewRequest("POST", url, strings.NewReader(string(paramStr)))
	if reqErr != nil {
		logUtils.Error(reqErr.Error())
		return nil, false
	}

	req.Header.Set("Content-Type", "application/json")

	resp, respErr := client.Do(req)
	if respErr != nil {
		logUtils.Error(respErr.Error())
		return nil, false
	}

	bodyStr, _ := ioutil.ReadAll(resp.Body)
	if consts.Verbose {
		logUtils.PrintUnicode(bodyStr)
	}

	var result domain.RpcResp
	json.Unmarshal(bodyStr, &result)

	defer resp.Body.Close()

	code := result.Code
	return result, code == consts.ResultCodeSuccess
}

func GenUrl(server string, path string) string {
	server = UpdateUrl(server)
	url := fmt.Sprintf("%sapi/v1/%s", server, path)
	return url
}

func UpdateUrl(url string) string {
	if strings.LastIndex(url, "/") < len(url)-1 {
		url += "/"
	}
	return url
}
