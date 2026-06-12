// Пакет filters предоставляет систему компонуемых фильтров для выбора объектов по условиям.
//
// Основные концепции:
//   - Filter - интерфейс для любого фильтра с методами And() и Eval().
//   - FilterItem[T] - базовый фильтр, работающий с одним полем и предикатом.
//   - CompositeFilter - композиция нескольких фильтров (логическое И).
//   - FieldDescriptor - описание поля для фильтрации.
//
// Пакет не привязан к vkc. Фильтры можно использовать где и как угодно.
//
// Пример использования:
//
//	filter := filters.New(
//		filters.Text("hello"),
//		filters.TextPrefix("hello")
//	)
//	ok, sideValues := filter.Eval(map[filters.FieldDescriptor]any{
//		filters.TextField: "hello",
//	})
package filters

// Описание поля входных данных фильтра
type FieldDescriptor string

// Определяет интерфейс для фильтра
type Filter interface {
	// Комбинирует текущий фильтр с другими фильтрами в новый композитный фильтр (см. [CompositeFilter])
	And(otherFilters ...Filter) Filter

	// Вычисляет результат фильтра по входным данным.
	//
	// Возвращает:
	//   - result: true если фильтр успешно применен, иначе false.
	//   - sideValues: побочные значения, которые могут быть использованы для анализа результата работы фильтра.
	Eval(input map[FieldDescriptor]any) (result bool, sideValues map[string]any)
}
