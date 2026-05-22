package services

import (
	"errors"
	"testing"

	"github.com/gabriel/ttrpg-toolkit/backend/internal/models"
)

func TestValidateProgressionConfig_LevelMinMax(t *testing.T) {
	cfg := &models.ProgressionConfig{
		Paradigm: models.ProgressionParadigmLevelBased,
		LevelBased: models.LevelBasedConfig{MinLevel: 5, MaxLevel: 1},
	}
	if err := validateProgressionConfig(cfg); !errors.Is(err, ErrInvalidProgression) {
		t.Fatalf("expected ErrInvalidProgression, got %v", err)
	}
}

func TestValidateProgressionConfig_PointBuyPool(t *testing.T) {
	cfg := &models.ProgressionConfig{
		Paradigm: models.ProgressionParadigmPointBuy,
		PointBuy: models.PointBuyConfig{StartingPool: -1},
	}
	if err := validateProgressionConfig(cfg); !errors.Is(err, ErrInvalidProgression) {
		t.Fatalf("expected ErrInvalidProgression, got %v", err)
	}
}

func TestSanitizeProgressionConfig_ClearsPointBuyForLevelBased(t *testing.T) {
	cfg := &models.ProgressionConfig{
		Paradigm: models.ProgressionParadigmLevelBased,
		PointBuy: models.PointBuyConfig{StartingPool: 10},
	}
	sanitizeProgressionConfig(cfg)
	if cfg.PointBuy.StartingPool != 0 {
		t.Fatalf("expected point_buy cleared")
	}
}
