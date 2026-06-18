package filters

import "maps"

func mergeResults(dst MatchResult, src MatchResult) {
	if src == nil {
		return
	}
	maps.Copy(dst, src)
}
