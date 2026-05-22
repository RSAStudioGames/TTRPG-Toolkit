package models

// Allowed resolution_type values.
const (
	ResolutionTypeSingleDieModifier          = "single_die_modifier"
	ResolutionTypeDicePoolCountSuccesses     = "dice_pool_count_successes"
	ResolutionTypeDicePoolTakeHighestLowest  = "dice_pool_take_highest_lowest"
	ResolutionTypeDicePoolSum                = "dice_pool_sum"
	ResolutionTypeStepDice                    = "step_dice"
	ResolutionTypeCardDraw                    = "card_draw"
	ResolutionTypeCoinFlip                    = "coin_flip"
	ResolutionTypeNarrative                   = "narrative"
	ResolutionTypeCustom                      = "custom"
)

// Allowed success_determination.method values.
const (
	SuccessMethodTargetNumber              = "target_number"
	SuccessMethodOpposedRoll               = "opposed_roll"
	SuccessMethodSuccessThresholdLadder    = "success_threshold_ladder"
	SuccessMethodSuccessCounting           = "success_counting"
	SuccessMethodBinaryPassFail            = "binary_pass_fail"
	SuccessMethodMarginOfSuccessFailure    = "margin_of_success_failure"
	SuccessMethodShiftsRaises              = "shifts_raises"
	SuccessMethodCustomTiers               = "custom_tiers"
)

// Allowed ladder_tier.operator values.
const (
	LadderOperatorLTE = "<="
	LadderOperatorLT  = "<"
	LadderOperatorGTE = ">="
	LadderOperatorGT  = ">"
	LadderOperatorEQ  = "=="
)

// Allowed advantage_disadvantage.mechanic_type values.
const (
	MechanicTypeKeepHighestN = "keep_highest_n"
	MechanicTypeKeepLowestN  = "keep_lowest_n"
	MechanicTypeAddModifier  = "add_modifier"
	MechanicTypeReroll       = "reroll"
	MechanicTypeShiftTier    = "shift_tier"
)

var allowedResolutionTypes = map[string]struct{}{
	ResolutionTypeSingleDieModifier:         {},
	ResolutionTypeDicePoolCountSuccesses:    {},
	ResolutionTypeDicePoolTakeHighestLowest: {},
	ResolutionTypeDicePoolSum:               {},
	ResolutionTypeStepDice:                  {},
	ResolutionTypeCardDraw:                  {},
	ResolutionTypeCoinFlip:                  {},
	ResolutionTypeNarrative:                 {},
	ResolutionTypeCustom:                    {},
}

var allowedSuccessMethods = map[string]struct{}{
	SuccessMethodTargetNumber:           {},
	SuccessMethodOpposedRoll:            {},
	SuccessMethodSuccessThresholdLadder: {},
	SuccessMethodSuccessCounting:       {},
	SuccessMethodBinaryPassFail:         {},
	SuccessMethodMarginOfSuccessFailure: {},
	SuccessMethodShiftsRaises:           {},
	SuccessMethodCustomTiers:            {},
}

var allowedLadderOperators = map[string]struct{}{
	LadderOperatorLTE: {},
	LadderOperatorLT:  {},
	LadderOperatorGTE: {},
	LadderOperatorGT:  {},
	LadderOperatorEQ:  {},
}

var allowedMechanicTypes = map[string]struct{}{
	MechanicTypeKeepHighestN: {},
	MechanicTypeKeepLowestN:  {},
	MechanicTypeAddModifier:  {},
	MechanicTypeReroll:       {},
	MechanicTypeShiftTier:    {},
}

// IsAllowedResolutionType reports whether v is a valid resolution_type.
func IsAllowedResolutionType(v string) bool {
	_, ok := allowedResolutionTypes[v]
	return ok
}

// IsAllowedSuccessMethod reports whether v is a valid success_determination.method.
func IsAllowedSuccessMethod(v string) bool {
	_, ok := allowedSuccessMethods[v]
	return ok
}

// IsAllowedLadderOperator reports whether v is a valid ladder_tier.operator.
func IsAllowedLadderOperator(v string) bool {
	_, ok := allowedLadderOperators[v]
	return ok
}

// IsAllowedMechanicType reports whether v is a valid advantage_disadvantage.mechanic_type.
func IsAllowedMechanicType(v string) bool {
	_, ok := allowedMechanicTypes[v]
	return ok
}
