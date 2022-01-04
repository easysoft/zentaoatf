package web

import (
	stdContext "context"
	"fmt"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/consts"
	"github.com/aaronchen2k/deeptest/internal/server/core/cache"
	"github.com/aaronchen2k/deeptest/internal/server/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/core/viper"
	serverZap "github.com/aaronchen2k/deeptest/internal/server/core/zap"
	myWs "github.com/aaronchen2k/deeptest/internal/server/modules/v1/controller"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	gorillaWs "github.com/gorilla/websocket"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos/gorilla"
	"net/http"
	"path/filepath"
	"sync"
	"testing"
	"time"

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
func Init() *WebServer {
	serverViper.Init()
	serverZap.Init()

	err := cache.Init()
	if err != nil {
		logUtils.Errorf("init redis cache failed, error %s", err.Error())
		return nil
	}

	app := iris.New()
	app.Validator = validator.New() //参数验证
	app.Logger().SetLevel(serverConsts.CONFIG.System.Level)
	idleConnClosed := make(chan struct{})
	iris.RegisterOnInterrupt(func() { //优雅退出
		timeout := 10 * time.Second
		ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
		defer cancel()
		app.Shutdown(ctx) // close all hosts
		close(idleConnClosed)
	})

	if serverConsts.CONFIG.System.Addr == "" { // 默认 8085
		serverConsts.CONFIG.System.Addr = "127.0.0.1:8085"
	}

	if serverConsts.CONFIG.System.StaticPath == "" { // 默认 /static/upload
		serverConsts.CONFIG.System.StaticPath = "/static/upload"
	}

	if serverConsts.CONFIG.System.StaticPrefix == "" { // 默认 /upload
		serverConsts.CONFIG.System.StaticPrefix = "/upload"
	}

	if serverConsts.CONFIG.System.WebPath == "" { // 默认 /./dist
		serverConsts.CONFIG.System.WebPath = "./dist"
	}

	if serverConsts.CONFIG.System.TimeFormat == "" { // 默认 80
		serverConsts.CONFIG.System.TimeFormat = time.RFC3339
	}

	// init grpc
	mvc.New(app)

	// init websocket
	websocketAPI := app.Party(serverConsts.WsPath)
	m := mvc.New(websocketAPI)
	m.Register(
		&service.PrefixedLogger{Prefix: ""},
	)
	m.HandleWebsocket(myWs.NewWsCtrl())
	websocketServer := websocket.New(
		gorilla.Upgrader(gorillaWs.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}), m)
	websocketAPI.Get("/", websocket.Handler(websocketServer))

	return &WebServer{
		app:               app,
		addr:              serverConsts.CONFIG.System.Addr,
		timeFormat:        serverConsts.CONFIG.System.TimeFormat,
		staticPrefix:      serverConsts.CONFIG.System.StaticPrefix,
		staticPath:        serverConsts.CONFIG.System.StaticPath,
		webPath:           serverConsts.CONFIG.System.WebPath,
		idleConnsClosed:   idleConnClosed,
		globalMiddlewares: []context.Handler{},
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

	// 添加上传文件路径
	webServer.app.Listen(
		webServer.addr,
		iris.WithoutInterruptHandler,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
		iris.WithTimeFormat(webServer.timeFormat),
	)
	<-webServer.idleConnsClosed
}
