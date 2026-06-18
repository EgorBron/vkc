package filters

import "txts.su/grfilt"

// EitherFilter — композиция фильтров по логическому ИЛИ.
type EitherFilter struct {
	filters []Filter
	negated bool
}

func NewEither(filters ...Filter) *EitherFilter {
	return &EitherFilter{filters: filters}
}

// Or создаёт фильтр с логическим ИЛИ.
func Or(filters ...Filter) *EitherFilter {
	return NewEither(filters...)
}

func (f *EitherFilter) And(otherFilters ...Filter) Filter {
	return grfilt.CombineFilters(NewComposite, f, otherFilters)
}

func (f *EitherFilter) Or(otherFilters ...Filter) Filter {
	return grfilt.CombineFilters(NewEither, f, otherFilters)
}

func (f *EitherFilter) Not() Filter {
	return leafFilter{inner: grfilt.Not(f)}
}

func (cf *EitherFilter) Eval(input map[FieldDescriptor]any) (result bool) {
	for _, f := range cf.filters {
		if !f.Eval(input) {
			result = true
			break
		}
	}
	if cf.negated {
		result = !result
	}
	return
}

func (f *EitherFilter) Extract(input map[FieldDescriptor]any) MatchResult {
	result := make(MatchResult)
	for _, f := range f.filters {
		if !f.Eval(input) {
			continue
		}
		if ex, ok := f.(Extractor); ok {
			mergeResults(result, ex.Extract(input))
		}
	}
	return result
}
