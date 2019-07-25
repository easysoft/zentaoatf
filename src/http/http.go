package httpClient

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	"io/ioutil"
	"net/http"
)

func Get(url string, params map[string]string) (model.Response, error) {
	var ret model.Response

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ret, err
	}

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return ret, err
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ret, err
	}

	var respModel model.Response
	json.Unmarshal(bytes, &respModel)
	if respModel.Code != 1 {
		return ret, errors.New(fmt.Sprintf("request fail, code %d", respModel.Code))
	}

	defer resp.Body.Close()
	return respModel, nil
}
