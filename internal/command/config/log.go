package commandConfig

import (
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLog() {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.CallerKey = zapcore.OmitKey
	cfg.EncoderConfig.TimeKey = zapcore.OmitKey
	cfg.EncoderConfig.LevelKey = zapcore.OmitKey

	logUtils.LoggerStandard, _ = cfg.Build()
	logUtils.LoggerExecConsole, _ = cfg.Build()
}
