package filters

import "txts.su/grfilt"

type extractableFilter[T any] struct {
	field   FieldDescriptor
	eval    func(T) bool
	extract func(T) MatchResult
}

func (f *extractableFilter[T]) Eval(input map[FieldDescriptor]any) bool {
	value, ok := input[f.field].(T)
	if !ok {
		return false
	}
	return f.eval(value)
}

func (f *extractableFilter[T]) Extract(input map[FieldDescriptor]any) MatchResult {
	value, ok := input[f.field].(T)
	if !ok {
		return nil
	}
	return f.extract(value)
}

func (f *extractableFilter[T]) Not() Filter {
	return leafFilter{inner: grfilt.Not(f)}
}

func (f *extractableFilter[T]) And(otherFilters ...Filter) Filter {
	return And(append([]Filter{f}, otherFilters...)...)
}

func (f *extractableFilter[T]) Or(otherFilters ...Filter) Filter {
	return Or(append([]Filter{f}, otherFilters...)...)
}

func NewExtractableFilter[T any](field FieldDescriptor, eval func(T) bool, extract func(T) MatchResult) Filter {
	return &extractableFilter[T]{
		field:   field,
		eval:    eval,
		extract: extract,
	}
}
