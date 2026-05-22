package services

import (
	"fmt"
	"strings"

	"github.com/gabriel/ttrpg-toolkit/backend/internal/models"
)

func sanitizeActionEconomyConfig(cfg *models.ActionEconomyConfig) {
	if cfg == nil {
		return
	}
	if cfg.TurnStructure != models.TurnStructureCustom {
		cfg.CustomTurnStructureName = ""
	}
	if cfg.TurnStructure != models.TurnStructureTokenBased {
		cfg.TokenTurn = nil
	}
	if cfg.RoundTimeDefinition != models.RoundTimeCustom {
		cfg.CustomRoundTimeDefinition = ""
	}
	if cfg.InitiativeSystem != models.InitiativeSystemStandardRolled &&
		cfg.InitiativeSystem != models.InitiativeSystemSideBasedIndividual {
		cfg.InitiativeExpression = ""
	}
	if cfg.InitiativeSystem != models.InitiativeSystemStatic {
		cfg.StaticInitiativeValue = ""
	}
	switch cfg.SystemType {
	case models.SystemTypePointPool:
		cfg.ActionSlots = nil
	case models.SystemTypeTypeSlots:
		cfg.PointPool = models.PointPoolConfig{}
	}
	for i := range cfg.ActionSlots {
		slot := &cfg.ActionSlots[i]
		if slot.CarryOver != models.CarryOverConvert {
			slot.ConvertTarget = ""
		}
		if !slot.IsReaction {
			slot.ReactionTrigger = ""
		}
		if !slot.IsFreeAction {
			slot.FreeActionLimits = ""
		}
		if slot.Allowance == models.AllowanceUnlimited {
			slot.AllowanceScope = ""
		}
	}
}

func validateActionEconomyConfig(cfg *models.ActionEconomyConfig) error {
	if cfg == nil {
		return fmt.Errorf("%w: config is required", ErrInvalidActionEconomy)
	}
	if !models.IsAllowedTurnStructure(cfg.TurnStructure) {
		return fmt.Errorf("%w: invalid turn_structure", ErrInvalidActionEconomy)
	}
	if cfg.TurnStructure == models.TurnStructureCustom && strings.TrimSpace(cfg.CustomTurnStructureName) == "" {
		return fmt.Errorf("%w: custom_turn_structure_name is required when turn_structure is custom", ErrInvalidActionEconomy)
	}
	if cfg.TurnStructure == models.TurnStructureTokenBased {
		if cfg.TokenTurn == nil {
			return fmt.Errorf("%w: token_turn is required when turn_structure is token_based", ErrInvalidActionEconomy)
		}
		if cfg.TokenTurn.TokensPerRound < 1 {
			return fmt.Errorf("%w: tokens_per_round must be at least 1", ErrInvalidActionEconomy)
		}
		if !models.IsAllowedTokenRefreshOn(cfg.TokenTurn.RefreshOn) {
			return fmt.Errorf("%w: invalid token_turn.refresh_on", ErrInvalidActionEconomy)
		}
	}
	if !models.IsAllowedSystemType(cfg.SystemType) {
		return fmt.Errorf("%w: invalid system_type", ErrInvalidActionEconomy)
	}
	if !models.IsAllowedRoundTimeDefinition(cfg.RoundTimeDefinition) {
		return fmt.Errorf("%w: invalid round_time_definition", ErrInvalidActionEconomy)
	}
	if cfg.RoundTimeDefinition == models.RoundTimeCustom && strings.TrimSpace(cfg.CustomRoundTimeDefinition) == "" {
		return fmt.Errorf("%w: custom_round_time_definition is required when round_time_definition is custom", ErrInvalidActionEconomy)
	}
	if !models.IsAllowedCombatTimeTrackingMode(cfg.CombatTimeTrackingMode) {
		return fmt.Errorf("%w: invalid combat_time_tracking_mode", ErrInvalidActionEconomy)
	}
	if !models.IsAllowedInitiativeSystem(cfg.InitiativeSystem) {
		return fmt.Errorf("%w: invalid initiative_system", ErrInvalidActionEconomy)
	}
	if !models.IsAllowedInitiativePersistence(cfg.InitiativePersistence) {
		return fmt.Errorf("%w: invalid initiative_persistence", ErrInvalidActionEconomy)
	}
	if !models.IsAllowedTieBreaking(cfg.TieBreaking) {
		return fmt.Errorf("%w: invalid tie_breaking", ErrInvalidActionEconomy)
	}
	if cfg.InitiativeSystem == models.InitiativeSystemStandardRolled ||
		cfg.InitiativeSystem == models.InitiativeSystemSideBasedIndividual {
		if strings.TrimSpace(cfg.InitiativeExpression) == "" {
			return fmt.Errorf("%w: initiative_expression is required", ErrInvalidActionEconomy)
		}
	}
	if cfg.InitiativeSystem == models.InitiativeSystemStatic && strings.TrimSpace(cfg.StaticInitiativeValue) == "" {
		return fmt.Errorf("%w: static_initiative_value is required when initiative_system is static", ErrInvalidActionEconomy)
	}
	switch cfg.SystemType {
	case models.SystemTypePointPool:
		if err := validatePointPoolConfig(&cfg.PointPool); err != nil {
			return err
		}
	case models.SystemTypeTypeSlots:
		if err := validateActionSlots(cfg.ActionSlots); err != nil {
			return err
		}
	}
	return nil
}

func validatePointPoolConfig(pool *models.PointPoolConfig) error {
	if pool == nil {
		return fmt.Errorf("%w: point_pool is required", ErrInvalidActionEconomy)
	}
	if pool.PointsPerPool < 0 {
		return fmt.Errorf("%w: points_per_pool cannot be negative", ErrInvalidActionEconomy)
	}
	if !models.IsAllowedRefreshScope(pool.RefreshScope) {
		return fmt.Errorf("%w: invalid point_pool.refresh_scope", ErrInvalidActionEconomy)
	}
	if len(pool.ActionCostTable) == 0 {
		return fmt.Errorf("%w: at least one action cost row is required", ErrInvalidActionEconomy)
	}
	for i, row := range pool.ActionCostTable {
		if strings.TrimSpace(row.Name) == "" {
			return fmt.Errorf("%w: action_cost_table[%d].name is required", ErrInvalidActionEconomy, i)
		}
		if row.Cost < 0 {
			return fmt.Errorf("%w: action_cost_table[%d].cost cannot be negative", ErrInvalidActionEconomy, i)
		}
	}
	return nil
}

func validateActionSlots(slots []models.ActionSlotEntry) error {
	if len(slots) == 0 {
		return fmt.Errorf("%w: at least one action slot is required", ErrInvalidActionEconomy)
	}
	names := make(map[string]struct{}, len(slots))
	for _, slot := range slots {
		n := strings.TrimSpace(slot.Name)
		if n != "" {
			names[n] = struct{}{}
		}
	}
	for i, slot := range slots {
		if strings.TrimSpace(slot.Name) == "" {
			return fmt.Errorf("%w: action_slots[%d].name is required", ErrInvalidActionEconomy, i)
		}
		if slot.Allowance != models.AllowanceUnlimited && slot.Allowance < 1 {
			return fmt.Errorf("%w: action_slots[%d].allowance must be at least 1 or unlimited (-1)", ErrInvalidActionEconomy, i)
		}
		if slot.Allowance != models.AllowanceUnlimited {
			if !models.IsAllowedRefreshScope(slot.AllowanceScope) {
				return fmt.Errorf("%w: action_slots[%d].allowance_scope is required", ErrInvalidActionEconomy, i)
			}
		}
		if !models.IsAllowedCarryOver(slot.CarryOver) {
			return fmt.Errorf("%w: invalid action_slots[%d].carry_over", ErrInvalidActionEconomy, i)
		}
		if slot.CarryOver == models.CarryOverConvert && strings.TrimSpace(slot.ConvertTarget) == "" {
			return fmt.Errorf("%w: action_slots[%d].convert_target is required when carry_over is convert", ErrInvalidActionEconomy, i)
		}
		if slot.IsReaction && strings.TrimSpace(slot.ReactionTrigger) == "" {
			return fmt.Errorf("%w: action_slots[%d].reaction_trigger is required when is_reaction is true", ErrInvalidActionEconomy, i)
		}
		if slot.IsFreeAction && strings.TrimSpace(slot.FreeActionLimits) == "" {
			return fmt.Errorf("%w: action_slots[%d].free_action_limits is required when is_free_action is true", ErrInvalidActionEconomy, i)
		}
		for j, combo := range slot.Combos {
			if strings.TrimSpace(combo.ComboName) == "" {
				return fmt.Errorf("%w: action_slots[%d].combos[%d].combo_name is required", ErrInvalidActionEconomy, i, j)
			}
			if len(combo.ComponentNames) == 0 {
				return fmt.Errorf("%w: action_slots[%d].combos[%d] requires at least one component", ErrInvalidActionEconomy, i, j)
			}
			for k, comp := range combo.ComponentNames {
				if _, ok := names[strings.TrimSpace(comp)]; !ok {
					return fmt.Errorf("%w: action_slots[%d].combos[%d].component_names[%d] does not match a defined action name", ErrInvalidActionEconomy, i, j, k)
				}
			}
		}
	}
	return nil
}

func initiativeNeedsFormula(system string) bool {
	return system == models.InitiativeSystemStandardRolled ||
		system == models.InitiativeSystemSideBasedIndividual
}
