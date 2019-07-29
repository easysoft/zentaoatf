package ui

import "sync"

const (
	LeftWidth = 36
	MinWidth  = 130
	MinHeight = 36
)

var Tabs []string

var ViewMap map[string][]string

func init() {
	var once sync.Once
	once.Do(func() {
		ViewMap = map[string][]string{
			"root":        {},
			"testingTab":  {},
			"projectsTab": {},
			"settingsTab": {},
			"import":      {},
		}

		Tabs = make([]string, 0)
		Tabs = append(Tabs, "testing", "projects", "settings")
	})
}
