package services

import (
	"errors"
	"testing"

	"github.com/gabriel/ttrpg-toolkit/backend/internal/models"
)

func TestValidateSkillConfig_MultiTierRequiresTiers(t *testing.T) {
	cfg := models.SkillConfig{}
	if err := validateSkillConfig(models.SkillTypeMultiTier, &cfg); !errors.Is(err, ErrInvalidSkill) {
		t.Fatalf("expected ErrInvalidSkill, got %v", err)
	}
}

func TestValidateSkillConfig_MultiTierRequiresTierName(t *testing.T) {
	cfg := models.SkillConfig{Tiers: []models.SkillTierEntry{{TierName: " ", NumericBacking: 1}}}
	if err := validateSkillConfig(models.SkillTypeMultiTier, &cfg); !errors.Is(err, ErrInvalidSkill) {
		t.Fatalf("expected ErrInvalidSkill, got %v", err)
	}
}

func TestValidateSkillConfig_NumericMinMax(t *testing.T) {
	cfg := models.SkillConfig{Min: 10, Max: 5}
	if err := validateSkillConfig(models.SkillTypeNumeric, &cfg); !errors.Is(err, ErrInvalidSkill) {
		t.Fatalf("expected ErrInvalidSkill, got %v", err)
	}
}

func TestSanitizeSkillConfig_SpecializationBonusDefault(t *testing.T) {
	cfg := models.SkillConfig{AllowSpecializations: true, SpecializationBonus: 0}
	sanitizeSkillConfig(models.SkillTypeBinary, &cfg)
	if cfg.SpecializationBonus != 2 {
		t.Fatalf("expected default bonus 2, got %d", cfg.SpecializationBonus)
	}
}

func TestSanitizeSkillConfig_ClearsBonusWhenDisabled(t *testing.T) {
	cfg := models.SkillConfig{AllowSpecializations: false, SpecializationBonus: 5}
	sanitizeSkillConfig(models.SkillTypeBinary, &cfg)
	if cfg.SpecializationBonus != 0 {
		t.Fatalf("expected bonus cleared, got %d", cfg.SpecializationBonus)
	}
}

func TestSanitizeSkillConfig_StripsNumericFieldsForMultiTier(t *testing.T) {
	cfg := models.SkillConfig{
		Tiers: []models.SkillTierEntry{{TierName: "Trained", NumericBacking: 2}},
		Min:   1,
		Max:   10,
	}
	sanitizeSkillConfig(models.SkillTypeMultiTier, &cfg)
	if cfg.Min != 0 || cfg.Max != 0 {
		t.Fatalf("expected min/max cleared for multi_tier")
	}
}
