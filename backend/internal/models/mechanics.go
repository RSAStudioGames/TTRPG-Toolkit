package models

import (
	"encoding/json"
)

// Typed config structs map to JSONB columns (fields added in later modules).

// ResolutionConfig maps to system_mechanics.resolution_config.
type ResolutionConfig struct {
	ResolutionType        string                      `json:"resolution_type"`
	RollExpression          string                      `json:"roll_expression"`
	CustomParadigmName      string                      `json:"custom_paradigm_name,omitempty"`
	SuccessDetermination  SuccessDetermination        `json:"success_determination"`
	CriticalMechanics     CriticalMechanics           `json:"critical_mechanics"`
	AdvantageDisadvantage []AdvantageDisadvantageEntry `json:"advantage_disadvantage"`
}

// SuccessDetermination defines how roll outcomes map to success.
type SuccessDetermination struct {
	Method                  string       `json:"method"`
	ThresholdLadder         []LadderTier `json:"threshold_ladder"`
	DefaultTargetVariable   string       `json:"default_target_variable"`
}

// LadderTier is one tier in a success threshold ladder.
type LadderTier struct {
	Label    string `json:"label"`
	Operator string `json:"operator"`
	Value    int    `json:"value"`
}

// CriticalMechanics configures critical success and failure behavior.
type CriticalMechanics struct {
	EnableCritSuccess         bool   `json:"enable_crit_success"`
	CritSuccessTrigger        string `json:"crit_success_trigger"`
	CritSuccessExceedAmount   int    `json:"crit_success_exceed_amount"`
	EnableCritFailure         bool   `json:"enable_crit_failure"`
	CritFailureTrigger        string `json:"crit_failure_trigger"`
}

// AdvantageDisadvantageEntry defines a named advantage/disadvantage mechanic.
type AdvantageDisadvantageEntry struct {
	Name         string `json:"name"`
	MechanicType string `json:"mechanic_type"`
}

// SaveResolutionConfigRequest is PUT /api/systems/{id}/mechanics/resolution body.
type SaveResolutionConfigRequest struct {
	ResolutionType        string                      `json:"resolution_type" validate:"required"`
	RollExpression          string                      `json:"roll_expression" validate:"required"`
	CustomParadigmName      string                      `json:"custom_paradigm_name,omitempty"`
	SuccessDetermination  SuccessDetermination        `json:"success_determination" validate:"required"`
	CriticalMechanics     CriticalMechanics           `json:"critical_mechanics" validate:"required"`
	AdvantageDisadvantage []AdvantageDisadvantageEntry `json:"advantage_disadvantage" validate:"required"`
}

// ToResolutionConfig converts the request DTO to a ResolutionConfig for persistence.
func (r SaveResolutionConfigRequest) ToResolutionConfig() ResolutionConfig {
	return ResolutionConfig{
		ResolutionType:        r.ResolutionType,
		RollExpression:          r.RollExpression,
		CustomParadigmName:      r.CustomParadigmName,
		SuccessDetermination:  r.SuccessDetermination,
		CriticalMechanics:     r.CriticalMechanics,
		AdvantageDisadvantage: r.AdvantageDisadvantage,
	}
}

// XPTableEntry is one level/XP pair for level-based progression.
type XPTableEntry struct {
	Level      int `json:"level"`
	XPRequired int `json:"xp_required"`
}

// CostTableEntry is one rating/cost pair for point-buy progression.
type CostTableEntry struct {
	Rating int `json:"rating"`
	Cost   int `json:"cost"`
}

// LevelBasedConfig is configuration for level-based and milestone paradigms.
type LevelBasedConfig struct {
	MinLevel       int            `json:"min_level"`
	MaxLevel       int            `json:"max_level"`
	XPTable        []XPTableEntry `json:"xp_table,omitempty"`
	AllowMilestone bool           `json:"allow_milestone"`
}

// PointBuyConfig is configuration for point-buy progression.
type PointBuyConfig struct {
	StartingPool int              `json:"starting_pool"`
	CostTable    []CostTableEntry `json:"cost_table,omitempty"`
}

// ProgressionConfig maps to system_mechanics.progression_config.
type ProgressionConfig struct {
	Paradigm   string           `json:"paradigm"`
	LevelBased LevelBasedConfig `json:"level_based,omitempty"`
	PointBuy   PointBuyConfig   `json:"point_buy,omitempty"`
	GMApproval bool             `json:"gm_approval"`
	AllowUndo  bool             `json:"allow_undo"`
}

// SaveProgressionConfigRequest is PUT /api/systems/{id}/mechanics/progression body.
type SaveProgressionConfigRequest struct {
	Paradigm   string           `json:"paradigm" validate:"required"`
	LevelBased LevelBasedConfig `json:"level_based"`
	PointBuy   PointBuyConfig   `json:"point_buy"`
	GMApproval bool             `json:"gm_approval"`
	AllowUndo  bool             `json:"allow_undo"`
}

// ToProgressionConfig converts the request DTO to ProgressionConfig for persistence.
func (r SaveProgressionConfigRequest) ToProgressionConfig() ProgressionConfig {
	return ProgressionConfig{
		Paradigm:   r.Paradigm,
		LevelBased: r.LevelBased,
		PointBuy:   r.PointBuy,
		GMApproval: r.GMApproval,
		AllowUndo:  r.AllowUndo,
	}
}

// ActionEconomyConfig maps to system_mechanics.action_economy_config.
type ActionEconomyConfig struct{}

// DescriptiveMapEntry is one label/value pair for descriptive attributes.
type DescriptiveMapEntry struct {
	Label string `json:"label"`
	Value int    `json:"value"`
}

// RankMapEntry is one rank name with numeric backing for rank_tier attributes.
type RankMapEntry struct {
	RankName       string `json:"rank_name"`
	NumericBacking int    `json:"numeric_backing"`
}

// AttributeConfig maps to system_attributes.config.
type AttributeConfig struct {
	Min                 int                   `json:"min,omitempty"`
	Max                 int                   `json:"max,omitempty"`
	NumericFormat       string                `json:"numeric_format,omitempty"`
	StepDice            []string              `json:"step_dice,omitempty"`
	DescriptiveMap      []DescriptiveMapEntry `json:"descriptive_map,omitempty"`
	RankMap             []RankMapEntry        `json:"rank_map,omitempty"`
	ModifierFormula     string                `json:"modifier_formula,omitempty"`
	ModifierDisplay     string                `json:"modifier_display,omitempty"`
	IsDerived           bool                  `json:"is_derived"`
	DerivationFormula   string                `json:"derivation_formula,omitempty"`
	CachingRule         string                `json:"caching_rule,omitempty"`
	RecalculateTriggers []string              `json:"recalculate_triggers,omitempty"`
}

// SkillTierEntry is one tier for multi_tier skills.
type SkillTierEntry struct {
	TierName       string `json:"tier_name"`
	NumericBacking int    `json:"numeric_backing"`
}

// SkillConfig maps to system_skills.config.
type SkillConfig struct {
	Tiers                 []SkillTierEntry `json:"tiers,omitempty"`
	Min                   int              `json:"min,omitempty"`
	Max                   int              `json:"max,omitempty"`
	AllowSpecializations  bool             `json:"allow_specializations"`
	SpecializationBonus   int              `json:"specialization_bonus,omitempty"`
}

// RecoveryScheduleEntry defines when and how a resource recovers.
type RecoveryScheduleEntry struct {
	Trigger    string `json:"trigger"`
	Amount     string `json:"amount"`
	Conditions string `json:"conditions,omitempty"`
}

// ResourceConfig maps to system_resources.config.
type ResourceConfig struct {
	CurrentMaxFormat  string                  `json:"current_max_format"`
	MinVal            int                     `json:"min_val"`
	MaxValFormula     string                  `json:"max_val_formula"`
	RecoverySchedules []RecoveryScheduleEntry `json:"recovery_schedules,omitempty"`
}

// SystemMechanics is the 1:1 mechanics row per game system.
type SystemMechanics struct {
	ID                    string          `json:"id" db:"id"`
	SystemID              string          `json:"system_id" db:"system_id"`
	ResolutionConfigJSON  json.RawMessage `json:"-" db:"resolution_config"`
	ProgressionConfigJSON json.RawMessage `json:"-" db:"progression_config"`
	ActionEconomyConfigJSON json.RawMessage `json:"-" db:"action_economy_config"`
}

// SystemAttribute is a definitional attribute for a system.
type SystemAttribute struct {
	ID                string          `json:"id" db:"id"`
	SystemID          string          `json:"system_id" db:"system_id"`
	GroupName         *string         `json:"group_name,omitempty" db:"group_name"`
	Name              string          `json:"name" db:"name"`
	Type              string          `json:"type" db:"type"`
	ConfigJSON        json.RawMessage `json:"-" db:"config"`
	SortOrder         int             `json:"sort_order" db:"sort_order"`
	ParentAttributeID *string         `json:"parent_attribute_id,omitempty" db:"parent_attribute_id"`
}

// SystemSkill is a skill definition linked optionally to an attribute.
type SystemSkill struct {
	ID                string          `json:"id" db:"id"`
	SystemID          string          `json:"system_id" db:"system_id"`
	Name              string          `json:"name" db:"name"`
	LinkedAttributeID *string         `json:"linked_attribute_id,omitempty" db:"linked_attribute_id"`
	Type              string          `json:"type" db:"type"`
	Category          *string         `json:"category,omitempty" db:"category"`
	ConfigJSON        json.RawMessage `json:"-" db:"config"`
	SortOrder         int             `json:"sort_order" db:"sort_order"`
}

// SystemResource is a resource pool/track/counter definition.
type SystemResource struct {
	ID         string          `json:"id" db:"id"`
	SystemID   string          `json:"system_id" db:"system_id"`
	Name       string          `json:"name" db:"name"`
	Type       string          `json:"type" db:"type"`
	ConfigJSON json.RawMessage `json:"-" db:"config"`
	SortOrder  int             `json:"sort_order" db:"sort_order"`
}

// EmptyJSONObject returns {} for repository inserts and normalization.
func EmptyJSONObject() json.RawMessage {
	return json.RawMessage("{}")
}

// MarshalConfig serializes a typed config to JSON for JSONB storage.
func MarshalConfig[T any](v T) (json.RawMessage, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	if len(b) == 0 {
		return EmptyJSONObject(), nil
	}
	return json.RawMessage(b), nil
}

// UnmarshalConfig deserializes JSONB into a typed config struct.
func UnmarshalConfig[T any](raw json.RawMessage) (T, error) {
	var out T
	if len(raw) == 0 {
		return out, nil
	}
	err := json.Unmarshal(raw, &out)
	return out, err
}

// UpsertMechanicsRequest is PUT /api/systems/{id}/mechanics body.
type UpsertMechanicsRequest struct {
	ResolutionConfig  ResolutionConfig  `json:"resolution_config"`
	ProgressionConfig   ProgressionConfig   `json:"progression_config"`
	ActionEconomyConfig ActionEconomyConfig `json:"action_economy_config"`
}

// MechanicsResponse is GET mechanics payload with typed configs.
type MechanicsResponse struct {
	ID                  string              `json:"id"`
	SystemID            string              `json:"system_id"`
	ResolutionConfig    ResolutionConfig    `json:"resolution_config"`
	ProgressionConfig   ProgressionConfig   `json:"progression_config"`
	ActionEconomyConfig ActionEconomyConfig `json:"action_economy_config"`
}

// CreateAttributeRequest is POST /api/systems/{id}/attributes body.
type CreateAttributeRequest struct {
	GroupName         *string         `json:"group_name"`
	ParentAttributeID *string         `json:"parent_attribute_id"`
	Name              string          `json:"name" validate:"required,min=3"`
	Type              string          `json:"type" validate:"required,oneof=numeric step_die descriptive rank_tier custom"`
	Config            AttributeConfig `json:"config" validate:"required"`
	SortOrder         int             `json:"sort_order"`
}

// UpdateAttributeRequest is PUT /api/systems/{id}/attributes/{attrId} body.
type UpdateAttributeRequest struct {
	GroupName         *string          `json:"group_name"`
	ParentAttributeID *string          `json:"parent_attribute_id"`
	Name              *string          `json:"name"`
	Type              *string          `json:"type"`
	Config            *AttributeConfig `json:"config"`
	SortOrder         *int             `json:"sort_order"`
}

// AttributeResponse is a single attribute with typed config.
type AttributeResponse struct {
	ID                string          `json:"id"`
	SystemID          string          `json:"system_id"`
	GroupName         *string         `json:"group_name,omitempty"`
	ParentAttributeID *string         `json:"parent_attribute_id,omitempty"`
	Name              string          `json:"name"`
	Type              string          `json:"type"`
	Config            AttributeConfig `json:"config"`
	SortOrder         int             `json:"sort_order"`
}

// ListAttributesResponse is GET /api/systems/{id}/attributes data.
type ListAttributesResponse struct {
	Items []AttributeResponse `json:"items"`
}

// CreateSkillRequest is POST /api/systems/{id}/skills body.
type CreateSkillRequest struct {
	Name              string      `json:"name" validate:"required"`
	LinkedAttributeID *string     `json:"linked_attribute_id"`
	Type              string      `json:"type" validate:"required,oneof=binary multi_tier numeric step_die rank_name"`
	Category          *string     `json:"category"`
	Config            SkillConfig `json:"config" validate:"required"`
	SortOrder         int         `json:"sort_order"`
}

// UpdateSkillRequest is PUT /api/systems/{id}/skills/{skillId} body.
type UpdateSkillRequest struct {
	Name              *string     `json:"name"`
	LinkedAttributeID *string     `json:"linked_attribute_id"`
	Type              *string     `json:"type"`
	Category          *string     `json:"category"`
	Config            *SkillConfig `json:"config"`
	SortOrder         *int        `json:"sort_order"`
}

// SkillResponse is a single skill with typed config.
type SkillResponse struct {
	ID                string      `json:"id"`
	SystemID          string      `json:"system_id"`
	Name              string      `json:"name"`
	LinkedAttributeID *string     `json:"linked_attribute_id,omitempty"`
	Type              string      `json:"type"`
	Category          *string     `json:"category,omitempty"`
	Config            SkillConfig `json:"config"`
	SortOrder         int         `json:"sort_order"`
}

// ListSkillsResponse is GET /api/systems/{id}/skills data.
type ListSkillsResponse struct {
	Items []SkillResponse `json:"items"`
}

// CreateResourceRequest is POST /api/systems/{id}/resources body.
type CreateResourceRequest struct {
	Name      string         `json:"name" validate:"required"`
	Type      string         `json:"type" validate:"required,oneof=pool slot_track counter currency custom"`
	Config    ResourceConfig `json:"config" validate:"required"`
	SortOrder int            `json:"sort_order"`
}

// UpdateResourceRequest is PUT /api/systems/{id}/resources/{resourceId} body.
type UpdateResourceRequest struct {
	Name      *string         `json:"name"`
	Type      *string         `json:"type"`
	Config    *ResourceConfig `json:"config"`
	SortOrder *int            `json:"sort_order"`
}

// ResourceResponse is a single resource with typed config.
type ResourceResponse struct {
	ID        string         `json:"id"`
	SystemID  string         `json:"system_id"`
	Name      string         `json:"name"`
	Type      string         `json:"type"`
	Config    ResourceConfig `json:"config"`
	SortOrder int            `json:"sort_order"`
}

// ListResourcesResponse is GET /api/systems/{id}/resources data.
type ListResourcesResponse struct {
	Items []ResourceResponse `json:"items"`
}

// ValidateFormulaRequest is POST /api/systems/{id}/validate-formula body.
type ValidateFormulaRequest struct {
	Formula string `json:"formula" validate:"required"`
}

// ValidateFormulaResponse is POST /api/systems/{id}/validate-formula success data.
type ValidateFormulaResponse struct {
	Valid bool `json:"valid"`
}
