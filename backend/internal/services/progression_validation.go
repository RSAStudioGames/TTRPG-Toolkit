package services

import (
	"fmt"

	"github.com/gabriel/ttrpg-toolkit/backend/internal/models"
)

func sanitizeProgressionConfig(cfg *models.ProgressionConfig) {
	if cfg == nil {
		return
	}
	switch cfg.Paradigm {
	case models.ProgressionParadigmLevelBased, models.ProgressionParadigmMilestone:
		cfg.PointBuy = models.PointBuyConfig{}
	case models.ProgressionParadigmPointBuy:
		cfg.LevelBased = models.LevelBasedConfig{}
	default:
		cfg.LevelBased = models.LevelBasedConfig{}
		cfg.PointBuy = models.PointBuyConfig{}
	}
}

func validateProgressionConfig(cfg *models.ProgressionConfig) error {
	if cfg == nil {
		return fmt.Errorf("%w: config is required", ErrInvalidProgression)
	}
	if !models.IsAllowedProgressionParadigm(cfg.Paradigm) {
		return fmt.Errorf("%w: invalid paradigm", ErrInvalidProgression)
	}
	switch cfg.Paradigm {
	case models.ProgressionParadigmLevelBased, models.ProgressionParadigmMilestone:
		if cfg.LevelBased.MinLevel > cfg.LevelBased.MaxLevel {
			return fmt.Errorf("%w: min_level cannot exceed max_level", ErrInvalidProgression)
		}
	case models.ProgressionParadigmPointBuy:
		if cfg.PointBuy.StartingPool < 0 {
			return fmt.Errorf("%w: starting_pool cannot be negative", ErrInvalidProgression)
		}
	}
	return nil
}
