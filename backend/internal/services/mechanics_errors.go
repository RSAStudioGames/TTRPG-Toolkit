package services

import "errors"

// ErrInvalidResolution indicates resolution config failed semantic validation.
var ErrInvalidResolution = errors.New("invalid resolution config")

// InvalidFormulaError is returned when roll_expression fails formula parsing.
type InvalidFormulaError struct {
	Errors []string
}

func (e *InvalidFormulaError) Error() string {
	return "invalid formula"
}

// ErrInvalidAttribute indicates attribute definition failed semantic validation.
var ErrInvalidAttribute = errors.New("invalid attribute")

// ErrInvalidSkill indicates skill definition failed semantic validation.
var ErrInvalidSkill = errors.New("invalid skill")

// ErrInvalidProgression indicates progression config failed semantic validation.
var ErrInvalidProgression = errors.New("invalid progression config")

// ErrInvalidResource indicates resource definition failed semantic validation.
var ErrInvalidResource = errors.New("invalid resource")

// AttributeFormulaError is returned when an attribute formula fails parsing.
type AttributeFormulaError struct {
	Field  string
	Errors []string
}

func (e *AttributeFormulaError) Error() string {
	switch e.Field {
	case "derivation_formula":
		return "invalid derivation formula"
	default:
		return "invalid modifier formula"
	}
}
