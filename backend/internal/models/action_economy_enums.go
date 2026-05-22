package models

const (
	AllowanceUnlimited = -1

	TurnStructureTurnBased    = "turn_based"
	TurnStructurePopcorn      = "popcorn"
	TurnStructureSideBased    = "side_based"
	TurnStructureTokenBased   = "token_based"
	TurnStructureNoFormal     = "no_formal"
	TurnStructureSimultaneous = "simultaneous"
	TurnStructureCustom       = "custom"

	SystemTypePointPool  = "point_pool"
	SystemTypeTypeSlots  = "type_slots"

	CarryOverExpire  = "expire"
	CarryOverConvert = "convert"

	RefreshScopePerTurn  = "per_turn"
	RefreshScopePerRound = "per_round"

	InitiativeSystemStandardRolled       = "standard_rolled"
	InitiativeSystemStatic               = "static"
	InitiativeSystemGroup                = "group"
	InitiativeSystemSideBasedIndividual  = "side_based_individual"
	InitiativeSystemCardBased            = "card_based"
	InitiativeSystemPopcornPass          = "popcorn_pass"

	InitiativePersistenceRerollEachRound = "reroll_each_round"
	InitiativePersistencePersistCombat   = "persist_combat"

	TieBreakingAttackerWins    = "attacker_wins"
	TieBreakingDefenderWins      = "defender_wins"
	TieBreakingRollAgain         = "roll_again"
	TieBreakingPartialSuccess    = "partial_success"
	TieBreakingStipulation       = "stipulation"

	RoundTimeSixSeconds  = "six_seconds"
	RoundTimeTenSeconds  = "ten_seconds"
	RoundTimeOneMinute   = "one_minute"
	RoundTimeNarrative   = "narrative"
	RoundTimeCustom      = "custom"

	CombatTimeRoundCounting   = "round_counting"
	CombatTimeRealTimeMapping = "real_time_mapping"
	CombatTimeNarrative       = "narrative"

	TokenRefreshPerRound = "per_round"
	TokenRefreshPerTurn  = "per_turn"
)

var (
	allowedTurnStructures = map[string]struct{}{
		TurnStructureTurnBased: {}, TurnStructurePopcorn: {}, TurnStructureSideBased: {},
		TurnStructureTokenBased: {}, TurnStructureNoFormal: {}, TurnStructureSimultaneous: {},
		TurnStructureCustom: {},
	}
	allowedSystemTypes = map[string]struct{}{
		SystemTypePointPool: {}, SystemTypeTypeSlots: {},
	}
	allowedCarryOver = map[string]struct{}{
		CarryOverExpire: {}, CarryOverConvert: {},
	}
	allowedRefreshScopes = map[string]struct{}{
		RefreshScopePerTurn: {}, RefreshScopePerRound: {},
	}
	allowedInitiativeSystems = map[string]struct{}{
		InitiativeSystemStandardRolled: {}, InitiativeSystemStatic: {},
		InitiativeSystemGroup: {}, InitiativeSystemSideBasedIndividual: {},
		InitiativeSystemCardBased: {}, InitiativeSystemPopcornPass: {},
	}
	allowedInitiativePersistence = map[string]struct{}{
		InitiativePersistenceRerollEachRound: {}, InitiativePersistencePersistCombat: {},
	}
	allowedTieBreaking = map[string]struct{}{
		TieBreakingAttackerWins: {}, TieBreakingDefenderWins: {}, TieBreakingRollAgain: {},
		TieBreakingPartialSuccess: {}, TieBreakingStipulation: {},
	}
	allowedRoundTimeDefinitions = map[string]struct{}{
		RoundTimeSixSeconds: {}, RoundTimeTenSeconds: {}, RoundTimeOneMinute: {},
		RoundTimeNarrative: {}, RoundTimeCustom: {},
	}
	allowedCombatTimeTracking = map[string]struct{}{
		CombatTimeRoundCounting: {}, CombatTimeRealTimeMapping: {}, CombatTimeNarrative: {},
	}
	allowedTokenRefresh = map[string]struct{}{
		TokenRefreshPerRound: {}, TokenRefreshPerTurn: {},
	}
)

func IsAllowedTurnStructure(v string) bool {
	_, ok := allowedTurnStructures[v]
	return ok
}

func IsAllowedSystemType(v string) bool {
	_, ok := allowedSystemTypes[v]
	return ok
}

func IsAllowedCarryOver(v string) bool {
	_, ok := allowedCarryOver[v]
	return ok
}

func IsAllowedRefreshScope(v string) bool {
	_, ok := allowedRefreshScopes[v]
	return ok
}

func IsAllowedInitiativeSystem(v string) bool {
	_, ok := allowedInitiativeSystems[v]
	return ok
}

func IsAllowedInitiativePersistence(v string) bool {
	_, ok := allowedInitiativePersistence[v]
	return ok
}

func IsAllowedTieBreaking(v string) bool {
	_, ok := allowedTieBreaking[v]
	return ok
}

func IsAllowedRoundTimeDefinition(v string) bool {
	_, ok := allowedRoundTimeDefinitions[v]
	return ok
}

func IsAllowedCombatTimeTrackingMode(v string) bool {
	_, ok := allowedCombatTimeTracking[v]
	return ok
}

func IsAllowedTokenRefreshOn(v string) bool {
	_, ok := allowedTokenRefresh[v]
	return ok
}
