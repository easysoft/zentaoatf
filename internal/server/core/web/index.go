package web

import (
	stdContext "context"
	"fmt"
	"net/http"
	"path/filepath"
	"sync"
	"testing"
	"time"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	langHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/lang"
	websocketHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/websocket"
	serverConfig "github.com/easysoft/zentaoatf/internal/server/config"
	"github.com/easysoft/zentaoatf/internal/server/core/cron"
	"github.com/easysoft/zentaoatf/internal/server/core/dao"
	"github.com/easysoft/zentaoatf/internal/server/core/module"
	v1 "github.com/easysoft/zentaoatf/internal/server/modules/v1"
	myWs "github.com/easysoft/zentaoatf/internal/server/modules/v1/controller"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	"github.com/facebookgo/inject"
	gorillaWs "github.com/gorilla/websocket"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos/gorilla"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/helper/tests"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

var client *tests.Client

// WebServer web 服务
// - app iris application
// - modules 服务的模块
// - idleConnsClosed
// - addr  服务访问地址
// - timeFormat  时间格式
// - globalMiddlewares  全局中间件
// - wg  sync.WaitGroup
// - staticPrefix  静态文件访问地址前缀
// - staticPath  静态文件地址
// - webPath  前端文件地址
type WebServer struct {
	app               *iris.Application
	modules           []module.WebModule
	idleConnsClosed   chan struct{}
	addr              string
	timeFormat        string
	globalMiddlewares []context.Handler
	wg                sync.WaitGroup
	staticPrefix      string
	staticPath        string
	webPath           string
}

// Init 初始化web服务
func Init(port int) *WebServer {
	serverConfig.Init()
	serverConfig.InitLog()
	i118Utils.Init(commConsts.Language, commConsts.AppServer)

	langHelper.GetExtToNameMap()
	langHelper.GetEditorExtToLangMap()

	app := iris.New()
	app.Validator = validator.New() //参数验证
	app.Logger().SetLevel(serverConfig.CONFIG.System.Level)
	idleConnClosed := make(chan struct{})
	iris.RegisterOnInterrupt(func() { //优雅退出
		timeout := 10 * time.Second
		ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
		defer cancel()
		app.Shutdown(ctx) // close all hosts
		close(idleConnClosed)
	})

	mvc.New(app)

	// init websocket
	websocketCtrl := myWs.NewWebSocketCtrl()
	injectWsModule(websocketCtrl)

	websocketAPI := app.Party(serverConfig.WsPath)
	m := mvc.New(websocketAPI)
	m.Register(
		&websocketHelper.PrefixedLogger{Prefix: ""},
	)
	m.HandleWebsocket(websocketCtrl)
	websocketServer := websocket.New(gorilla.Upgrader(gorillaWs.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}), m)
	websocketAPI.Get("/", websocket.Handler(websocketServer))

	// init http
	if port != 0 {
		serverConfig.CONFIG.System.Addr = fmt.Sprintf(":%d", port)
	}
	webServer := &WebServer{
		app:               app,
		addr:              serverConfig.CONFIG.System.Addr,
		timeFormat:        serverConfig.CONFIG.System.TimeFormat,
		staticPrefix:      serverConfig.CONFIG.System.StaticPrefix,
		staticPath:        serverConfig.CONFIG.System.StaticPath,
		webPath:           serverConfig.CONFIG.System.WebPath,
		idleConnsClosed:   idleConnClosed,
		globalMiddlewares: []context.Handler{},
	}

	injectModule(webServer)

	return webServer
}

func injectModule(ws *WebServer) {
	var g inject.Graph
	g.Logger = logUtils.LoggerStandard.Sugar()

	cron := cron.NewServerCron()
	cron.Init()
	indexModule := v1.NewIndexModule()

	// inject objects
	if err := g.Provide(
		&inject.Object{Value: dao.GetDB()},
		&inject.Object{Value: cron},
		&inject.Object{Value: indexModule},
	); err != nil {
		logUtils.Fatalf("provide usecase objects to the Graph: %v", err)
	}
	err := g.Populate()
	if err != nil {
		logUtils.Fatalf("populate the incomplete Objects: %v", err)
	}

	ws.AddModule(indexModule.Party())

	logUtils.Infof("start server")
}

func injectWsModule(websocketCtrl *myWs.WebSocketCtrl) {
	var g inject.Graph
	g.Logger = logUtils.LoggerStandard.Sugar()

	// inject objects
	if err := g.Provide(
		&inject.Object{Value: dao.GetDB()},
		&inject.Object{Value: websocketCtrl},
	); err != nil {
		logUtils.Fatalf("provide usecase objects to the Graph: %v", err)
	}
	err := g.Populate()
	if err != nil {
		logUtils.Fatalf("populate the incomplete Objects: %v", err)
	}
}

// GetStaticPath 获取静态路径
func (webServer *WebServer) GetStaticPath() string {
	return webServer.staticPath
}

// GetWebPath 获取前端路径
func (webServer *WebServer) GetWebPath() string {
	return webServer.webPath
}

// GetAddr 获取web服务地址
func (webServer *WebServer) GetAddr() string {
	return webServer.addr
}

// AddModule 添加模块
func (webServer *WebServer) AddModule(module ...module.WebModule) {
	webServer.modules = append(webServer.modules, module...)
}

// AddStatic 添加静态文件
func (webServer *WebServer) AddStatic(requestPath string, fsOrDir interface{}, opts ...iris.DirOptions) {
	webServer.app.HandleDir(requestPath, fsOrDir, opts...)
}

// AddWebStatic 添加前端访问地址
func (webServer *WebServer) AddWebStatic(requestPath string) {
	fsOrDir := iris.Dir(filepath.Join(dir.GetCurrentAbPath(), webServer.webPath))
	webServer.AddStatic(requestPath, fsOrDir, iris.DirOptions{
		IndexName: "index.html",
		SPA:       true,
	})
}

// AddUploadStatic 添加上传文件访问地址
func (webServer *WebServer) AddUploadStatic() {
	fsOrDir := iris.Dir(filepath.Join(dir.GetCurrentAbPath(), webServer.staticPath))
	webServer.AddStatic(webServer.staticPrefix, fsOrDir)
}

// GetModules 获取模块
func (webServer *WebServer) GetModules() []module.WebModule {
	return webServer.modules
}

// GetTestAuth 获取测试验证客户端
func (webServer *WebServer) GetTestAuth(t *testing.T) *tests.Client {
	var once sync.Once
	once.Do(
		func() {
			client = tests.New(str.Join("http://", webServer.addr), t, webServer.app)
			if client == nil {
				t.Fatalf("client is nil")
			}
		},
	)

	return client
}

// GetTestLogin 测试登录web服务
func (webServer *WebServer) GetTestLogin(t *testing.T, url string, res tests.Responses, datas ...map[string]interface{}) *tests.Client {
	client := webServer.GetTestAuth(t)
	err := client.Login(url, res, datas...)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// Run 启动web服务
func (webServer *WebServer) Run() {
	webServer.app.UseGlobal(webServer.globalMiddlewares...)
	err := webServer.InitRouter()
	if err != nil {
		fmt.Printf("初始化路由错误： %v\n", err)
		panic(err)
	}

	//logUtils.Info(i118Utils.Sprintf("start_server", "localhost",
	//	strings.Replace(webServer.addr, ":", "", -1)))

	webServer.app.Listen(
		webServer.addr,
		iris.WithoutInterruptHandler,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
		iris.WithTimeFormat(webServer.timeFormat),
	)
	<-webServer.idleConnsClosed
}
