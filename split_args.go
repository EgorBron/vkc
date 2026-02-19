package vkc

import (
	"strings"
)

// Разбиение строки на аргументы по пробелам.
func SplitArgs(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return []string{}
	}
	return strings.Fields(s)
}
