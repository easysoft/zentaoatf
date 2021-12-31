package agentConfig

import myZap "github.com/aaronchen2k/deeptest/internal/pkg/core/zap"

type Config struct {
	System System    `mapstructure:"system" json:"system" yaml:"system"`
	Zap    myZap.Zap `mapstructure:"zap" json:"zap" yaml:"zap"`
}

type System struct {
	Level      string `mapstructure:"level" json:"level" yaml:"level"` // debug,release,test
	Addr       string `mapstructure:"addr" json:"addr" yaml:"addr"`
	TimeFormat string `mapstructure:"time-format" json:"timeFormat" yaml:"time-format"`
}
