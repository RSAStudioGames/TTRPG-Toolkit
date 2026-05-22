package services

import (
	"errors"
	"testing"

	"github.com/gabriel/ttrpg-toolkit/backend/internal/models"
)

func validNumericAttributeConfig() models.AttributeConfig {
	return models.AttributeConfig{
		NumericFormat:   models.NumericFormatInteger,
		ModifierFormula: "1d20 + {strength}",
	}
}

func TestValidateAttributeConfig_numericValid(t *testing.T) {
	cfg := validNumericAttributeConfig()
	if err := validateAttributeConfig(models.AttributeTypeNumeric, &cfg); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestValidateAttributeConfig_numericMissingFormat(t *testing.T) {
	cfg := models.AttributeConfig{}
	if err := validateAttributeConfig(models.AttributeTypeNumeric, &cfg); !errors.Is(err, ErrInvalidAttribute) {
		t.Fatalf("expected ErrInvalidAttribute, got %v", err)
	}
}

func TestValidateAttributeConfig_rankTierRequiresMap(t *testing.T) {
	cfg := models.AttributeConfig{}
	if err := validateAttributeConfig(models.AttributeTypeRankTier, &cfg); !errors.Is(err, ErrInvalidAttribute) {
		t.Fatalf("expected ErrInvalidAttribute, got %v", err)
	}
}

func TestValidateAttributeConfig_rankTierValid(t *testing.T) {
	cfg := models.AttributeConfig{
		RankMap: []models.RankMapEntry{
			{RankName: "Novice", NumericBacking: 1},
		},
	}
	if err := validateAttributeConfig(models.AttributeTypeRankTier, &cfg); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestValidateAttributeConfig_derivedRequiresCaching(t *testing.T) {
	cfg := models.AttributeConfig{
		IsDerived:         true,
		DerivationFormula: "2 + 2",
	}
	if err := validateAttributeConfig(models.AttributeTypeCustom, &cfg); !errors.Is(err, ErrInvalidAttribute) {
		t.Fatalf("expected ErrInvalidAttribute, got %v", err)
	}
}

func TestValidateAttributeConfig_derivedOnTriggerRequiresTriggers(t *testing.T) {
	cfg := models.AttributeConfig{
		IsDerived:           true,
		DerivationFormula:   "2 + 2",
		CachingRule:         models.CachingRuleOnTrigger,
		RecalculateTriggers: nil,
	}
	if err := validateAttributeConfig(models.AttributeTypeCustom, &cfg); !errors.Is(err, ErrInvalidAttribute) {
		t.Fatalf("expected ErrInvalidAttribute, got %v", err)
	}
}

func TestValidateAttributeConfig_stepDieRequiresDice(t *testing.T) {
	cfg := models.AttributeConfig{}
	if err := validateAttributeConfig(models.AttributeTypeStepDie, &cfg); !errors.Is(err, ErrInvalidAttribute) {
		t.Fatalf("expected ErrInvalidAttribute, got %v", err)
	}
}

func TestSanitizeAttributeConfig_stripsDerivedFields(t *testing.T) {
	cfg := models.AttributeConfig{
		IsDerived:           false,
		DerivationFormula:   "orphan",
		CachingRule:         models.CachingRuleLive,
		RecalculateTriggers: []string{"x"},
	}
	sanitizeAttributeConfig(&cfg)
	if cfg.DerivationFormula != "" || cfg.CachingRule != "" || len(cfg.RecalculateTriggers) > 0 {
		t.Fatal("expected derived fields cleared when not derived")
	}
}

func TestValidateAttributeFormulas_invalidModifier(t *testing.T) {
	svc := NewMechanicsService(nil, nil, NewFormulaService())
	cfg := models.AttributeConfig{ModifierFormula: "2d"}
	formulaErr := svc.validateAttributeFormulas(&cfg)
	if formulaErr == nil {
		t.Fatal("expected formula error")
	}
	var ferr *AttributeFormulaError
	if !errors.As(formulaErr, &ferr) || ferr.Field != "modifier_formula" {
		t.Fatalf("expected AttributeFormulaError modifier_formula, got %v", formulaErr)
	}
}

func TestAttributeFormulaError_message(t *testing.T) {
	e := &AttributeFormulaError{Field: "derivation_formula"}
	if e.Error() != "invalid derivation formula" {
		t.Fatalf("got %q", e.Error())
	}
}
