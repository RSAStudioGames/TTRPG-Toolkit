package formula

import (
	"fmt"
	"strings"
	"sync"

	"github.com/alecthomas/participle/v2"
)

var (
	parser     *participle.Parser[TopExpr]
	parserOnce sync.Once
	parserErr  error
)

func getParser() (*participle.Parser[TopExpr], error) {
	parserOnce.Do(func() {
		parser, parserErr = participle.Build[TopExpr](
			participle.Lexer(lexerDef),
			participle.Elide("Whitespace"),
		)
	})
	return parser, parserErr
}

// Parse parses a formula string into an AST. Returns semantic errors as *SemanticError.
func Parse(input string) (*Formula, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("formula is empty")
	}
	p, err := getParser()
	if err != nil {
		return nil, err
	}
	expr, err := p.ParseString("", input)
	if err != nil {
		return nil, err
	}
	f := &Formula{Expr: expr}
	if errs := ValidateSemantics(f); len(errs) > 0 {
		return nil, &SemanticError{Messages: errs}
	}
	if errs := normalizeTop(f.Expr); len(errs) > 0 {
		return nil, &SemanticError{Messages: errs}
	}
	return f, nil
}

func normalizeTop(e *TopExpr) []string {
	if e == nil {
		return nil
	}
	if e.If != nil {
		var errs []string
		for _, o := range []*OrExpr{e.If.Condition, e.If.Then, e.If.Else} {
			errs = append(errs, normalizeOr(o)...)
		}
		return errs
	}
	return normalizeOr(e.Or)
}

func normalizeOr(o *OrExpr) []string {
	if o == nil {
		return nil
	}
	if o.Chain == nil {
		return nil
	}
	var errs []string
	errs = append(errs, normalizeAnd(o.Chain.Left)...)
	for _, op := range o.Chain.Right {
		errs = append(errs, normalizeAnd(op.Right)...)
	}
	return errs
}

func normalizeAnd(a *AndExpr) []string {
	if a == nil {
		return nil
	}
	var errs []string
	errs = append(errs, normalizeNot(a.Left)...)
	for _, op := range a.Right {
		errs = append(errs, normalizeNot(op.Right)...)
	}
	return errs
}

func normalizeNot(n *NotExpr) []string {
	if n == nil || n.Cmp == nil {
		return nil
	}
	var errs []string
	errs = append(errs, normalizeAdd(n.Cmp.Left)...)
	if n.Cmp.Tail != nil {
		errs = append(errs, normalizeAdd(n.Cmp.Tail.Right)...)
	}
	return errs
}

func normalizeAdd(a *AddExpr) []string {
	if a == nil {
		return nil
	}
	var errs []string
	errs = append(errs, normalizeMul(a.Left)...)
	for _, op := range a.Right {
		errs = append(errs, normalizeMul(op.Right)...)
	}
	return errs
}

func normalizeMul(m *MulExpr) []string {
	if m == nil {
		return nil
	}
	var errs []string
	errs = append(errs, normalizeMulOperand(m.Left)...)
	for _, op := range m.Right {
		errs = append(errs, normalizeMulOperand(op.Right)...)
	}
	return errs
}

func normalizeMulOperand(op *MulOperand) []string {
	if op == nil {
		return nil
	}
	if op.Func != nil {
		return nil
	}
	return normalizePrimaryNode(op.Primary)
}

func normalizePrimaryNode(pr *Primary) []string {
	if pr == nil {
		return nil
	}
	if pr.Dice != "" {
		if _, err := ParseDiceToken(pr.Dice); err != nil {
			return []string{err.Error()}
		}
	}
	if pr.Variable != "" {
		if len(pr.Variable) < 3 || pr.Variable[0] != '{' || pr.Variable[len(pr.Variable)-1] != '}' {
			return []string{fmt.Sprintf("invalid variable %q", pr.Variable)}
		}
	}
	if pr.Subexpression != nil {
		return normalizeAdd(pr.Subexpression)
	}
	return nil
}

// SemanticError holds post-parse validation failures.
type SemanticError struct {
	Messages []string
}

func (e *SemanticError) Error() string {
	return strings.Join(e.Messages, "; ")
}
