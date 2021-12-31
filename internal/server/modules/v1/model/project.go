package model

type Project struct {
	BaseModel

	Name string `json:"name"`
	Desc string `json:"desc" gorm:"column:descr"`

	SchemaId uint `json:"schemaId"`
	OrgId    uint `json:"orgId"`

	//Products []*Product `json:"products" gorm:"many2many:biz_project_product_r;"`
}

func (Project) TableName() string {
	return "biz_project"
}
