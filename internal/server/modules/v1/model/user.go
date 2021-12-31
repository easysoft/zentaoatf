package model

type SysUser struct {
	BaseModel

	BaseUser
	Password string `gorm:"type:varchar(250)" json:"password" validate:"required"`
	RoleIds  []uint `gorm:"-" json:"role_ids"`
}

type BaseUser struct {
	Username string `gorm:"uniqueIndex;not null;type:varchar(60)" json:"username" validate:"required"`
	Name     string `gorm:"index;not null; type:varchar(60)" json:"name"`
	Intro    string `gorm:"not null; type:varchar(512)" json:"intro"`
	Avatar   string `gorm:"type:varchar(1024)" json:"avatar"`
}

type Avatar struct {
	Avatar string `json:"avatar"`
}

func (SysUser) TableName() string {
	return "sys_user"
}
