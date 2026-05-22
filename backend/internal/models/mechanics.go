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

// ActionCostEntry is one row in the action point cost table.
type ActionCostEntry struct {
	Name string `json:"name"`
	Cost int    `json:"cost"`
}

// ActionComboEntry links named actions into a combo.
type ActionComboEntry struct {
	ComboName      string   `json:"combo_name"`
	ComponentNames []string `json:"component_names"`
}

// PointPoolConfig is the action point pool paradigm.
type PointPoolConfig struct {
	PointsPerPool   int               `json:"points_per_pool"`
	RefreshScope    string            `json:"refresh_scope"`
	ActionCostTable []ActionCostEntry `json:"action_cost_table,omitempty"`
}

// ActionSlotEntry is one action type slot definition.
type ActionSlotEntry struct {
	Name              string             `json:"name"`
	Allowance         int                `json:"allowance"`
	AllowanceScope    string             `json:"allowance_scope,omitempty"`
	CarryOver         string             `json:"carry_over"`
	ConvertTarget     string             `json:"convert_target,omitempty"`
	IsReaction        bool               `json:"is_reaction"`
	ReactionTrigger   string             `json:"reaction_trigger,omitempty"`
	IsFreeAction      bool               `json:"is_free_action"`
	FreeActionLimits  string             `json:"free_action_limits,omitempty"`
	InterruptionRules string             `json:"interruption_rules,omitempty"`
	DelayReadyRules   string             `json:"delay_ready_rules,omitempty"`
	Combos            []ActionComboEntry `json:"combos,omitempty"`
}

// TokenTurnConfig is used when turn_structure is token_based.
type TokenTurnConfig struct {
	TokensPerRound int    `json:"tokens_per_round"`
	RefreshOn      string `json:"refresh_on"`
}

// ActionEconomyConfig maps to system_mechanics.action_economy_config.
type ActionEconomyConfig struct {
	TurnStructure             string            `json:"turn_structure"`
	CustomTurnStructureName   string            `json:"custom_turn_structure_name,omitempty"`
	TokenTurn                 *TokenTurnConfig  `json:"token_turn,omitempty"`
	SystemType                string            `json:"system_type"`
	PointPool                 PointPoolConfig   `json:"point_pool,omitempty"`
	ActionSlots               []ActionSlotEntry `json:"action_slots,omitempty"`
	RoundTimeDefinition       string            `json:"round_time_definition"`
	CustomRoundTimeDefinition string            `json:"custom_round_time_definition,omitempty"`
	CombatTimeTrackingMode    string            `json:"combat_time_tracking_mode"`
	TimeEscalationRules       string            `json:"time_escalation_rules,omitempty"`
	InitiativeSystem          string            `json:"initiative_system"`
	InitiativePersistence     string            `json:"initiative_persistence"`
	InitiativeExpression      string            `json:"initiative_expression,omitempty"`
	StaticInitiativeValue       string            `json:"static_initiative_value,omitempty"`
	InitiativeModifiers         string            `json:"initiative_modifiers,omitempty"`
	TieBreaking                 string            `json:"tie_breaking"`
}

// SaveActionEconomyConfigRequest is PUT /api/systems/{id}/mechanics/action-economy body.
type SaveActionEconomyConfigRequest struct {
	TurnStructure             string            `json:"turn_structure" validate:"required"`
	CustomTurnStructureName   string            `json:"custom_turn_structure_name,omitempty"`
	TokenTurn                 *TokenTurnConfig  `json:"token_turn,omitempty"`
	SystemType                string            `json:"system_type" validate:"required"`
	PointPool                 PointPoolConfig   `json:"point_pool"`
	ActionSlots               []ActionSlotEntry `json:"action_slots"`
	RoundTimeDefinition       string            `json:"round_time_definition" validate:"required"`
	CustomRoundTimeDefinition string            `json:"custom_round_time_definition,omitempty"`
	CombatTimeTrackingMode    string            `json:"combat_time_tracking_mode" validate:"required"`
	TimeEscalationRules       string            `json:"time_escalation_rules,omitempty"`
	InitiativeSystem          string            `json:"initiative_system" validate:"required"`
	InitiativePersistence     string            `json:"initiative_persistence" validate:"required"`
	InitiativeExpression      string            `json:"initiative_expression,omitempty"`
	StaticInitiativeValue       string            `json:"static_initiative_value,omitempty"`
	InitiativeModifiers         string            `json:"initiative_modifiers,omitempty"`
	TieBreaking                 string            `json:"tie_breaking" validate:"required"`
}

// ToActionEconomyConfig converts the request DTO to ActionEconomyConfig for persistence.
func (r SaveActionEconomyConfigRequest) ToActionEconomyConfig() ActionEconomyConfig {
	return ActionEconomyConfig{
		TurnStructure:             r.TurnStructure,
		CustomTurnStructureName:   r.CustomTurnStructureName,
		TokenTurn:                 r.TokenTurn,
		SystemType:                r.SystemType,
		PointPool:                 r.PointPool,
		ActionSlots:               r.ActionSlots,
		RoundTimeDefinition:       r.RoundTimeDefinition,
		CustomRoundTimeDefinition: r.CustomRoundTimeDefinition,
		CombatTimeTrackingMode:    r.CombatTimeTrackingMode,
		TimeEscalationRules:       r.TimeEscalationRules,
		InitiativeSystem:          r.InitiativeSystem,
		InitiativePersistence:     r.InitiativePersistence,
		InitiativeExpression:      r.InitiativeExpression,
		StaticInitiativeValue:       r.StaticInitiativeValue,
		InitiativeModifiers:         r.InitiativeModifiers,
		TieBreaking:                 r.TieBreaking,
	}
}

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

// AttributesConfig maps to system_mechanics.attributes_config.
type AttributesConfig struct {
	EnabledDerived bool `json:"enabled_derived"`
}

// SaveAttributesConfigRequest is PUT /api/systems/{id}/mechanics/attributes-config body.
type SaveAttributesConfigRequest struct {
	EnabledDerived bool `json:"enabled_derived"`
}

// ToAttributesConfig converts the request DTO to AttributesConfig for persistence.
func (r SaveAttributesConfigRequest) ToAttributesConfig() AttributesConfig {
	return AttributesConfig{EnabledDerived: r.EnabledDerived}
}

// SystemMechanics is the 1:1 mechanics row per game system.
type SystemMechanics struct {
	ID                      string          `json:"id" db:"id"`
	SystemID                string          `json:"system_id" db:"system_id"`
	ResolutionConfigJSON    json.RawMessage `json:"-" db:"resolution_config"`
	ProgressionConfigJSON   json.RawMessage `json:"-" db:"progression_config"`
	ActionEconomyConfigJSON json.RawMessage `json:"-" db:"action_economy_config"`
	AttributesConfigJSON    json.RawMessage `json:"-" db:"attributes_config"`
}

// SystemAttributeGroup is a named group for organizing attributes.
type SystemAttributeGroup struct {
	ID        string `json:"id" db:"id"`
	SystemID  string `json:"system_id" db:"system_id"`
	Name      string `json:"name" db:"name"`
	SortOrder int    `json:"sort_order" db:"sort_order"`
}

// SystemAttribute is a definitional attribute for a system.
type SystemAttribute struct {
	ID                string          `json:"id" db:"id"`
	SystemID          string          `json:"system_id" db:"system_id"`
	GroupName         *string         `json:"group_name,omitempty" db:"group_name"`
	AttributeGroupID  *string         `json:"attribute_group_id,omitempty" db:"attribute_group_id"`
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
	AttributesConfig    AttributesConfig    `json:"attributes_config"`
}

// CreateAttributeGroupRequest is POST /api/systems/{id}/attribute-groups body.
type CreateAttributeGroupRequest struct {
	Name      string `json:"name" validate:"required,min=1"`
	SortOrder int    `json:"sort_order"`
}

// UpdateAttributeGroupRequest is PUT /api/systems/{id}/attribute-groups/{groupId} body.
type UpdateAttributeGroupRequest struct {
	Name      *string `json:"name"`
	SortOrder *int    `json:"sort_order"`
}

// AttributeGroupResponse is a single attribute group.
type AttributeGroupResponse struct {
	ID        string `json:"id"`
	SystemID  string `json:"system_id"`
	Name      string `json:"name"`
	SortOrder int    `json:"sort_order"`
}

// ListAttributeGroupsResponse is GET attribute-groups payload.
type ListAttributeGroupsResponse struct {
	Items []AttributeGroupResponse `json:"items"`
}

// CreateAttributeRequest is POST /api/systems/{id}/attributes body.
type CreateAttributeRequest struct {
	GroupName         *string         `json:"group_name"`
	AttributeGroupID  *string         `json:"attribute_group_id"`
	ParentAttributeID *string         `json:"parent_attribute_id"`
	Name              string          `json:"name" validate:"required,min=3"`
	Type              string          `json:"type" validate:"required,oneof=numeric step_die descriptive rank_tier custom"`
	Config            AttributeConfig `json:"config" validate:"required"`
	SortOrder         int             `json:"sort_order"`
}

// UpdateAttributeRequest is PUT /api/systems/{id}/attributes/{attrId} body.
type UpdateAttributeRequest struct {
	GroupName         *string          `json:"group_name"`
	AttributeGroupID  *string          `json:"attribute_group_id"`
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
	AttributeGroupID  *string         `json:"attribute_group_id,omitempty"`
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
