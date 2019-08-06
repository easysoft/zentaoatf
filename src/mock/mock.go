package mock

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

var (
	BaseUrl = "http://localhost:8888/"
)

func GetUrl(uri string) string {
	return BaseUrl + uri
}

func Launch() {
	r := mux.NewRouter()

	r.HandleFunc("/"+utils.UrlZentaoSettings, zentaoSettings)
	r.HandleFunc("/"+utils.UrlImportProject, importProject)
	r.HandleFunc("/"+utils.UrlSubmitResult, submitResult)
	r.HandleFunc("/"+utils.UrlReportBug, reportBug)

	r.Methods("POST")

	err := http.ListenAndServe("0.0.0.0:8888", r)
	if err != nil {
		log.Fatalln("ListenAndServe err:", err)
	}
}

func importProject(w http.ResponseWriter, r *http.Request) {
	printRequestBody(r.Body)

	jsonString := utils.ReadFile("src/mock/json/case-from-prodoct.json")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, jsonString)
}

func submitResult(w http.ResponseWriter, r *http.Request) {
	printRequestBody(r.Body)

	jsonString := utils.ReadFile("src/mock/json/success.json")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, jsonString)
}

func reportBug(w http.ResponseWriter, r *http.Request) {
	printRequestBody(r.Body)

	jsonString := utils.ReadFile("src/mock/json/success.json")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, jsonString)
}

func zentaoSettings(w http.ResponseWriter, r *http.Request) {
	printRequestBody(r.Body)

	jsonString := utils.ReadFile("src/mock/json/zentao-settings.json")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, jsonString)
}

func printRequestBody(rd io.ReadCloser) {
	var body map[string]interface{}
	json.NewDecoder(rd).Decode(&body)
	fmt.Printf("%v\n", body)
}
