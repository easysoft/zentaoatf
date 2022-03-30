package httpUtils

import (
	"encoding/json"
	"errors"
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/ajg/form"
	"github.com/bitly/go-simplejson"
	"github.com/fatih/color"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func Get(url string) (ret []byte, err error) {
	if commConsts.Verbose {
		logUtils.Info(url)
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logUtils.Infof(color.RedString("get request failed, error: %s.", err.Error()))
		return
	}

	if strings.Index(url, "user-login") < 0 && strings.Index(url, "mode=getconfig") < 0 {
		req.Header.Add(commConsts.Token, commConsts.SessionId)
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		logUtils.Infof(color.RedString("get request failed, error: %s.", err.Error()))
		return
	}

	ret, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logUtils.Infof(color.RedString("read response failed, error ", err.Error()))
		return
	}

	if commConsts.Verbose {
		logUtils.Infof(i118Utils.Sprintf("request_response"))
		logUtils.Infof(logUtils.ConvertUnicode(ret))
	}

	jsn, err := simplejson.NewJson(ret)
	if err != nil {
		return
	}
	errMsg, _ := jsn.Get("error").String()
	if strings.ToLower(errMsg) == "unauthorized" {
		err = errors.New(commConsts.UnAuthorizedErr.Message)
		return
	}

	return
}

func Post(url string, data interface{}) (ret []byte, err error) {
	return PostOrPut(url, "POST", data)
}
func Put(url string, data interface{}) (ret []byte, err error) {
	return PostOrPut(url, "PUT", data)
}

func PostOrPut(url string, method string, data interface{}) (ret []byte, err error) {
	if commConsts.Verbose {
		logUtils.Infof(i118Utils.Sprintf("server_address") + url)
	}
	client := &http.Client{}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		logUtils.Infof(color.RedString("marshal request failed, error: %s.", err.Error()))
		return
	}

	dataStr := string(dataBytes)
	if commConsts.Verbose {
		logUtils.Infof(i118Utils.Sprintf("request_content"))
		logUtils.Infof(dataStr)
	}

	req, err := http.NewRequest(method, url, strings.NewReader(dataStr))
	if err != nil {
		logUtils.Infof(color.RedString("post request failed, error: %s.", err.Error()))
		return
	}

	if strings.Index(url, "/tokens") < 0 {
		req.Header.Add(commConsts.Token, commConsts.SessionId)
	}
	//req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		logUtils.Infof(color.RedString("post request failed, error: %s.", err.Error()))
		return
	}

	ret, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		logUtils.Infof(color.RedString("read response failed, error: %s.", err.Error()))
		return
	}

	if commConsts.Verbose {
		logUtils.Infof(i118Utils.Sprintf("request_response"))
		logUtils.Infof(logUtils.ConvertUnicode(ret))
	}

	jsn, err := simplejson.NewJson(ret)
	if err != nil {
		return
	}
	errMsg, _ := jsn.Get("error").String()
	if strings.ToLower(errMsg) == "unauthorized" {
		err = errors.New(commConsts.UnAuthorizedErr.Message)
		return
	}

	return
}

func PostWithFormat(url string, data interface{}, useFormFormat bool) (ret []byte, err error) {
	url = AddToken(url)

	if commConsts.Verbose {
		logUtils.Infof(i118Utils.Sprintf("server_address") + url)
	}
	client := &http.Client{}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		logUtils.Infof(color.RedString("marshal request failed, error: %s.", err.Error()))
		return
	}

	dataStr := string(dataBytes)
	if useFormFormat {
		dataStr, _ = form.EncodeToString(data)
		// convert data to post fomat
		re3, _ := regexp.Compile(`([^&]*?)=`)
		dataStr = re3.ReplaceAllStringFunc(dataStr, replacePostData)
	}

	if commConsts.Verbose {
		logUtils.Infof(i118Utils.Sprintf("request_content"))
		logUtils.Infof(dataStr)
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(dataStr))
	if err != nil {
		logUtils.Infof(color.RedString("post request failed, error: %s.", err.Error()))
		return
	}

	//req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		logUtils.Infof(color.RedString("post request failed, error: %s.", err.Error()))
		return
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logUtils.Infof(color.RedString("read response failed, error: %s.", err.Error()))
		return
	}

	if commConsts.Verbose {
		logUtils.Infof(i118Utils.Sprintf("request_response"))
		logUtils.Infof(logUtils.ConvertUnicode(bodyBytes))
	}

	defer resp.Body.Close()

	ret, err = GetRespErr(bodyBytes, url)

	return
}

func PostStr(url string, params map[string]string) (ret []byte, err error) {
	url = AddToken(url)

	if commConsts.Verbose {
		logUtils.Infof(i118Utils.Sprintf("server_address") + url)
	}
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

	if commConsts.Verbose {
		logUtils.Infof(i118Utils.Sprintf("request_content"))
		logUtils.Infof(paramStr)
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(paramStr))
	if err != nil {
		if commConsts.Verbose {
			logUtils.Infof(color.RedString("post string failed, error: %s.", err.Error()))
		}
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("cookie", commConsts.SessionVar+"="+commConsts.SessionId)

	resp, err := client.Do(req)
	if err != nil {
		if commConsts.Verbose {
			logUtils.Infof(color.RedString("post string failed, error: %s.", err.Error()))
		}
		return
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logUtils.Infof(color.RedString("read response failed, error ", err.Error()))
		return
	}

	if commConsts.Verbose {
		logUtils.Infof(i118Utils.Sprintf("request_response"))
		logUtils.Infof(logUtils.ConvertUnicode(bodyBytes))
	}

	return
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

func replacePostData(str string) string {
	str = strings.ToLower(str[:1]) + str[1:]

	if strings.Index(str, ".") > -1 {
		str = strings.Replace(str, ".", "[", -1)
		str = strings.Replace(str, "=", "]=", -1)
	}
	return str
}

func AddToken(url string) (ret string) {
	address := url
	hash := ""

	index := strings.Index(url, "#")
	if index > 0 {
		address = url[:index]
		hash = url[index:]
	}

	paramPir := commConsts.SessionVar + "=" + commConsts.SessionId

	if commConsts.RequestType == commConsts.PathInfo && strings.Index(address, "?") < 0 {
		address = address + "?" + paramPir
	} else {
		address = address + "&" + paramPir
	}

	ret = address + "&XDEBUG_SESSION_START=PHPSTORM" + hash

	return
}

func GetRespErr(bytes []byte, url string) (ret []byte, err error) {
	ret = bytes

	if len(bytes) == 0 {
		return
	}

	var zentaoResp serverDomain.ZentaoResp
	err = json.Unmarshal(bytes, &zentaoResp)
	if err != nil {
		err = errors.New("Wrong Zentao response, unmarshal to serverDomain.ZentaoResp failed, error " + err.Error())
		if commConsts.Verbose {
			if strings.Index(url, "login") < 0 { // jsonErr caused by login request return a html
				logUtils.Infof(color.RedString(err.Error()))
			}
		}
		return
	}

	// 嵌套结构，map[status:success, data:{}]
	status := zentaoResp.Status
	if status != "" {
		ret = []byte(zentaoResp.Data)
		if status == "success" {
			err = errors.New(zentaoResp.Data)
		}
		return
	}

	// 非嵌套结构，map[result:success]
	var respData = serverDomain.ZentaoRespData{}
	err = json.Unmarshal(bytes, &respData)

	if err == nil && (respData.Result != "" && respData.Result != "success") {
		err = errors.New(string(bytes))
	}

	return
}
