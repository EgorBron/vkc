package filters

import "txts.su/grfilt"

type Filter = grfilt.Filter
type FieldDescriptor = grfilt.FieldDescriptor

type MatchResult map[string]any

type Extractor interface {
	Extract(input map[FieldDescriptor]any) MatchResult
}
