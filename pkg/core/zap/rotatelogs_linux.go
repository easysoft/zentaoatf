package myZap

import (
	"os"
	"path"
	"time"

	"go.uber.org/zap/zapcore"

	zaprotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

// GetWriteSyncer zap logger中加入file-rotatelogs
func GetWriteSyncer() (zapcore.WriteSyncer, error) {
	fileWriter, err := zaprotatelogs.New(
		path.Join(ZapInst.Director, "%Y-%m-%d.log"),
		zaprotatelogs.WithLinkName(ZapInst.LinkName),
		zaprotatelogs.WithMaxAge(7*24*time.Hour),
		zaprotatelogs.WithRotationTime(24*time.Hour),
	)
	if ZapInst.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter)), err
	}
	return zapcore.AddSync(fileWriter), err
}
