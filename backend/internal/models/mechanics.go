package models

import (
	"encoding/json"
)

// Typed config structs map to JSONB columns (fields added in later modules).

// ResolutionConfig maps to system_mechanics.resolution_config.
type ResolutionConfig struct{}

// ProgressionConfig maps to system_mechanics.progression_config.
type ProgressionConfig struct{}

// ActionEconomyConfig maps to system_mechanics.action_economy_config.
type ActionEconomyConfig struct{}

// AttributeConfig maps to system_attributes.config.
type AttributeConfig struct{}

// SkillConfig maps to system_skills.config.
type SkillConfig struct{}

// ResourceConfig maps to system_resources.config.
type ResourceConfig struct{}

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
	ID         string          `json:"id" db:"id"`
	SystemID   string          `json:"system_id" db:"system_id"`
	GroupName  *string         `json:"group_name,omitempty" db:"group_name"`
	Name       string          `json:"name" db:"name"`
	Type       string          `json:"type" db:"type"`
	ConfigJSON json.RawMessage `json:"-" db:"config"`
	SortOrder  int             `json:"sort_order" db:"sort_order"`
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
	GroupName *string         `json:"group_name"`
	Name      string          `json:"name"`
	Type      string          `json:"type"`
	Config    AttributeConfig `json:"config"`
	SortOrder int             `json:"sort_order"`
}

// UpdateAttributeRequest is PUT /api/systems/{id}/attributes/{attrId} body.
type UpdateAttributeRequest struct {
	GroupName *string          `json:"group_name"`
	Name      *string          `json:"name"`
	Type      *string          `json:"type"`
	Config    *AttributeConfig `json:"config"`
	SortOrder *int             `json:"sort_order"`
}

// AttributeResponse is a single attribute with typed config.
type AttributeResponse struct {
	ID        string          `json:"id"`
	SystemID  string          `json:"system_id"`
	GroupName *string         `json:"group_name,omitempty"`
	Name      string          `json:"name"`
	Type      string          `json:"type"`
	Config    AttributeConfig `json:"config"`
	SortOrder int             `json:"sort_order"`
}

// ListAttributesResponse is GET /api/systems/{id}/attributes data.
type ListAttributesResponse struct {
	Items []AttributeResponse `json:"items"`
}

// CreateSkillRequest is POST /api/systems/{id}/skills body.
type CreateSkillRequest struct {
	Name              string     `json:"name"`
	LinkedAttributeID *string    `json:"linked_attribute_id"`
	Type              string     `json:"type"`
	Category          *string    `json:"category"`
	Config            SkillConfig `json:"config"`
	SortOrder         int        `json:"sort_order"`
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
	Name      string         `json:"name"`
	Type      string         `json:"type"`
	Config    ResourceConfig `json:"config"`
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
