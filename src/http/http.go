package httpClient

import (
	"io/ioutil"
	"net/http"
)

func GetBuf(url string, params map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte(""), err
	}

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return []byte(""), err
	}

	bytes, _ := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()
	return bytes, nil
}

func Get(url string, params map[string]string) (string, error) {
	bts, err := GetBuf(url, params)
	if err != nil {
		return "", err
	}

	return string(bts), nil
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
