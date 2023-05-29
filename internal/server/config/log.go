package serverConfig

import (
	"errors"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commonUtils "github.com/easysoft/zentaoatf/pkg/lib/common"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/snowlyg/helper/dir"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	WinFileSchema = "winfile:///"
)

func InitLog() {
	CONFIG.Zap.Director = filepath.Join(commConsts.WorkDir, CONFIG.Zap.Director)
	if !dir.IsExist(CONFIG.Zap.Director) { // 判断是否有Director文件夹
		dir.InsureDir(CONFIG.Zap.Director)
	}

	zapConfig := getLogConfig()

	// print to console and info、err files
	var err error
	zapConfig.EncoderConfig.EncodeLevel = nil
	zapConfig.DisableCaller = true
	zapConfig.EncoderConfig.TimeKey = ""

	logUtils.LoggerStandard, err = zapConfig.Build()
	if err != nil {
		log.Println("init console logger fail " + err.Error())
	}
	logUtils.LoggerStandard = logUtils.LoggerStandard.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))

	// print to console and info、err files without stacktrace detail
	// by set DisableStacktrace to true
	zapConfig.DisableStacktrace = true
	logUtils.LoggerExecConsole, err = zapConfig.Build(zap.WrapCore(zapCoreInFile))
	if err != nil {
		log.Println("init exec console logger fail " + err.Error())
	}

	logUtils.LoggerExecConsole = logUtils.LoggerExecConsole.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))
}

// write exec results by using zap log
func InitExecLog(workspacePath string) {
	if commConsts.ExecFrom == commConsts.FromClient {
		commConsts.ExecLogDir = logUtils.GetLogDir(workspacePath)
	} else {
		commConsts.ExecLogDir = logUtils.GetLogDir(workspacePath)
	}
	config := getLogConfig()
	config.EncoderConfig.EncodeLevel = nil

	// print to test log file
	logPathInfo := filepath.Join(commConsts.ExecLogDir, commConsts.LogText)
	if commonUtils.IsWin() {
		logPathInfo = filepath.Join(WinFileSchema, logPathInfo)
		zap.RegisterSink("winfile", newWinFileSink)
	}

	config.OutputPaths = []string{logPathInfo}
	var err error
	logUtils.LoggerExecFile, err = config.Build()
	logUtils.LoggerExecFile = logUtils.LoggerExecFile.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))

	if err != nil {
		log.Println("init exec file logger fail " + err.Error())
	}

	config.DisableCaller = true
	config.DisableStacktrace = true
	config.EncoderConfig.TimeKey = ""

	// print to test result file
	logPathResult := filepath.Join(commConsts.ExecLogDir, commConsts.ResultText)
	if commonUtils.IsWin() {
		logPathResult = filepath.Join(WinFileSchema, logPathResult)
		zap.RegisterSink("winfile", newWinFileSink)
	}
	config.OutputPaths = []string{logPathResult}
	logUtils.LoggerExecResult, err = config.Build()
	if err != nil {
		log.Println("init exec result logger fail " + err.Error())
	}

	// logUtils.LoggerExecResult = logUtils.LoggerExecResult.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))
}

// flush buffer and release the file.
func SyncExecLog() {
	if logUtils.LoggerExecFile != nil {
		logUtils.LoggerExecFile.Sync()
		logUtils.LoggerExecFile = nil
	}

	if logUtils.LoggerExecResult != nil {
		logUtils.LoggerExecResult.Sync()
		logUtils.LoggerExecResult = nil
	}
}

func getLogConfig() (config zap.Config) {
	var level zapcore.Level

	switch CONFIG.Zap.Level { // 初始化配置文件的Level
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, //这里可以指定颜色
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 路径编码器
		//EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
	}

	config = zap.Config{
		Level:         zap.NewAtomicLevelAt(level), // 日志级别
		Development:   true,                        // 开发模式，堆栈跟踪
		Encoding:      "console",                   // 输出格式 console 或 json
		EncoderConfig: encoderConfig,               // 编码器配置
		//InitialFields:    map[string]interface{}{"test_machine": "pc1"}, // 初始化字段
	}

	//if commonUtils.IsWin() {
	//	zap.RegisterSink("winfile", newWinFileSink)
	//}
	//
	//logPathInfo := filepath.Join(CONFIG.Zap.Director, "info.log")
	//logPathErr := filepath.Join(CONFIG.Zap.Director, "err.log")
	//if commonUtils.IsWin() {
	//	logPathInfo = filepath.Join(WinFileSchema, logPathInfo)
	//	logPathErr = filepath.Join(WinFileSchema, logPathErr)
	//}
	//config.OutputPaths = []string{"stdout", logPathInfo}
	//config.ErrorOutputPaths = []string{"stderr", logPathErr}

	return
}

func newWinFileSink(u *url.URL) (zap.Sink, error) {
	// Remove leading slash left by url.Parse()
	var name string
	if u.Path != "" {
		name = u.Path[1:]
	} else if u.Opaque != "" {
		name = u.Opaque[1:]
	} else {
		return nil, errors.New("path error")
	}
	return os.OpenFile(name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
}

func zapCoreInFile(c zapcore.Core) zapcore.Core {
	logPathInfo := filepath.Join(CONFIG.Zap.Director, "ztf.log")

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(&lumberjack.Logger{
				Filename:   logPathInfo,
				LocalTime:  true,
				MaxSize:    300, // M
				MaxAge:     30,  // days
				MaxBackups: 30,
			}),
			zapcore.AddSync(os.Stdout),
		),
		zap.DebugLevel,
	)
	cores := zapcore.NewTee(c, core)

	return cores
}
