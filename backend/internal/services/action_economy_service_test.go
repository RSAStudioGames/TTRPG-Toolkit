package services

import (
	"errors"
	"testing"

	"github.com/gabriel/ttrpg-toolkit/backend/internal/models"
)

func validActionEconomyConfig() *models.ActionEconomyConfig {
	return &models.ActionEconomyConfig{
		TurnStructure:          models.TurnStructureTurnBased,
		SystemType:             models.SystemTypePointPool,
		PointPool: models.PointPoolConfig{
			PointsPerPool: 3,
			RefreshScope:  models.RefreshScopePerTurn,
			ActionCostTable: []models.ActionCostEntry{
				{Name: "Attack", Cost: 2},
				{Name: "Move", Cost: 1},
			},
		},
		RoundTimeDefinition:    models.RoundTimeOneMinute,
		CombatTimeTrackingMode: models.CombatTimeRoundCounting,
		InitiativeSystem:       models.InitiativeSystemStandardRolled,
		InitiativePersistence:  models.InitiativePersistencePersistCombat,
		InitiativeExpression:   "1d20 + {dexterity_mod}",
		TieBreaking:            models.TieBreakingAttackerWins,
	}
}

func TestValidateActionEconomyConfig_ValidPointPool(t *testing.T) {
	if err := validateActionEconomyConfig(validActionEconomyConfig()); err != nil {
		t.Fatalf("expected valid config, got %v", err)
	}
}

func TestValidateActionEconomyConfig_ZeroCostAllowed(t *testing.T) {
	cfg := validActionEconomyConfig()
	cfg.PointPool.ActionCostTable = append(cfg.PointPool.ActionCostTable, models.ActionCostEntry{Name: "Speak", Cost: 0})
	if err := validateActionEconomyConfig(cfg); err != nil {
		t.Fatalf("expected zero cost allowed, got %v", err)
	}
}

func TestValidateActionEconomyConfig_NegativeCostRejected(t *testing.T) {
	cfg := validActionEconomyConfig()
	cfg.PointPool.ActionCostTable[0].Cost = -1
	if err := validateActionEconomyConfig(cfg); !errors.Is(err, ErrInvalidActionEconomy) {
		t.Fatalf("expected ErrInvalidActionEconomy, got %v", err)
	}
}

func TestValidateActionEconomyConfig_ComboSelfReference(t *testing.T) {
	cfg := validActionEconomyConfig()
	cfg.SystemType = models.SystemTypeTypeSlots
	cfg.PointPool = models.PointPoolConfig{}
	cfg.ActionSlots = []models.ActionSlotEntry{
		{
			Name:           "Strike",
			Allowance:      1,
			AllowanceScope: models.RefreshScopePerTurn,
			CarryOver:      models.CarryOverExpire,
			Combos: []models.ActionComboEntry{
				{ComboName: "Double Strike", ComponentNames: []string{"Strike", "Strike"}},
			},
		},
	}
	if err := validateActionEconomyConfig(cfg); err != nil {
		t.Fatalf("expected self-referencing combo allowed, got %v", err)
	}
}

func TestSanitizeActionEconomyConfig_UnlimitedStripsScope(t *testing.T) {
	cfg := &models.ActionEconomyConfig{
		SystemType: models.SystemTypeTypeSlots,
		ActionSlots: []models.ActionSlotEntry{
			{
				Name:           "Dash",
				Allowance:      models.AllowanceUnlimited,
				AllowanceScope: models.RefreshScopePerRound,
				CarryOver:      models.CarryOverExpire,
			},
		},
	}
	sanitizeActionEconomyConfig(cfg)
	if cfg.ActionSlots[0].AllowanceScope != "" {
		t.Fatalf("expected allowance_scope stripped for unlimited")
	}
}

func TestValidateActionEconomyConfig_RefreshScopeRequired(t *testing.T) {
	cfg := validActionEconomyConfig()
	cfg.PointPool.RefreshScope = ""
	if err := validateActionEconomyConfig(cfg); !errors.Is(err, ErrInvalidActionEconomy) {
		t.Fatalf("expected ErrInvalidActionEconomy, got %v", err)
	}
}
