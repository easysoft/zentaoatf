package ui

import "sync"

var ModuleTabs []string

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

		ModuleTabs = make([]string, 0)
		ModuleTabs = append(ModuleTabs, "testing", "projects", "settings")
	})
}
