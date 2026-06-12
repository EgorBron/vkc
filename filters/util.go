package filters

func combineFilters(existing []Filter, additional []Filter) *CompositeFilter {
	allFilters := make([]Filter, 0, len(existing)+len(additional))
	allFilters = append(allFilters, existing...)
	allFilters = append(allFilters, additional...)
	return New(allFilters...)
}
