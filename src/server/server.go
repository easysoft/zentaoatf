package server

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zentaoatf/src/server/domain"
	"github.com/easysoft/zentaoatf/src/server/service"
	serverUtils "github.com/easysoft/zentaoatf/src/server/utils/common"
	serverConst "github.com/easysoft/zentaoatf/src/server/utils/const"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"io"
	"io/ioutil"
	"net/http"
)

type Server struct {
	commonService *service.CommonService
	agentService  *service.AgentService
	buildService  *service.BuildService
	taskService   *service.TaskService
	cronService   *service.CronService
}

func NewServer() *Server {
	commonService := service.NewCommonService()
	agentService := service.NewAgentService()

	taskService := service.NewTaskService()
	buildService := service.NewBuildService(taskService)

	cronService := service.NewCronService(commonService)
	cronService.Init()

	return &Server{commonService: commonService, agentService: agentService,
		buildService: buildService, taskService: taskService,
		cronService: cronService}
}

func (s *Server) Init() {
	if vari.AgentDir != "" {
		return
	}

	home, _ := serverUtils.GetUserHome()
	vari.AgentDir = fileUtils.AddPathSepIfNeeded(home + constant.PthSep + serverConst.AgentDir)
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
	ret := domain.RespData{Code: 1, Msg: "success"}
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

func (s *Server) get(req *http.Request) (resp domain.RespData, err error) {
	method, _ := serverUtils.ParserGetParams(req)

	switch method {

	case "listTask":
		resp.Msg = "listTask"
	case "listHistory":
		resp.Msg = "listHistory"

	case "":
		resp.Code = 0
		resp.Msg = "METHOD IS EMPTY"
	default:
		resp.Code = 0
		resp.Msg = "METHOD NOT FOUND"
	}

	return
}

func (s *Server) post(req *http.Request) (resp domain.RespData, err error) {
	body, err := ioutil.ReadAll(req.Body)
	if len(body) == 0 {
		return
	}

	reqData := domain.ReqData{}
	err = serverUtils.ParserJsonReq(body, &reqData)
	if err != nil {
		return
	}

	method := reqData.Action

	switch method {

	case "addTask":
		s.buildService.Add(reqData)

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
