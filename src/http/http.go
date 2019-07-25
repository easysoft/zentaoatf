package httpClient

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func Get(url string, params map[string]string) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.String())

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return resp.Status
	}

	bytes, _ := ioutil.ReadAll(resp.Body)

	return string(bytes)
}

func GetMock(url string, params map[string]string) string {
	resp, err := http.Get(url)

	if err != nil {
		return err.Error()
	}

	bytes, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return string(bytes)
}
