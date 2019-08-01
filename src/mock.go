package main

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/mock"
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/"+mock.ImportProject, importProject)
	r.HandleFunc("/"+mock.ReportBugs, reportBugs)
	r.Methods("GET")

	err := http.ListenAndServe("0.0.0.0:8888", r)
	if err != nil {
		log.Fatalln("ListenAndServe err:", err)
	}
}

func importProject(w http.ResponseWriter, r *http.Request) {
	jsonString := utils.ReadFile("src/mock/json/case-from-prodoct.json")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, jsonString)
}

func reportBugs(w http.ResponseWriter, r *http.Request) {
	jsonString := utils.ReadFile("src/mock/json/success.json")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, jsonString)
}
