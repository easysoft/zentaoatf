package ui

import "sync"

const (
	LeftWidth = 36
)

var Tabs []string

var ViewMap map[string][]string

//var EventMap map[string][][]interface{}

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
		//EventMap = map[string][][]interface{}{
		//	"root": make([][]interface{}, 0, 2),
		//	"testingTab": make([][]interface{}, 0, 2),
		//	"projectsTab": make([][]interface{}, 0, 2),
		//	"settingsTab": make([][]interface{}, 0, 2),
		//	"import": make([][]interface{}, 0, 2),
		//}

		Tabs = make([]string, 0)
		Tabs = append(Tabs, "testing", "projects", "settings")
	})
}
