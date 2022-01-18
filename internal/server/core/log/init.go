package serverLog

import (
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	serverConfig "github.com/aaronchen2k/deeptest/internal/server/config"
	"path/filepath"

	"github.com/snowlyg/helper/dir"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init() {
	var level zapcore.Level

	if !dir.IsExist(serverConfig.CONFIG.Zap.Director) { // 判断是否有Director文件夹
		dir.InsureDir(serverConfig.CONFIG.Zap.Director)
	}

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

	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(level), // 日志级别
		Development: true,                        // 开发模式，堆栈跟踪
		//Encoding:         "json",                                              // 输出格式 console 或 json
		Encoding:         "console",          // 输出格式 console 或 json
		EncoderConfig:    encoderConfig,      // 编码器配置
		OutputPaths:      []string{"stdout"}, // 输出到指定文件 stdout（标准输出，正常颜色） stderr（错误输出，红色）
		ErrorOutputPaths: []string{"stderr"},
		//InitialFields:    map[string]interface{}{"test_machine": "pc1"}, // 初始化字段
	}
	config.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder //这里可以指定颜色

	// print console
	logUtils.LoggerConsole, _ = config.Build()

	// print to exec log file
	config.EncoderConfig.EncodeLevel = nil
	config.OutputPaths = []string{filepath.Join(serverConfig.CONFIG.Zap.Director, "log.txt")}
	logUtils.LoggerLog, _ = config.Build()

	// print to test result file
	config.DisableCaller = true
	config.DisableStacktrace = true
	config.EncoderConfig.TimeKey = ""
	config.OutputPaths = []string{filepath.Join(serverConfig.CONFIG.Zap.Director, "result.txt")}
	logUtils.LoggerResult, _ = config.Build()

}
