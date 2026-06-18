package filters

import "txts.su/grfilt"

type leafFilter struct {
	inner Filter
}

func (f leafFilter) Eval(input map[FieldDescriptor]any) bool {
	return f.inner.Eval(input)
}

func (f leafFilter) Not() Filter {
	return leafFilter{inner: grfilt.Not(f.inner)}
}

func (f leafFilter) And(otherFilters ...Filter) Filter {
	return And(append([]Filter{f}, otherFilters...)...)
}

func (f leafFilter) Or(otherFilters ...Filter) Filter {
	return Or(append([]Filter{f}, otherFilters...)...)
}

// NewFilter создаёт листовой фильтр по полю и предикату (только bool, без Extractor).
func NewFilter[T any](fieldName FieldDescriptor, fn func(value T) bool) Filter {
	return leafFilter{inner: grfilt.NewFilter(fieldName, fn)}
}
