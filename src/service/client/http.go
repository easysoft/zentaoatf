package client

import (
	"encoding/json"
	"fmt"
	"github.com/ajg/form"
	"github.com/easysoft/zentaoatf/src/model"
	printUtils "github.com/easysoft/zentaoatf/src/utils/print"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func Get(url string, params map[string]string) (string, bool) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", false
	}

	req.Header.Set("cookie", vari.SessionVar+"="+vari.SessionId)

	q := req.URL.Query()
	q.Add(vari.SessionVar, vari.SessionId)
	if params != nil {
		for pkey, pval := range params {
			q.Add(pkey, pval)
		}

	}
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return "", false
	}

	bodyStr, _ := ioutil.ReadAll(resp.Body)
	printUtils.PrintUnicode(bodyStr)

	var bodyJson model.ZentaoResponse
	json.Unmarshal(bodyStr, &bodyJson)

	defer resp.Body.Close()

	status := bodyJson.Status
	if status == "" { // 非嵌套结构
		return string(bodyStr), true
	} else { // 嵌套结构
		dataStr := bodyJson.Data
		return dataStr, status == "success"
	}
}

func PostObject(url string, params interface{}) (string, bool) {
	client := &http.Client{}

	val, _ := form.EncodeToValues(params)
	// convert data to post fomat
	re3, _ := regexp.Compile("\\.(.*?)=")
	data := re3.ReplaceAllStringFunc(val.Encode(), replacePostData)
	printUtils.PrintToCmd(fmt.Sprintf("%s", data))

	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return "", false
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("cookie", vari.SessionVar+"="+vari.SessionId)

	resp, err := client.Do(req)
	if err != nil {
		return "", false
	}

	bodyStr, _ := ioutil.ReadAll(resp.Body)
	printUtils.PrintUnicode(bodyStr)

	var bodyJson model.ZentaoResponse
	json.Unmarshal(bodyStr, &bodyJson)

	defer resp.Body.Close()

	status := bodyJson.Status
	if status == "" { // 非嵌套结构
		return string(bodyStr), true
	} else { // 嵌套结构
		dataStr := bodyJson.Data
		return dataStr, status == "success"
	}
}

func PostStr(url string, params map[string]string) (string, bool) {
	client := &http.Client{}

	paramStr := ""
	idx := 0
	for pkey, pval := range params {
		if idx > 0 {
			paramStr += "&"
		}
		paramStr = paramStr + pkey + "=" + pval
		idx++
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(paramStr))
	if err != nil {
		return "", false
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("cookie", vari.SessionVar+"="+vari.SessionId)

	resp, err := client.Do(req)
	if err != nil {
		return "", false
	}

	bodyStr, _ := ioutil.ReadAll(resp.Body)
	printUtils.PrintUnicode(bodyStr)

	var bodyJson model.ZentaoResponse
	json.Unmarshal(bodyStr, &bodyJson)

	defer resp.Body.Close()

	status := bodyJson.Status
	if status == "" { // 非嵌套结构
		return string(bodyStr), true
	} else { // 嵌套结构
		dataStr := bodyJson.Data
		return dataStr, status == "success"
	}
}

func replacePostData(str string) string {
	str = strings.Replace(str, ".", "[", -1)
	str = strings.Replace(str, "=", "]=", -1)
	return str
}
