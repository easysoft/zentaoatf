package httpClient

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	"io/ioutil"
	"net/http"
	"strings"
)

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
		return model.Response{}, nil
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
