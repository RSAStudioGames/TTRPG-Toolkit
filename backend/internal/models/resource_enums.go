package models

const (
	ResourceTypePool      = "pool"
	ResourceTypeSlotTrack = "slot_track"
	ResourceTypeCounter   = "counter"
	ResourceTypeCurrency    = "currency"
	ResourceTypeCustom      = "custom"

	ResourceFormatInteger = "integer"
	ResourceFormatFloat   = "float"
	ResourceFormatDiePool = "die_pool"
	ResourceFormatStepDie = "step_die"

	RecoveryTriggerShortRest = "short_rest"
	RecoveryTriggerLongRest  = "long_rest"
	RecoveryTriggerDawn      = "dawn"
	RecoveryTriggerCustom    = "custom"
)

var allowedResourceTypes = map[string]struct{}{
	ResourceTypePool:      {},
	ResourceTypeSlotTrack: {},
	ResourceTypeCounter:   {},
	ResourceTypeCurrency:  {},
	ResourceTypeCustom:    {},
}

var allowedResourceFormats = map[string]struct{}{
	ResourceFormatInteger: {},
	ResourceFormatFloat:   {},
	ResourceFormatDiePool: {},
	ResourceFormatStepDie: {},
}

var allowedRecoveryTriggers = map[string]struct{}{
	RecoveryTriggerShortRest: {},
	RecoveryTriggerLongRest:  {},
	RecoveryTriggerDawn:      {},
	RecoveryTriggerCustom:    {},
}

func IsAllowedResourceType(v string) bool {
	_, ok := allowedResourceTypes[v]
	return ok
}

func IsAllowedResourceFormat(v string) bool {
	_, ok := allowedResourceFormats[v]
	return ok
}

func IsAllowedRecoveryTrigger(v string) bool {
	_, ok := allowedRecoveryTriggers[v]
	return ok
}
