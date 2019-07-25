package httpClient

import (
	"io/ioutil"
	"net/http"
)

func GetBuf(url string, params map[string]string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte("")
	}

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	// fmt.Println(req.URL.String())

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return []byte(err.Error())
	}

	bytes, _ := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()
	return bytes
}

func Get(url string, params map[string]string) string {
	bts := GetBuf(url, params)
	return string(bts)
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
