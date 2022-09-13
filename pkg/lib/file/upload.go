package fileUtils

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"

	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
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
		logUtils.Error(i118Utils.Sprintf("fail_to_upload_files", err.Error()))
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logUtils.Error(i118Utils.Sprintf("fail_to_parse_upload_response", err.Error()))
	}

	logUtils.Info(i118Utils.Sprintf("upload_status", resp.Status, string(respBody)))
}

func UploadWithResp(url string, files []string, extraParams map[string]string) map[string]interface{} {
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
		logUtils.Error(i118Utils.Sprintf("fail_to_upload_files", err.Error()))
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logUtils.Error(i118Utils.Sprintf("fail_to_parse_upload_response", err.Error()))
	}

	logUtils.Info(i118Utils.Sprintf("upload_status", resp.Status, string(respBody)))
	respMap := make(map[string]interface{})
	_ = json.Unmarshal(respBody, &respMap)
	return respMap
}
