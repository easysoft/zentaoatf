package ui

import "sync"

const (
	LeftWidth = 36
)

var Tabs []string

var ViewMap map[string][]string
var EventMap map[string][][]interface{}

func init() {
	var once sync.Once
	once.Do(func() {
		Tabs = make([]string, 0)
		Tabs = append(Tabs, "testing", "projects", "settings")
	})
}
