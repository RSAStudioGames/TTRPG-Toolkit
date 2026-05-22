package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gabriel/ttrpg-toolkit/backend/internal/models"
	"github.com/gabriel/ttrpg-toolkit/backend/internal/repository"
)

const maxAttributeNestDepth = 32

func sanitizeAttributeConfig(cfg *models.AttributeConfig) {
	if cfg == nil {
		return
	}
	if !cfg.IsDerived {
		cfg.DerivationFormula = ""
		cfg.CachingRule = ""
		cfg.RecalculateTriggers = nil
	}
}

func validateAttributeConfig(attrType string, cfg *models.AttributeConfig) error {
	if cfg == nil {
		return fmt.Errorf("%w: config is required", ErrInvalidAttribute)
	}
	if !models.IsAllowedAttributeType(attrType) {
		return fmt.Errorf("%w: invalid type", ErrInvalidAttribute)
	}
	if cfg.ModifierDisplay != "" && !models.IsAllowedModifierDisplay(cfg.ModifierDisplay) {
		return fmt.Errorf("%w: invalid modifier_display", ErrInvalidAttribute)
	}

	switch attrType {
	case models.AttributeTypeNumeric:
		if !models.IsAllowedNumericFormat(cfg.NumericFormat) {
			return fmt.Errorf("%w: numeric_format is required for numeric attributes", ErrInvalidAttribute)
		}
		if cfg.Min > cfg.Max {
			return fmt.Errorf("%w: min cannot exceed max", ErrInvalidAttribute)
		}
	case models.AttributeTypeStepDie:
		if len(cfg.StepDice) == 0 {
			return fmt.Errorf("%w: step_dice is required for step_die attributes", ErrInvalidAttribute)
		}
	case models.AttributeTypeDescriptive:
		if len(cfg.DescriptiveMap) == 0 {
			return fmt.Errorf("%w: descriptive_map is required for descriptive attributes", ErrInvalidAttribute)
		}
		for i, e := range cfg.DescriptiveMap {
			if strings.TrimSpace(e.Label) == "" {
				return fmt.Errorf("%w: descriptive_map[%d].label is required", ErrInvalidAttribute, i)
			}
		}
	case models.AttributeTypeRankTier:
		if len(cfg.RankMap) == 0 {
			return fmt.Errorf("%w: rank_map is required for rank_tier attributes", ErrInvalidAttribute)
		}
		for i, e := range cfg.RankMap {
			if strings.TrimSpace(e.RankName) == "" {
				return fmt.Errorf("%w: rank_map[%d].rank_name is required", ErrInvalidAttribute, i)
			}
		}
	}

	if cfg.IsDerived {
		if strings.TrimSpace(cfg.DerivationFormula) == "" {
			return fmt.Errorf("%w: derivation_formula is required when is_derived is true", ErrInvalidAttribute)
		}
		if !models.IsAllowedCachingRule(cfg.CachingRule) {
			return fmt.Errorf("%w: caching_rule is required when is_derived is true", ErrInvalidAttribute)
		}
		if cfg.CachingRule == models.CachingRuleOnTrigger && len(cfg.RecalculateTriggers) == 0 {
			return fmt.Errorf("%w: recalculate_triggers is required when caching_rule is on_trigger", ErrInvalidAttribute)
		}
	}

	return nil
}

func (svc *MechanicsService) validateAttributeFormulas(cfg *models.AttributeConfig) error {
	if svc.formula == nil || cfg == nil {
		return nil
	}
	if f := strings.TrimSpace(cfg.ModifierFormula); f != "" {
		valid, errs := svc.formula.ValidateFormula(f)
		if !valid {
			return &AttributeFormulaError{Field: "modifier_formula", Errors: errs}
		}
	}
	if f := strings.TrimSpace(cfg.DerivationFormula); f != "" {
		valid, errs := svc.formula.ValidateFormula(f)
		if !valid {
			return &AttributeFormulaError{Field: "derivation_formula", Errors: errs}
		}
	}
	return nil
}

func (svc *MechanicsService) validateParentAttribute(systemID, attrID string, parentID *string) error {
	if parentID == nil || strings.TrimSpace(*parentID) == "" {
		return nil
	}
	pid := strings.TrimSpace(*parentID)
	if attrID != "" && pid == attrID {
		return fmt.Errorf("%w: attribute cannot be its own parent", ErrInvalidAttribute)
	}
	if _, err := svc.repo.GetAttribute(systemID, pid); err != nil {
		return err
	}
	if attrID == "" {
		return nil
	}
	visited := map[string]struct{}{attrID: {}}
	current := pid
	for depth := 0; depth < maxAttributeNestDepth; depth++ {
		if _, ok := visited[current]; ok {
			return fmt.Errorf("%w: circular attribute nesting", ErrInvalidAttribute)
		}
		visited[current] = struct{}{}
		row, err := svc.repo.GetAttribute(systemID, current)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return fmt.Errorf("%w: invalid parent_attribute_id", ErrInvalidAttribute)
			}
			return err
		}
		if row.ParentAttributeID == nil || strings.TrimSpace(*row.ParentAttributeID) == "" {
			return nil
		}
		current = strings.TrimSpace(*row.ParentAttributeID)
	}
	return fmt.Errorf("%w: attribute nesting exceeds maximum depth", ErrInvalidAttribute)
}

func (svc *MechanicsService) ensureSystemExists(systemID string) error {
	if svc.systems == nil {
		return nil
	}
	_, err := svc.systems.GetByID(systemID)
	return err
}

func (svc *MechanicsService) ensureUniqueAttributeName(systemID, name, excludeID string) error {
	exists, err := svc.repo.AttributeNameExists(systemID, name, excludeID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%w: attribute name already in use", ErrInvalidAttribute)
	}
	return nil
}
