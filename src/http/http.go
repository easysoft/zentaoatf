package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/easysoft/zentaoatf/src/model"
	"io/ioutil"
	"net/http"
	"strings"
)

func Get(url string, params map[string]string) (bool, *simplejson.Json, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, nil, err
	}

	if params != nil {
		q := req.URL.Query()
		for pkey, pval := range params {
			q.Add(pkey, pval)
		}
		req.URL.RawQuery = q.Encode()
	}

	//req.Header.Set("Cookie", "name=anny")

	resp, err := client.Do(req)
	if err != nil {
		return false, nil, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	json, err := simplejson.NewJson([]byte(body))
	status, err := json.Get("status").String()

	pass := status == "success"

	if err != nil && pass {
		return false, nil, err
	}

	dataStr, _ := json.Get("data").String()
	data, _ := simplejson.NewJson([]byte(dataStr))

	//if respModel.Code != 1 {
	//	return ret, errors.New(fmt.Sprintf("request fail, code %d", respModel.Code))
	//}

	defer resp.Body.Close()
	return pass, data, nil
}

func Post(url string, jsonStr string) (model.Response, error) {
	client := &http.Client{}
	var ret model.Response

	req, err := http.NewRequest("POST", url, strings.NewReader(jsonStr))
	if err != nil {
		return ret, err
	}

	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Cookie", "name=anny")

	resp, err := client.Do(req)
	if err != nil {
		return model.Response{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	var respModel model.Response
	json.Unmarshal(body, &respModel)
	if respModel.Code != 1 {
		return ret, errors.New(fmt.Sprintf("request fail, code %d", respModel.Code))
	}

	defer resp.Body.Close()
	return respModel, nil
}
