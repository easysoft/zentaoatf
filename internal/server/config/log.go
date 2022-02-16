package serverConfig

import (
	"errors"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commonUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/snowlyg/helper/dir"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLog() {
	CONFIG.Zap.Director = filepath.Join(commConsts.WorkDir, CONFIG.Zap.Director)
	if !dir.IsExist(CONFIG.Zap.Director) { // 判断是否有Director文件夹
		dir.InsureDir(CONFIG.Zap.Director)
	}

	zapConfig := getLogConfig()

	// print to console
	var err error
	zapConfig.EncoderConfig.EncodeLevel = nil
	zapConfig.DisableCaller = true
	zapConfig.EncoderConfig.TimeKey = ""
	logUtils.LoggerStandard, err = zapConfig.Build()
	if err != nil {
		log.Println("init console logger fail " + err.Error())
	}

	// print to console without detail
	zapConfig.DisableStacktrace = true
	logUtils.LoggerExecConsole, err = zapConfig.Build()
	if err != nil {
		log.Println("init exec console logger fail " + err.Error())
	}
}

// 执行日志，用于具体的测试执行
func InitExecLog(projectPath string) {
	commConsts.ExecLogDir = logUtils.GetLogDir(projectPath)
	config := getLogConfig()
	config.EncoderConfig.EncodeLevel = nil

	// print to test log file
	logPath := filepath.Join(commConsts.ExecLogDir, commConsts.LogText)
	if commonUtils.IsWin() {
		logPath = filepath.Join("winfile:///", logPath)
		zap.RegisterSink("winfile", newWinFileSink)
	}

	config.OutputPaths = []string{logPath}
	var err error
	logUtils.LoggerExecFile, err = config.Build()
	if err != nil {
		log.Println("init exec file logger fail " + err.Error())
	}

	config.DisableCaller = true
	config.DisableStacktrace = true
	config.EncoderConfig.TimeKey = ""

	// print to test result file
	logPathResult := filepath.Join(commConsts.ExecLogDir, commConsts.ResultText)
	if commonUtils.IsWin() {
		logPathResult = filepath.Join("winfile:///", logPathResult)
		zap.RegisterSink("winfile", newWinFileSink)
	}
	config.OutputPaths = []string{logPathResult}
	logUtils.LoggerExecResult, err = config.Build()
	if err != nil {
		log.Println("init exec result logger fail " + err.Error())
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
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		//EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, //这里可以指定颜色
		EncodeTime:     zapcore.ISO8601TimeEncoder,       // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
	}

	logPathInfo := filepath.Join(CONFIG.Zap.Director, "info.log")
	logPathErr := filepath.Join(CONFIG.Zap.Director, "err.log")
	if commonUtils.IsWin() {
		logPathInfo = filepath.Join("winfile:///", logPathInfo)
		logPathErr = filepath.Join("winfile:///", logPathErr)
	}
	log.Printf("%s %s", logPathInfo, logPathErr)

	config = zap.Config{
		Level:         zap.NewAtomicLevelAt(level), // 日志级别
		Development:   true,                        // 开发模式，堆栈跟踪
		Encoding:      "console",                   // 输出格式 console 或 json
		EncoderConfig: encoderConfig,               // 编码器配置
		//InitialFields:    map[string]interface{}{"test_machine": "pc1"}, // 初始化字段
	}
	config.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder //这里可以指定颜色
	if commonUtils.IsWin() {
		zap.RegisterSink("winfile", newWinFileSink)
	}
	config.OutputPaths = []string{"stdout", logPathInfo}
	config.ErrorOutputPaths = []string{"stderr", logPathErr}

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
	return os.OpenFile(name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
}
