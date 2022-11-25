package commandConfig

import (
	"os"

	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLog() {
	zapConfig := zap.Config{
		Level:         zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:   true,
		Encoding:      "console",
		EncoderConfig: zapcore.EncoderConfig{},
	}

	logUtils.LoggerStandard, _ = zapConfig.Build(zap.WrapCore(zapCoreConsole))
	logUtils.LoggerExecConsole, _ = zapConfig.Build(zap.WrapCore(zapCoreConsole))
}

func zapCoreConsole(c zapcore.Core) zapcore.Core {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.TimeKey = zapcore.OmitKey
	encoderConfig.LevelKey = zapcore.OmitKey
	encoderConfig.CallerKey = zapcore.OmitKey

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
		),
		zap.DebugLevel,
	)
	cores := zapcore.NewTee(c, core)

	return cores
}
