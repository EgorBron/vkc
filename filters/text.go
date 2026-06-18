package filters

import (
	"regexp"
	"strings"
)

const (
	TextField FieldDescriptor = "text"
)

// Создает новый фильтр, срабатывающий только при полном совпадении строки.
//
// Требует поле фильтра: [TextField]
func Text(match string) Filter {
	return NewFilter(TextField, func(value string) bool {
		return value == match
	})
}

// Создает новый фильтр, срабатывающий только при совпадении с регулярным выражением.
//
// Требует поле фильтра: [TextField]
//
// Возвращаемые побочные значения:
//   - `command_regex_filter_groups` - группы захвата в регулярном выражении.
func TextRegexCompiled(re *regexp.Regexp) Filter {
	return NewExtractableFilter(TextField,
		func(value string) bool {
			return re.FindStringSubmatch(value) != nil
		},
		func(value string) MatchResult {
			matches := re.FindStringSubmatch(value)
			groups := []string{}
			if len(matches) > 1 {
				groups = append(groups, matches[1:]...)
			}
			return MatchResult{"command_regex_filter_groups": groups}
		},
	)
}

// Создает новый фильтр, срабатывающий только при совпадении с регулярным выражением в виде строки.
// При вычислении может вызвать панику, если регулярное выражение сформировано неверно.
//
// Требует поле фильтра: [TextField]
//
// Возвращаемые побочные значения:
//   - `command_regex_filter_groups` - группы захвата в регулярном выражении.
func TextRegex(pattern string) Filter {
	re := regexp.MustCompile(pattern)
	return TextRegexCompiled(re)
}

// Создает новый фильтр, срабатывающий только при совпадении по префиксу.
//
// Требует поле фильтра: [TextField]
//
// Возвращаемые побочные значения:
//   - `prefix_filter_remainder` - остаток после исключения префикса.
func TextPrefix(prefix string) Filter {
	return NewExtractableFilter(TextField,
		func(value string) bool {
			return strings.HasPrefix(value, prefix)
		},
		func(value string) MatchResult {
			return MatchResult{"prefix_filter_remainder": value[len(prefix):]}
		},
	)
}

func matchingTextPrefix(value string, prefixes []string) (prefix string, matched bool) {
	for _, p := range prefixes {
		if strings.HasPrefix(value, p) {
			return p, true
		}
	}
	return "", false
}

// Создает новый фильтр, срабатывающий только при совпадении с любым префиксом из среза.
// Префиксы, находящиеся в срезе первее, при поиске имеют приоритет над префиксами после них.
//
// Требует поле фильтра: [TextField]
//
// Возвращаемые побочные значения:
//   - `prefix_filter_match` - совпавший префикс;
//   - `prefix_filter_remainder` - остаток после исключения префикса.
func TextPrefixListOf(prefixList []string) Filter {
	return NewExtractableFilter(TextField,
		func(value string) bool {
			_, matched := matchingTextPrefix(value, prefixList)
			return matched
		},
		func(value string) MatchResult {
			prefix, matched := matchingTextPrefix(value, prefixList)
			if !matched {
				return nil
			}
			return MatchResult{
				"prefix_filter_match":     prefix,
				"prefix_filter_remainder": value[len(prefix):],
			}
		},
	)
}
