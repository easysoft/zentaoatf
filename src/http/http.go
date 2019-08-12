package http

import (
	"github.com/easysoft/zentaoatf/src/utils"
	"io/ioutil"
	"net/http"
	"strings"
)

func Get(url string, params map[string]string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("cookie", utils.SessionVar+"="+utils.SessionId)

	q := req.URL.Query()
	q.Add(utils.SessionVar, utils.SessionId)
	if params != nil {
		for pkey, pval := range params {
			q.Add(pkey, pval)
		}

	}
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	println(string(body))

	defer resp.Body.Close()
	return body, nil
}

func Post(url string, params map[string]string) ([]byte, error) {
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
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("cookie", utils.SessionVar+"="+utils.SessionId)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	println(string(body))

	defer resp.Body.Close()
	return body, nil
}
