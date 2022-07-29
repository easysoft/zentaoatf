package model

type Statistic struct {
	BaseModel

	Path     string `json:"path" gorm:"unique"`
	Total    int    `json:"total"`
	Succ     int    `json:"succ"`
	Fail     int    `json:"fail"`
	SuccLogs string `json:"succ_logs"`
	FailLogs string `json:"fail_logs"`
}

func (Statistic) TableName() string {
	return "biz_statistics"
}
