package model

type Preference struct {
	Language string
	WorkDir  string

	Width  int
	Height int

	WorkHistories []WorkHistory
}

type WorkHistory struct {
	Id          string
	ProjectName string
	ProjectPath string

	EntityType string
	EntityVal  string
}
