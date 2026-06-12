package filters

import "slices"

// Создает новый фильтр, срабатывающий только при совпадении с любым элементом из среза.
// Элементы, стоящие в срезе первее, имеют приоритет при поиске.
//
// Требует указания поля фильтра.
//
// Возвращаемые побочные значения:
//   - `list_filter_match_index` - индекс совпавшего элемента или -1 при отсутствии совпадений.
func InList[T comparable](field FieldDescriptor, list []T) Filter {
	return NewFilter(field, func(value T) (bool, map[string]any) {
		idx := slices.Index(list, value)
		return idx >= 0, map[string]any{"list_filter_match_index": idx}
	})
}
