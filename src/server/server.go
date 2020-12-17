package server

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zentaoatf/src/server/model"
	"github.com/easysoft/zentaoatf/src/server/service"
	serverUtils "github.com/easysoft/zentaoatf/src/server/utils"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"io"
	"io/ioutil"
	"net/http"
)

type Server struct {
	agentService *service.AgentService
}

func NewServer() *Server {
	agentService := service.NewAgentService()

	return &Server{agentService: agentService}
}

func (s *Server) Run() {
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", vari.Port),
		Handler: s.Handler(),
	}

	httpServer.ListenAndServe()
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.handle)

	return mux
}

func (s *Server) handle(writer http.ResponseWriter, req *http.Request) {
	ret := model.ResData{Code: 1, Msg: "success"}
	var err error

	serverUtils.SetupCORS(&writer, req)

	if req.Method == "GET" {
		ret, err = s.get(req)
		if err != nil {
			serverUtils.OutputErr(err, writer)
			return
		}

	} else if req.Method == "POST" {
		ret, err = s.post(req)
		if err != nil {
			serverUtils.OutputErr(err, writer)
			return
		}
	}

	bytes, _ := json.Marshal(ret)
	io.WriteString(writer, string(bytes))
}

func (s *Server) get(req *http.Request) (resp model.ResData, err error) {
	params := serverUtils.ParserGetParams(req.URL.Query())
	method := params["method"]

	switch method {

	case "ListTask":
		resp.Msg = "WORKING"

	case "":
		resp.Code = 0
		resp.Msg = "METHOD IS EMPTY"
	default:
		resp.Code = 0
		resp.Msg = "API NOT FOUND"
	}

	return
}

func (s *Server) post(req *http.Request) (resp model.ResData, err error) {
	body, err := ioutil.ReadAll(req.Body)
	if len(body) == 0 {
		return
	}

	reqData := model.ReqData{}
	err = serverUtils.ParserJsonReq(body, &reqData)
	if err != nil {
		return
	}

	method := reqData.Action

	switch method {

	case "ListTask":

	default:
		resp.Code = 0
		resp.Msg = "API NOT FOUND"
	}
	if err != nil {
		resp.Code = 0
		resp.Msg = "API ERROR: " + err.Error()
	}

	return
}
