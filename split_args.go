package vkc

import (
	"strings"

	"txts.su/vkc/filters"
)

// Разбиение строки на аргументы по пробелам.
func SplitArgs(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return []string{}
	}
	return strings.Fields(s)
}

// Формирует аргументы команды из результатов извлечения фильтра.
func ArgumentsFromMatch(m filters.MatchResult) []string {
	if m == nil {
		return []string{}
	}
	if rem, ok := m["prefix_filter_remainder"].(string); ok {
		return SplitArgs(rem)
	}
	if groups, ok := m["command_regex_filter_groups"].([]string); ok {
		return groups
	}
	return []string{}
}
