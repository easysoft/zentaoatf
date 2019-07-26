package ui

import (
	"strings"
)

func GetNextView(name string, views []string) string {
	i := 0
	found := false
	for true {
		if name == views[i] {
			found = true
			i++
			i = i % len(views)
			continue
		}

		if found {
			if strings.Index(views[i], "Input") > -1 {
				return views[i]
			}
		}

		i++
		i = i % len(views)
	}

	return ""
}
