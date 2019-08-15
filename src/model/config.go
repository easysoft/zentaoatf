package model

type Config struct {
	Url        string
	EntityType string
	EntityVal  string
	ProductId  int
	ProjectId  int
	LangType   string
	SingleFile bool

	Account  string
	Password string

	ProjectName string
}
