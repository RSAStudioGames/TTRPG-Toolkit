package services

import "testing"

func TestFormulaServiceValidateFormula(t *testing.T) {
	svc := NewFormulaService()
	valid, errs := svc.ValidateFormula("2d6 + {strength_mod}")
	if !valid || len(errs) > 0 {
		t.Fatalf("expected valid formula, got valid=%v errs=%v", valid, errs)
	}
	valid, errs = svc.ValidateFormula("2d")
	if valid {
		t.Fatal("expected invalid formula")
	}
	if len(errs) == 0 {
		t.Fatal("expected error messages")
	}
}
