package mock

import (
	"encoding/json"
	"fmt"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	zentaoUtils "github.com/easysoft/zentaoatf/src/utils/zentao"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	BaseUrl = "client://ztpmp.ngtesting.com/"

	caseJson     = fmt.Sprintf("res%sjson%scase-from-prodoct.json", string(os.PathSeparator), string(os.PathSeparator))
	settingsJson = fmt.Sprintf("res%sjson%szentao-settings.json", string(os.PathSeparator), string(os.PathSeparator))
	successJson  = fmt.Sprintf("res%sjson%ssuccess.json", string(os.PathSeparator), string(os.PathSeparator))
)

func GetUrl(uri string) string {
	return BaseUrl + uri
}

func Launch() {
	r := mux.NewRouter()

	r.HandleFunc("/"+constant.UrlZentaoSettings, zentaoSettings)
	r.HandleFunc("/"+constant.UrlImportProject, importProject)
	r.HandleFunc("/"+constant.UrlSubmitResult, submitResult)
	r.HandleFunc("/"+constant.UrlReportBug, reportBug)

	r.Methods("POST")

	err := http.ListenAndServe("0.0.0.0:8888", r)
	if err != nil {
		log.Fatalln("ListenAndServe err:", err)
	}
}

func importProject(w http.ResponseWriter, r *http.Request) {
	printRequestBody(r.Body)

	jsonString := zentaoUtils.ReadResData(caseJson)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, jsonString)
}

func submitResult(w http.ResponseWriter, r *http.Request) {
	printRequestBody(r.Body)

	jsonString := zentaoUtils.ReadResData(successJson)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, jsonString)
}

func reportBug(w http.ResponseWriter, r *http.Request) {
	printRequestBody(r.Body)

	jsonString := zentaoUtils.ReadResData(successJson)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, jsonString)
}

func zentaoSettings(w http.ResponseWriter, r *http.Request) {
	printRequestBody(r.Body)

	jsonString := zentaoUtils.ReadResData(settingsJson)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, jsonString)
}

func printRequestBody(rd io.ReadCloser) {
	var body map[string]interface{}
	json.NewDecoder(rd).Decode(&body)
	fmt.Printf("%v\n", body)
}
