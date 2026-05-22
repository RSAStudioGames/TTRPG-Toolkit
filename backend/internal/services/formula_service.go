package services

import (
	"strings"

	formulapkg "github.com/gabriel/ttrpg-toolkit/backend/pkg/formula"
)

// FormulaService validates formula syntax using the shared parser package.
type FormulaService struct{}

// NewFormulaService returns a stateless formula validation service.
func NewFormulaService() *FormulaService {
	return &FormulaService{}
}

// ValidateFormula parses the formula into an AST and returns syntax/semantic errors.
func (s *FormulaService) ValidateFormula(formula string) (bool, []string) {
	_, err := formulapkg.Parse(strings.TrimSpace(formula))
	if err != nil {
		return false, formulapkg.FormatErrors(err)
	}
	return true, nil
}

// Evaluate parses, validates, and evaluates a formula with the given variables.
func (s *FormulaService) Evaluate(
	formula string,
	vars map[string]float64,
	rng formulapkg.RandSource,
) (float64, error) {
	f, err := formulapkg.Parse(strings.TrimSpace(formula))
	if err != nil {
		return 0, err
	}
	return formulapkg.Eval(f.Expr, vars, rng)
}
