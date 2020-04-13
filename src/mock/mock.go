package mock

import (
	"fmt"
	"github.com/gorilla/mux"
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

func Launch() {
	r := mux.NewRouter()

	//r.HandleFunc("/"+constant.UrlZentaoSettings, zentaoSettings)

	r.Methods("POST")

	err := http.ListenAndServe("0.0.0.0:8888", r)
	if err != nil {
		log.Fatalln("ListenAndServe err:", err)
	}
}
