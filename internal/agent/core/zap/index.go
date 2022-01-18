package agentZap

import (
	"github.com/aaronchen2k/deeptest/internal/agent/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/core/zap"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"

	"github.com/snowlyg/helper/dir"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// level 日志级别
var level zapcore.Level

// Init 初始化日志服务
func Init() {
	var logger *zap.Logger

	if !dir.IsExist(agentConsts.CONFIG.Zap.Director) { // 判断是否有Director文件夹
		dir.InsureDir(agentConsts.CONFIG.Zap.Director)
	}

	switch agentConsts.CONFIG.Zap.Level { // 初始化配置文件的Level
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

	if level == zap.DebugLevel || level == zap.ErrorLevel {
		logger = zap.New(myZap.GetEncoderCore(level), zap.AddStacktrace(level))
	} else {
		logger = zap.New(myZap.GetEncoderCore(level))
	}
	if agentConsts.CONFIG.Zap.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}

	logUtils.LoggerConsole = logger
}
