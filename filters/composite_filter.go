package filters

import "maps"

// Представляет композицию нескольких фильтров.
//
// Фильтр срабатывает только если все входящие в него фильтры срабатывают успешно (логическое И).
type CompositeFilter struct {
	filters []Filter
}

// Создает новый композитный фильтр из переданных фильтров
func New(items ...Filter) *CompositeFilter {
	return &CompositeFilter{filters: items}
}

// Добавляет дополнительные фильтры к текущим и возвращает новый композитный фильтр
func (cf *CompositeFilter) And(otherFilters ...Filter) Filter {
	return combineFilters(cf.filters, otherFilters)
}

// Вычисляет результат композитного фильтра путем последовательной проверки всех входящих фильтров.
//
// Возвращает false если любой из фильтров не прошел проверку.
//
// Побочные значения от всех фильтров объединяются в итоговый словарь.
func (cf *CompositeFilter) Eval(input map[FieldDescriptor]any) (bool, map[string]any) {
	sideValues := make(map[string]any)

	for _, f := range cf.filters {
		ok, side := f.Eval(input)
		if !ok {
			return false, sideValues
		}

		maps.Copy(sideValues, side)
	}
	return true, sideValues
}
