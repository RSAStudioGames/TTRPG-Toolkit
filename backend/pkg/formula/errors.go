package formula

import (
	"errors"
	"fmt"

	"github.com/alecthomas/participle/v2"
)

// FormatErrors turns parse or semantic errors into human-readable messages.
func FormatErrors(err error) []string {
	if err == nil {
		return nil
	}
	var sem *SemanticError
	if errors.As(err, &sem) {
		return sem.Messages
	}
	var pe *participle.ParseError
	if errors.As(err, &pe) {
		return []string{formatParticipleError(pe)}
	}
	return []string{err.Error()}
}

func formatParticipleError(pe *participle.ParseError) string {
	pos := pe.Position()
	if pos.Line != 0 || pos.Column != 0 {
		return fmt.Sprintf("line %d, column %d: %s", pos.Line, pos.Column, pe.Message())
	}
	return pe.Message()
}
