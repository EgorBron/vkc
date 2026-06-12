package filters

// Определяет поведение фильтра при его вычислении.
//
//  - (по умолчанию) если установлено в `false`, фильтр считается успешно выполненным, только если запрашиваемое поле было во входных данных;
//  - если установлено в `true`, фильтр считается успешно выполненным, даже если запрашиваемое поле отсутсвовало во входных данных. Такой фильтр стоит воспринимать как "необязательный".
var FiltersReturnTrueIfFieldNotFound = false

// Представляет самую малую единицу фильтра - пару запрашиваемого поля и предиката
type FilterItem[T any] struct {
	// Описание запрашиваемого поля
	fieldName FieldDescriptor
	// Предикат фильтра. Может вернуть побочные значения, которые могут использоваться для понимания результата работы предиката.
	fn func(T) (bool, map[string]any)
}

// Комбинирует текущий фильтр с переданными в композицию фильтров (см. [CompositeFilter])
func (fi *FilterItem[T]) And(otherFilters ...Filter) Filter {
	return combineFilters([]Filter{fi}, otherFilters)
}

// Вычисляет значение фильтра.
//
// Если при вычислении запрашиваемое поле было найдено во входных данных, предикат выполняется. В противном случае возвращается значение переменной [FiltersReturnTrueIfFieldNotFound].
//
// Предикат может вернуть побочные значения, которые могут использоваться для понимания результата работы фильтра.
func (fi *FilterItem[T]) Eval(input map[FieldDescriptor]any) (bool, map[string]any) {
	val, ok := input[fi.fieldName].(T)
	if !ok {
		return FiltersReturnTrueIfFieldNotFound, nil
	}

	return fi.fn(val)
}

// Создает новый фильтр
func NewFilter[T any](fieldName FieldDescriptor, fn func(value T) (result bool, sideValues map[string]any)) Filter {
	return &FilterItem[T]{
		fieldName: fieldName,
		fn:        fn,
	}
}
