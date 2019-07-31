package ui

import "sync"

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
