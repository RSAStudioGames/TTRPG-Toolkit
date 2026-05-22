package services

import (
	"fmt"
	"strings"

	"github.com/gabriel/ttrpg-toolkit/backend/internal/models"
)

func sanitizeSkillConfig(skillType string, cfg *models.SkillConfig) {
	if cfg == nil {
		return
	}
	if !cfg.AllowSpecializations {
		cfg.SpecializationBonus = 0
	} else if cfg.SpecializationBonus == 0 {
		cfg.SpecializationBonus = 2
	}
	switch skillType {
	case models.SkillTypeMultiTier:
		cfg.Min = 0
		cfg.Max = 0
	case models.SkillTypeNumeric:
		cfg.Tiers = nil
	default:
		cfg.Tiers = nil
		cfg.Min = 0
		cfg.Max = 0
	}
}

func validateSkillConfig(skillType string, cfg *models.SkillConfig) error {
	if cfg == nil {
		return fmt.Errorf("%w: config is required", ErrInvalidSkill)
	}
	if !models.IsAllowedSkillType(skillType) {
		return fmt.Errorf("%w: invalid type", ErrInvalidSkill)
	}
	switch skillType {
	case models.SkillTypeMultiTier:
		if len(cfg.Tiers) == 0 {
			return fmt.Errorf("%w: tiers is required for multi_tier skills", ErrInvalidSkill)
		}
		for i, t := range cfg.Tiers {
			if strings.TrimSpace(t.TierName) == "" {
				return fmt.Errorf("%w: tiers[%d].tier_name is required", ErrInvalidSkill, i)
			}
		}
	case models.SkillTypeNumeric:
		if cfg.Min > cfg.Max {
			return fmt.Errorf("%w: min cannot exceed max", ErrInvalidSkill)
		}
	}
	return nil
}

func (svc *MechanicsService) validateLinkedAttribute(systemID string, linkedID *string) error {
	if linkedID == nil || strings.TrimSpace(*linkedID) == "" {
		return nil
	}
	id := strings.TrimSpace(*linkedID)
	if _, err := svc.repo.GetAttribute(systemID, id); err != nil {
		return fmt.Errorf("%w: invalid linked_attribute_id", ErrInvalidSkill)
	}
	return nil
}

func (svc *MechanicsService) ensureUniqueSkillName(systemID, name, excludeID string) error {
	exists, err := svc.repo.SkillNameExists(systemID, name, excludeID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%w: skill name already in use", ErrInvalidSkill)
	}
	return nil
}
