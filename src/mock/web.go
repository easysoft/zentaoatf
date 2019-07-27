package mock

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/utils"
	"net/http"
	"net/http/httptest"
)

var Server *httptest.Server

func CreateServer(interf string) *httptest.Server {
	jsonString := utils.ReadFile("src/mock/json/" + interf)

	f := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, jsonString)
	}
	return httptest.NewServer(http.HandlerFunc(f))
}
