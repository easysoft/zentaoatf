package server

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zentaoatf/src/server/cron"
	"github.com/easysoft/zentaoatf/src/server/domain"
	"github.com/easysoft/zentaoatf/src/server/service"
	serverUtils "github.com/easysoft/zentaoatf/src/server/utils/common"
	serverConst "github.com/easysoft/zentaoatf/src/server/utils/const"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Server struct {
	commonService *service.CommonService
	configService *service.ConfigService
	agentService  *service.AgentService
	buildService  *service.BuildService
	taskService   *service.TaskService
	cronService   *cron.CronService
}

func NewServer() *Server {
	commonService := service.NewCommonService()
	configService := service.NewConfigService()
	agentService := service.NewAgentService()
	heartBeatService := service.NewHeartBeatService()

	taskService := service.NewTaskService()
	buildService := service.NewBuildService(taskService)
	execService := service.NewExecService()
	upgradeService := service.NewUpgradeService()

	cronService := cron.NewCronService(heartBeatService, buildService, taskService, execService, upgradeService)
	cronService.Init()

	return &Server{commonService: commonService, configService: configService, agentService: agentService,
		buildService: buildService, taskService: taskService,
		cronService: cronService}
}

func (s *Server) Init() {
	vari.IP, vari.MAC = serverUtils.GetIp()

	vari.AgentLogDir = vari.ExeDir + serverConst.AgentLogDir + constant.PthSep
	err := fileUtils.MkDirIfNeeded(vari.AgentLogDir)
	if err != nil {
		logUtils.PrintTof("mkdir %s error %s", vari.AgentLogDir, err.Error())
	}
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
	resp := domain.RespData{Code: 1, Msg: "success"}
	var err error

	serverUtils.SetupCORS(&writer, req)

	if req.Method == "GET" {
		resp, err = s.get(writer, req)
		if err != nil {
			serverUtils.OutputErr(err, writer)
			return
		}

	} else if req.Method == "POST" {
		resp, err = s.post(req)
		if err != nil {
			serverUtils.OutputErr(err, writer)
			return
		}
	}

	bytes, _ := json.Marshal(resp)
	io.WriteString(writer, string(bytes))
}

func (s *Server) get(writer http.ResponseWriter, req *http.Request) (resp domain.RespData, err error) {
	resp = domain.RespData{Code: 1, Msg: "success"}
	method, params := serverUtils.ParserGetParams(req)

	switch method {

	case "listTask":
		resp.Data = s.taskService.ListTask()

	case "listHistory":
		resp.Data = s.taskService.ListHistory()

	case "download":
		Download(writer, params["f"])

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
	resp = domain.RespData{Code: 1, Msg: "success"}

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
	if method == "" {
		method, _ = serverUtils.ParserGetParams(req)
	}

	switch method {

	case "addTask":
		s.buildService.Add(reqData)

	case "config":
		s.configService.Update(reqData)

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

func Download(w http.ResponseWriter, fi string) {
	logDir := vari.ExeDir + "log-agent" + constant.PthSep
	file, _ := os.Open(logDir + strings.Replace(fi, "-", "/", 1))
	defer file.Close()

	fileHeader := make([]byte, 512)
	file.Read(fileHeader)

	fileStat, _ := file.Stat()

	w.Header().Set("Content-Disposition", "attachment; filename="+fi)
	w.Header().Set("Content-Type", http.DetectContentType(fileHeader))
	w.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))

	file.Seek(0, 0)
	io.Copy(w, file)

	return
}
