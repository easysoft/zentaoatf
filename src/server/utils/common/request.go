package serverUtils

import (
	"bytes"
	"encoding/json"
	"fmt"
	serverModel "github.com/easysoft/zentaoatf/src/server/domain"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func SetupCORS(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func OutputErr(err error, writer http.ResponseWriter) {
	errRes := ErrRes(err.Error())
	WriteRes(errRes, writer)
}

func WriteRes(ret serverModel.RespData, writer http.ResponseWriter) {
	jsonStr, _ := json.Marshal(ret)
	io.WriteString(writer, string(jsonStr))
}

func ErrRes(msg string) serverModel.RespData {
	return serverModel.RespData{Code: 0, Msg: msg}
}

func ParserJsonReq(bytes []byte, obj *serverModel.ReqData) (err error) {
	err = json.Unmarshal(bytes, &obj)
	if err != nil {
		log.Println(fmt.Sprintf("parse json error %s", err))
		return
	}

	return
}

func ParserGetParams(req *http.Request) (method string, params map[string]string) {
	path := req.URL.Path
	arr := strings.Split(path, "/")
	method = arr[1]

	values := req.URL.Query()
	params = map[string]string{}
	for key, items := range values {
		value := items[len(items)-1]
		params[key] = value
	}
	return
}

func ParserGetParam(values url.Values, name, short string) (val string) {
	for key, item := range values {
		if key == name || key == short {
			val = item[len(item)-1]
		}
	}
	return val
}

func ParserPostParam(req *http.Request, paramName1, paramName2 string, dft string, isFile bool) (ret string) {
	if paramName2 != "" && req.FormValue(paramName2) != "" {
		ret = req.FormValue(paramName2)
	} else if paramName1 != "" && req.FormValue(paramName1) != "" { // high priority than paramName2
		ret = req.FormValue(paramName1)
	}

	if isFile && ret == "" {
		postFile, _, _ := req.FormFile(paramName2)
		if postFile != nil {
			defer postFile.Close()

			buf := bytes.NewBuffer(nil)
			io.Copy(buf, postFile)
			ret = buf.String()
		}

		if ret == "" {
			postFile, _, _ = req.FormFile(paramName1)
			if postFile != nil {
				defer postFile.Close()

				buf := bytes.NewBuffer(nil)
				io.Copy(buf, postFile)
				ret = buf.String()
			}
		}
	}

	if ret == "" {
		ret = dft
	}

	return
}
