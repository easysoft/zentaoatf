package serverConfig

import (
	"fmt"
	myZap "github.com/aaronchen2k/deeptest/internal/pkg/core/zap"
)

type Config struct {
	MaxSize int64     `mapstructure:"max-size" json:"burst" yaml:"max-size"`
	System  System    `mapstructure:"system" json:"system" yaml:"system"`
	Limit   Limit     `mapstructure:"limit" json:"limit" yaml:"limit"`
	Zap     myZap.Zap `mapstructure:"zap" json:"zap" yaml:"zap"`
	Redis   Redis     `mapstructure:"redis" json:"redis" yaml:"redis"`
	Mysql   Mysql     `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Captcha Captcha   `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
}

type System struct {
	Level        string `mapstructure:"level" json:"level" yaml:"level"` // debug,release,test
	Addr         string `mapstructure:"addr" json:"addr" yaml:"addr"`
	StaticPrefix string `mapstructure:"static-prefix" json:"staticPrefix" yaml:"static-prefix"`
	StaticPath   string `mapstructure:"static-path" json:"staticPath" yaml:"static-path"`
	WebPath      string `mapstructure:"web-path" json:"webPath" yaml:"web-path"`
	DbType       string `mapstructure:"db-type" json:"dbType" yaml:"db-type"`
	CacheType    string `mapstructure:"cache-type" json:"cacheType" yaml:"cache-type"`
	TimeFormat   string `mapstructure:"time-format" json:"timeFormat" yaml:"time-format"`
}

type Limit struct {
	Disable bool    `mapstructure:"disable" json:"disable" yaml:"disable"`
	Limit   float64 `mapstructure:"limit" json:"limit" yaml:"limit"`
	Burst   int     `mapstructure:"burst" json:"burst" yaml:"burst"`
}

type Redis struct {
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`
	Addr     string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	PoolSize int    `mapstructure:"pool-size" json:"poolSize" yaml:"pool-size"`
}

type Mysql struct {
	Path         string `mapstructure:"path" json:"path" yaml:"path"`
	Config       string `mapstructure:"config" json:"config" yaml:"config"`
	Dbname       string `mapstructure:"db-name" json:"dbname" yaml:"db-name"`
	Username     string `mapstructure:"username" json:"username" yaml:"username"`
	Password     string `mapstructure:"password" json:"password" yaml:"password"`
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"maxIdleConns" yaml:"max-idle-conns"`
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"maxOpenConns" yaml:"max-open-conns"`
	LogMode      bool   `mapstructure:"log-mode" json:"logMode" yaml:"log-mode"`
	LogZap       string `mapstructure:"log-zap" json:"logZap" yaml:"log-zap"` //silent,error,warn,info,zap
}

func (m *Mysql) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", m.Username, m.Password, m.Path, m.Dbname, m.Config)
}

type Captcha struct {
	KeyLong   int `mapstructure:"key-long" json:"keyLong" yaml:"key-long"`
	ImgWidth  int `mapstructure:"img-width" json:"imgWidth" yaml:"img-width"`
	ImgHeight int `mapstructure:"img-height" json:"imgHeight" yaml:"img-height"`
}
