package helper

import "strings"

func ToSlugFormat(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	return strings.Join(strings.Split(s, " "), "-")
}
