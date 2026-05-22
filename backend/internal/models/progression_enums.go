package models

const (
	ProgressionParadigmLevelBased    = "level_based"
	ProgressionParadigmPointBuy        = "point_buy"
	ProgressionParadigmMilestone       = "milestone"
	ProgressionParadigmLifepath        = "lifepath"
	ProgressionParadigmProgressTrack   = "progress_track"
	ProgressionParadigmNoAdvancement   = "no_advancement"
	ProgressionParadigmCustom          = "custom"
)

var allowedProgressionParadigms = map[string]struct{}{
	ProgressionParadigmLevelBased:  {},
	ProgressionParadigmPointBuy:    {},
	ProgressionParadigmMilestone:   {},
	ProgressionParadigmLifepath:    {},
	ProgressionParadigmProgressTrack: {},
	ProgressionParadigmNoAdvancement: {},
	ProgressionParadigmCustom:      {},
}

func IsAllowedProgressionParadigm(v string) bool {
	_, ok := allowedProgressionParadigms[v]
	return ok
}
