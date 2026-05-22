package services

import (
	"errors"
	"testing"

	"github.com/gabriel/ttrpg-toolkit/backend/internal/models"
)

func validResolutionConfig() models.ResolutionConfig {
	return models.ResolutionConfig{
		ResolutionType: models.ResolutionTypeSingleDieModifier,
		RollExpression: "1d20 + {modifier}",
		SuccessDetermination: models.SuccessDetermination{
			Method:                models.SuccessMethodTargetNumber,
			ThresholdLadder:       nil,
			DefaultTargetVariable: "target",
		},
		CriticalMechanics: models.CriticalMechanics{
			EnableCritSuccess: false,
			EnableCritFailure: false,
		},
		AdvantageDisadvantage: []models.AdvantageDisadvantageEntry{},
	}
}

func TestValidateResolutionConfig_valid(t *testing.T) {
	cfg := validResolutionConfig()
	if err := validateResolutionConfig(&cfg); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestValidateResolutionConfig_invalidResolutionType(t *testing.T) {
	cfg := validResolutionConfig()
	cfg.ResolutionType = "d20_roll"
	if err := validateResolutionConfig(&cfg); !errors.Is(err, ErrInvalidResolution) {
		t.Fatalf("expected ErrInvalidResolution, got %v", err)
	}
}

func TestValidateResolutionConfig_invalidMethod(t *testing.T) {
	cfg := validResolutionConfig()
	cfg.SuccessDetermination.Method = "fixed_dc"
	if err := validateResolutionConfig(&cfg); !errors.Is(err, ErrInvalidResolution) {
		t.Fatalf("expected ErrInvalidResolution, got %v", err)
	}
}

func TestValidateResolutionConfig_ladderRequired(t *testing.T) {
	cfg := validResolutionConfig()
	cfg.SuccessDetermination.Method = models.SuccessMethodSuccessThresholdLadder
	cfg.SuccessDetermination.ThresholdLadder = nil
	if err := validateResolutionConfig(&cfg); !errors.Is(err, ErrInvalidResolution) {
		t.Fatalf("expected ErrInvalidResolution, got %v", err)
	}
}

func TestValidateResolutionConfig_invalidLadderOperator(t *testing.T) {
	cfg := validResolutionConfig()
	cfg.SuccessDetermination.Method = models.SuccessMethodSuccessThresholdLadder
	cfg.SuccessDetermination.ThresholdLadder = []models.LadderTier{
		{Label: "Success", Operator: "!=", Value: 10},
	}
	if err := validateResolutionConfig(&cfg); !errors.Is(err, ErrInvalidResolution) {
		t.Fatalf("expected ErrInvalidResolution, got %v", err)
	}
}

func TestValidateResolutionConfig_validLadderOperator(t *testing.T) {
	cfg := validResolutionConfig()
	cfg.SuccessDetermination.Method = models.SuccessMethodSuccessThresholdLadder
	cfg.SuccessDetermination.ThresholdLadder = []models.LadderTier{
		{Label: "Success", Operator: ">=", Value: 10},
	}
	if err := validateResolutionConfig(&cfg); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestValidateResolutionConfig_invalidMechanicType(t *testing.T) {
	cfg := validResolutionConfig()
	cfg.AdvantageDisadvantage = []models.AdvantageDisadvantageEntry{
		{Name: "Advantage", MechanicType: "advantage"},
	}
	if err := validateResolutionConfig(&cfg); !errors.Is(err, ErrInvalidResolution) {
		t.Fatalf("expected ErrInvalidResolution, got %v", err)
	}
}

func TestValidateResolutionConfig_critTriggers(t *testing.T) {
	cfg := validResolutionConfig()
	cfg.CriticalMechanics.EnableCritSuccess = true
	if err := validateResolutionConfig(&cfg); !errors.Is(err, ErrInvalidResolution) {
		t.Fatalf("expected ErrInvalidResolution for missing crit_success_trigger, got %v", err)
	}
	cfg.CriticalMechanics.CritSuccessTrigger = "natural_20"
	if err := validateResolutionConfig(&cfg); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestValidateResolutionConfig_customRequiresName(t *testing.T) {
	cfg := validResolutionConfig()
	cfg.ResolutionType = models.ResolutionTypeCustom
	cfg.CustomParadigmName = ""
	if err := validateResolutionConfig(&cfg); !errors.Is(err, ErrInvalidResolution) {
		t.Fatalf("expected ErrInvalidResolution, got %v", err)
	}
	cfg.CustomParadigmName = "  My System  "
	if err := validateResolutionConfig(&cfg); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestSanitizeResolutionConfig_stripsCustomName(t *testing.T) {
	cfg := validResolutionConfig()
	cfg.CustomParadigmName = "Orphan"
	sanitizeResolutionConfig(&cfg)
	if cfg.CustomParadigmName != "" {
		t.Fatalf("expected empty custom_paradigm_name, got %q", cfg.CustomParadigmName)
	}
	cfg.ResolutionType = models.ResolutionTypeCustom
	cfg.CustomParadigmName = "Keep Me"
	sanitizeResolutionConfig(&cfg)
	if cfg.CustomParadigmName != "Keep Me" {
		t.Fatalf("expected name preserved for custom type, got %q", cfg.CustomParadigmName)
	}
}

func TestSaveResolutionConfig_invalidFormula(t *testing.T) {
	svc := NewMechanicsService(nil, nil, NewFormulaService())
	cfg := validResolutionConfig()
	cfg.RollExpression = "2d"
	_, err := svc.SaveResolutionConfig("00000000-0000-0000-0000-000000000001", cfg)
	var formulaErr *InvalidFormulaError
	if !errors.As(err, &formulaErr) {
		t.Fatalf("expected InvalidFormulaError, got %v", err)
	}
	if len(formulaErr.Errors) == 0 {
		t.Fatal("expected formula error details")
	}
}
