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
	return NewFilter(TextField, func(value string) (bool, map[string]any) {
		return value == match, nil
	})
}

// Создает новый фильтр, срабатывающий только при совпадении с регулярным выражением.
//
// Требует поле фильтра: [TextField]
//
// Возвращаемые побочные значения:
//   - `command_regex_filter_groups` - группы захвата в регулярном выражении.
func TextRegexCompiled(re *regexp.Regexp) Filter {
	return NewFilter(TextField, func(value string) (bool, map[string]any) {
		matches := re.FindStringSubmatch(value)
		if matches == nil {
			return false, nil
		}
		groups := []string{}
		if len(matches) > 1 {
			groups = append(groups, matches[1:]...)
		}
		return true, map[string]any{"command_regex_filter_groups": groups}
	})
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
	return NewFilter(TextField, func(value string) (bool, map[string]any) {
		if !strings.HasPrefix(value, prefix) {
			return false, nil
		}
		remainder := value[len(prefix):]
		return true, map[string]any{"prefix_filter_remainder": remainder}
	})
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
	return NewFilter(TextField, func(value string) (bool, map[string]any) {
		for _, p := range prefixList {
			if strings.HasPrefix(value, p) {
				remainder := value[len(p):]
				return true, map[string]any{
					"prefix_filter_match":     p,
					"prefix_filter_remainder": remainder,
				}
			}
		}
		return false, nil
	})
}
