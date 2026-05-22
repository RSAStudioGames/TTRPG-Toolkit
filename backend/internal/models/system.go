package models

import (
	"encoding/json"
	"time"
)

const (
	StatusDraft     = "draft"
	StatusPublished = "published"
	StatusLocked    = "locked"
	StatusArchived  = "archived"
)

// GameSystem is the persisted game system entity.
type GameSystem struct {
	ID              string    `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	Slug            string    `json:"slug" db:"slug"`
	Edition         *string   `json:"edition,omitempty" db:"edition"`
	Publisher       *string   `json:"publisher,omitempty" db:"publisher"`
	Description     *string   `json:"description,omitempty" db:"description"`
	LicenseType     *string   `json:"license_type,omitempty" db:"license_type"`
	Version         string    `json:"version" db:"version"`
	Playstyle       *string   `json:"playstyle,omitempty" db:"playstyle"`
	Complexity      *int      `json:"complexity,omitempty" db:"complexity"`
	MeasurementUnit *string   `json:"measurement_unit,omitempty" db:"measurement_unit"`
	CurrencySymbol  *string   `json:"currency_symbol,omitempty" db:"currency_symbol"`
	Status          string    `json:"status" db:"status"`
	IsActive        bool      `json:"is_active" db:"is_active"`
	SystemFamily    *string   `json:"system_family,omitempty" db:"system_family"`
	PlayerCountMin  *int      `json:"player_count_min,omitempty" db:"player_count_min"`
	PlayerCountMax  *int      `json:"player_count_max,omitempty" db:"player_count_max"`
	OfficialLinks   []string  `json:"official_links" db:"-"`
	Tags            []string  `json:"tags" db:"-"`
	CoreRulebooks   []string  `json:"core_rulebooks" db:"-"`
	IconURL         *string   `json:"icon_url,omitempty" db:"icon_url"`
	CoverURL        *string   `json:"cover_url,omitempty" db:"cover_url"`
	ParentSystemID  *string   `json:"parent_system_id,omitempty" db:"parent_system_id"`
	IsCore          bool      `json:"is_core" db:"is_core"`
	IsProtected     bool      `json:"is_protected" db:"is_protected"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`

	OfficialLinksJSON json.RawMessage `json:"-" db:"official_links"`
	TagsJSON          json.RawMessage `json:"-" db:"tags"`
	CoreRulebooksJSON json.RawMessage `json:"-" db:"core_rulebooks"`
}

// CreateSystemRequest is the body for POST /api/systems.
type CreateSystemRequest struct {
	Name            string   `json:"name" validate:"required,min=3,max=120"`
	Slug            string   `json:"slug" validate:"omitempty,max=120"`
	Edition         string   `json:"edition" validate:"omitempty,max=255"`
	Publisher       string   `json:"publisher" validate:"omitempty,max=255"`
	Description     string   `json:"description"`
	LicenseType     string   `json:"license_type" validate:"omitempty,oneof=OGL SRD Proprietary Homebrew Custom"`
	Version         string   `json:"version" validate:"omitempty,max=32"`
	Playstyle       string   `json:"playstyle" validate:"omitempty,oneof=Tactical Narrative Simulationist Gamist Custom"`
	Complexity      *int     `json:"complexity" validate:"omitempty,min=1,max=5"`
	MeasurementUnit string   `json:"measurement_unit" validate:"omitempty,oneof=Imperial Metric"`
	CurrencySymbol  string   `json:"currency_symbol" validate:"omitempty,max=16"`
	IsActive        *bool    `json:"is_active"`
	SystemFamily    string   `json:"system_family" validate:"omitempty,max=255"`
	PlayerCountMin  *int     `json:"player_count_min"`
	PlayerCountMax  *int     `json:"player_count_max"`
	OfficialLinks   []string `json:"official_links" validate:"dive,url"`
	Tags            []string `json:"tags"`
	CoreRulebooks   []string `json:"core_rulebooks"`
	IsCore          *bool    `json:"is_core"`
	ParentSystemID  string   `json:"parent_system_id" validate:"omitempty,uuid"`
}

// UpdateSystemRequest is the body for PUT /api/systems/{id}.
type UpdateSystemRequest struct {
	Name            *string  `json:"name" validate:"omitempty,min=3,max=120"`
	Slug            *string  `json:"slug" validate:"omitempty,max=120"`
	Edition         *string  `json:"edition" validate:"omitempty,max=255"`
	Publisher       *string  `json:"publisher" validate:"omitempty,max=255"`
	Description     *string  `json:"description"`
	LicenseType     *string  `json:"license_type" validate:"omitempty,oneof=OGL SRD Proprietary Homebrew Custom"`
	Version         *string  `json:"version" validate:"omitempty,max=32"`
	Playstyle       *string  `json:"playstyle" validate:"omitempty,oneof=Tactical Narrative Simulationist Gamist Custom"`
	Complexity      *int     `json:"complexity" validate:"omitempty,min=1,max=5"`
	MeasurementUnit *string  `json:"measurement_unit" validate:"omitempty,oneof=Imperial Metric"`
	CurrencySymbol  *string  `json:"currency_symbol" validate:"omitempty,max=16"`
	IsActive        *bool    `json:"is_active"`
	SystemFamily    *string  `json:"system_family" validate:"omitempty,max=255"`
	PlayerCountMin  *int     `json:"player_count_min"`
	PlayerCountMax  *int     `json:"player_count_max"`
	OfficialLinks   []string `json:"official_links" validate:"dive,url"`
	Tags            []string `json:"tags"`
	CoreRulebooks   []string `json:"core_rulebooks"`
	IsCore          *bool    `json:"is_core"`
	ParentSystemID  *string  `json:"parent_system_id" validate:"omitempty,uuid"`
	IsProtected     *bool    `json:"is_protected"`
}

// SystemListResponse is paginated list data.
type SystemListResponse struct {
	Items      []GameSystem `json:"items"`
	Page       int          `json:"page"`
	PerPage    int          `json:"per_page"`
	Total      int          `json:"total"`
	TotalPages int          `json:"total_pages"`
}

// DeletePreview describes cascade impact.
type DeletePreview struct {
	TagCount       int `json:"tag_count"`
	RulebookCount  int `json:"rulebook_count"`
	LinkCount      int `json:"link_count"`
	ChildCount     int `json:"child_count"`
	TotalAssociated int `json:"total_associated"`
}

// SaveTemplateRequest is POST save-template body.
type SaveTemplateRequest struct {
	TemplateName        string `json:"template_name" validate:"required,min=1,max=120"`
	TemplateDescription string `json:"template_description"`
}
