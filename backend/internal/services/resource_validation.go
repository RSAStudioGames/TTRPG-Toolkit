package services

import (
	"fmt"
	"strings"

	"github.com/gabriel/ttrpg-toolkit/backend/internal/models"
)

func sanitizeResourceConfig(cfg *models.ResourceConfig) {
	if cfg == nil {
		return
	}
	if cfg.RecoverySchedules == nil {
		cfg.RecoverySchedules = []models.RecoveryScheduleEntry{}
	}
}

func validateResourceConfig(cfg *models.ResourceConfig) error {
	if cfg == nil {
		return fmt.Errorf("%w: config is required", ErrInvalidResource)
	}
	if !models.IsAllowedResourceFormat(cfg.CurrentMaxFormat) {
		return fmt.Errorf("%w: invalid current_max_format", ErrInvalidResource)
	}
	for i, e := range cfg.RecoverySchedules {
		if !models.IsAllowedRecoveryTrigger(e.Trigger) {
			return fmt.Errorf("%w: invalid recovery_schedules[%d].trigger", ErrInvalidResource, i)
		}
		if strings.TrimSpace(e.Amount) == "" {
			return fmt.Errorf("%w: recovery_schedules[%d].amount is required", ErrInvalidResource, i)
		}
	}
	return nil
}

func (svc *MechanicsService) ensureUniqueResourceName(systemID, name, excludeID string) error {
	exists, err := svc.repo.ResourceNameExists(systemID, name, excludeID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%w: resource name already in use", ErrInvalidResource)
	}
	return nil
}
