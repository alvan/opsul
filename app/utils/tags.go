package utils

import "strings"

func Tags(str string, sep string) (tags []string) {
	for _, term := range strings.Split(str, sep) {
		term = strings.TrimSpace(term)
		if term != "" {
			tags = append(tags, term)
		}
	}

	return
}
