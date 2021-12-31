package serverDomain

type DataRequest struct {
	Sys       DataSys   `json:"sys"`
	Db        DataDb    `json:"db"`
	SqlType   string    `json:"sqlType" validate:"required"`
	Cache     DataCache `json:"cache"`
	CacheType string    `json:"cacheType"  validate:"required"`
	Level     string    `json:"level"` // debug,release,test
	Addr      string    `json:"addr"`
	ClearData bool      `json:"clearData"`
}

type DataSys struct {
	AdminPassword string `json:"adminPassword"`
}

type DataDb struct {
	Host     string `json:"host"  validate:"required"`
	Port     string `json:"port"  validate:"required"`
	UserName string `json:"userName" validate:"required"`
	Password string `json:"password"  validate:"required"`
	DBName   string `json:"dbName" validate:"required"`
	LogMode  bool   `json:"logMode"`
}
type DataCache struct {
	Host     string `json:"host"  validate:"required"`
	Port     string `json:"port"  validate:"required"`
	Password string `json:"password"`
	PoolSize int    `json:"poolSize"`
	DB       int    `json:"db"`
}
