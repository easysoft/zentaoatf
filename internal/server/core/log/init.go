package serverLog

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	serverConfig "github.com/aaronchen2k/deeptest/internal/server/config"
	"log"
	"path/filepath"

	"github.com/snowlyg/helper/dir"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init() {
	if !dir.IsExist(serverConfig.CONFIG.Zap.Director) { // 判断是否有Director文件夹
		dir.InsureDir(serverConfig.CONFIG.Zap.Director)
	}

	config := getLogConfig()

	// print to console
	var err error
	config.EncoderConfig.EncodeLevel = nil
	config.DisableCaller = true
	config.EncoderConfig.TimeKey = ""
	logUtils.LoggerStandard, err = config.Build()
	if err != nil {
		log.Println("init console logger fail " + err.Error())
	}

	// print to console without detail
	config.DisableStacktrace = true
	logUtils.LoggerExecConsole, err = config.Build()
	if err != nil {
		log.Println("init exec console logger fail " + err.Error())
	}
}

func InitExecLog(projectPath string) {
	config := getLogConfig()

	commConsts.ExecLogDir = logUtils.GetLogDir(projectPath)

	// print to exec log file
	config.EncoderConfig.EncodeLevel = nil
	config.OutputPaths = []string{filepath.Join(commConsts.ExecLogDir, "log.txt")}
	var err error
	logUtils.LoggerExecFile, err = config.Build()
	if err != nil {
		log.Println("init exec file logger fail " + err.Error())
	}

	config.DisableCaller = true
	config.DisableStacktrace = true
	config.EncoderConfig.TimeKey = ""

	// print to test result file
	config.OutputPaths = []string{filepath.Join(commConsts.ExecLogDir, "result.txt")}
	logUtils.LoggerExecResult, err = config.Build()
	if err != nil {
		log.Println("init exec result logger fail " + err.Error())
	}
}

func getLogConfig() (config zap.Config) {
	var level zapcore.Level

	switch serverConfig.CONFIG.Zap.Level { // 初始化配置文件的Level
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

	config = zap.Config{
		Level:       zap.NewAtomicLevelAt(level), // 日志级别
		Development: true,                        // 开发模式，堆栈跟踪
		//Encoding:         "json",               // 输出格式 console 或 json
		Encoding:         "console",          // 输出格式 console 或 json
		EncoderConfig:    encoderConfig,      // 编码器配置
		OutputPaths:      []string{"stdout"}, // 输出到指定文件 stdout（标准输出，正常颜色） stderr（错误输出，红色）
		ErrorOutputPaths: []string{"stderr"},
		//InitialFields:    map[string]interface{}{"test_machine": "pc1"}, // 初始化字段
	}
	config.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder //这里可以指定颜色

	return
}
