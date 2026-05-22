package services

import (
	"errors"
	"testing"

	"github.com/gabriel/ttrpg-toolkit/backend/internal/models"
)

func TestValidateResourceConfig_InvalidFormat(t *testing.T) {
	cfg := &models.ResourceConfig{CurrentMaxFormat: "invalid"}
	if err := validateResourceConfig(cfg); !errors.Is(err, ErrInvalidResource) {
		t.Fatalf("expected ErrInvalidResource, got %v", err)
	}
}

func TestValidateResourceConfig_StepDieFormat(t *testing.T) {
	cfg := &models.ResourceConfig{
		CurrentMaxFormat: models.ResourceFormatStepDie,
		MaxValFormula:    "10",
	}
	if err := validateResourceConfig(cfg); err != nil {
		t.Fatalf("expected valid step_die format, got %v", err)
	}
}

func TestValidateResourceConfig_RecoveryAmountRequired(t *testing.T) {
	cfg := &models.ResourceConfig{
		CurrentMaxFormat: models.ResourceFormatInteger,
		RecoverySchedules: []models.RecoveryScheduleEntry{
			{Trigger: models.RecoveryTriggerShortRest, Amount: " "},
		},
	}
	if err := validateResourceConfig(cfg); !errors.Is(err, ErrInvalidResource) {
		t.Fatalf("expected ErrInvalidResource, got %v", err)
	}
}
