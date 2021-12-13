package fileUtils

import (
	"bytes"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

func Upload(url string, files []string, extraParams map[string]string) {
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)

	for _, file := range files {
		fw, _ := bodyWriter.CreateFormFile("file", file)
		f, _ := os.Open(file)
		defer f.Close()
		io.Copy(fw, f)
	}

	for key, value := range extraParams {
		_ = bodyWriter.WriteField(key, value)
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(url, contentType, bodyBuffer)
	defer resp.Body.Close()

	if err != nil {
		logUtils.Error(i118Utils.Sprintf("fail_to_upload_file", err.Error()))
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logUtils.Error(i118Utils.Sprintf("fail_to_parse_upload_file_response", err.Error()))
	}

	logUtils.Log(i118Utils.Sprintf("upload_file_result", resp.Status, string(respBody)))
}
