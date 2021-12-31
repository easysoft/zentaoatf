package model

type Product struct {
	BaseModel

	Name string `json:"name"`
	Desc string `json:"desc" gorm:"column:descr"`

	SchemaId uint `json:"schemaId"`
	OrgId    uint `json:"orgId"`
}

func (Product) TableName() string {
	return "biz_product"
}
