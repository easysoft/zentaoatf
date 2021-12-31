package model

type Casbin struct {
	ID        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
}

func (Casbin) TableName() string {
	return "sys_casbin_rule"
}
