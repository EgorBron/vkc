package filters

import "txts.su/grfilt"

type CompositeFilter struct {
	filters []Filter
	negated bool
}

func NewComposite(items ...Filter) *CompositeFilter {
	return &CompositeFilter{filters: items}
}

func (f *CompositeFilter) Not() Filter {
	f.negated = !f.negated
	return f
}

func And(filters ...Filter) *CompositeFilter {
	return NewComposite(filters...)
}

func (f *CompositeFilter) And(otherFilters ...Filter) Filter {
	return grfilt.CombineFilters(NewComposite, f, otherFilters)
}

func (f *CompositeFilter) Or(otherFilters ...Filter) Filter {
	return grfilt.CombineFilters(NewEither, f, otherFilters)
}

func (cf *CompositeFilter) Eval(input map[FieldDescriptor]any) (result bool) {
	result = true
	for _, f := range cf.filters {
		if !f.Eval(input) {
			result = false
			break
		}
	}
	if cf.negated {
		result = !result
	}
	return
}

func (cf *CompositeFilter) Extract(input map[FieldDescriptor]any) MatchResult {
	result := make(MatchResult)
	for _, f := range cf.filters {
		if ex, ok := f.(Extractor); ok {
			mergeResults(result, ex.Extract(input))
		}
	}
	return result
}
