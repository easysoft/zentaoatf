package myZap

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	ZapInst Zap
)

// getEncoderConfig 获取zapcore.EncoderConfig
func GetEncoderConfig() (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  ZapInst.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	switch {
	case ZapInst.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case ZapInst.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case ZapInst.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case ZapInst.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return config
}

// getEncoder 获取zapcore.Encoder
func getEncoder() zapcore.Encoder {
	if ZapInst.Format == "json" {
		return zapcore.NewJSONEncoder(GetEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(GetEncoderConfig())
}

// getEncoderCore 获取Encoder的zapcore.Core
func GetEncoderCore(level zapcore.Level) (core zapcore.Core) {
	writer, err := GetWriteSyncer() // 使用file-rotatelogs进行日志分割
	if err != nil {
		fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
		return
	}
	return zapcore.NewCore(getEncoder(), writer, level)
}

// 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(ZapInst.Prefix + "2006/01/02 - 15:04:05.000"))
}

type StringsArray [][]string

// MarshalLogArray 序列化数组日志
func (ss StringsArray) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range ss {
		for ii := range ss[i] {
			arr.AppendString(ss[i][ii])
		}
	}
	return nil
}

// Strings constructs a field that carries a slice of strings.
func Strings(key string, ss [][]string) zap.Field {
	return zap.Array(key, StringsArray(ss))
}
