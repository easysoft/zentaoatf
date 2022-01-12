package serverConfig

import (
	myZap "github.com/aaronchen2k/deeptest/internal/pkg/core/zap"
)

type Config struct {
	MaxSize int64     `mapstructure:"max-size" json:"burst" yaml:"max-size"`
	System  System    `mapstructure:"system" json:"system" yaml:"system"`
	Limit   Limit     `mapstructure:"limit" json:"limit" yaml:"limit"`
	Zap     myZap.Zap `mapstructure:"zap" json:"zap" yaml:"zap"`
}

type System struct {
	Version      string `mapstructure:"version" json:"version" yaml:"version"`
	Level        string `mapstructure:"level" json:"level" yaml:"level"` // debug,release,test
	Addr         string `mapstructure:"addr" json:"addr" yaml:"addr"`
	StaticPrefix string `mapstructure:"static-prefix" json:"staticPrefix" yaml:"static-prefix"`
	StaticPath   string `mapstructure:"static-path" json:"staticPath" yaml:"static-path"`
	WebPath      string `mapstructure:"web-path" json:"webPath" yaml:"web-path"`
	DbType       string `mapstructure:"db-type" json:"dbType" yaml:"db-type"`
	TimeFormat   string `mapstructure:"time-format" json:"timeFormat" yaml:"time-format"`
}

type Limit struct {
	Disable bool    `mapstructure:"disable" json:"disable" yaml:"disable"`
	Limit   float64 `mapstructure:"limit" json:"limit" yaml:"limit"`
	Burst   int     `mapstructure:"burst" json:"burst" yaml:"burst"`
}
