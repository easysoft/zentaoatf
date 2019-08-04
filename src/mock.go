package main

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
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
	fmt.Println(r.Body)

	jsonString := utils.ReadFile("src/mock/json/case-from-prodoct.json")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, jsonString)
}

func submitResult(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Body)

	jsonString := utils.ReadFile("src/mock/json/success.json")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, jsonString)
}

func reportBug(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Body)

	jsonString := utils.ReadFile("src/mock/json/success.json")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, jsonString)
}
