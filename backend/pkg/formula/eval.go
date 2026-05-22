package formula

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/adam-lavrik/go-imath/ix"
)

// RandSource supplies randomness for dice evaluation (testable).
type RandSource interface {
	Intn(n int) int
}

// Eval parses and evaluates a formula AST. Boolean expressions (comparisons,
// and/or/not) return 1.0 for true and 0.0 for false. The "if" construct returns
// the then-branch value when the condition is non-zero, else the else-branch value.
// Dice evaluation requires a non-nil RandSource; pass nil only when the AST has no dice nodes.
func Eval(ast *TopExpr, vars map[string]float64, rng RandSource) (float64, error) {
	if ast == nil {
		return 0, fmt.Errorf("empty expression")
	}
	if astHasDice(ast) && rng == nil {
		return 0, fmt.Errorf("dice evaluation requires a RandSource")
	}
	return evalTop(ast, vars, rng)
}

func astHasDice(ast *TopExpr) bool {
	if ast == nil {
		return false
	}
	if ast.If != nil {
		return exprHasDiceOr(ast.If.Condition) ||
			exprHasDiceOr(ast.If.Then) || exprHasDiceOr(ast.If.Else)
	}
	return exprHasDiceOr(ast.Or)
}

func exprHasDiceOr(o *OrExpr) bool {
	if o == nil || o.Chain == nil {
		return false
	}
	if exprHasDiceAnd(o.Chain.Left) {
		return true
	}
	for _, op := range o.Chain.Right {
		if exprHasDiceAnd(op.Right) {
			return true
		}
	}
	return false
}

func exprHasDiceAnd(a *AndExpr) bool {
	if a == nil {
		return false
	}
	if exprHasDiceNot(a.Left) {
		return true
	}
	for _, op := range a.Right {
		if exprHasDiceNot(op.Right) {
			return true
		}
	}
	return false
}

func exprHasDiceNot(n *NotExpr) bool {
	if n == nil || n.Cmp == nil {
		return false
	}
	if addHasDice(n.Cmp.Left) {
		return true
	}
	if n.Cmp.Tail != nil && addHasDice(n.Cmp.Tail.Right) {
		return true
	}
	return false
}

func addHasDice(a *AddExpr) bool {
	if a == nil {
		return false
	}
	if mulHasDice(a.Left) {
		return true
	}
	for _, op := range a.Right {
		if mulHasDice(op.Right) {
			return true
		}
	}
	return false
}

func mulHasDice(m *MulExpr) bool {
	if m == nil {
		return false
	}
	if mulOperandHasDice(m.Left) {
		return true
	}
	for _, op := range m.Right {
		if mulOperandHasDice(op.Right) {
			return true
		}
	}
	return false
}

func mulOperandHasDice(o *MulOperand) bool {
	if o == nil {
		return false
	}
	if o.Primary != nil {
		return primaryHasDice(o.Primary)
	}
	return funcValueHasDice(o.Func)
}

func primaryHasDice(p *Primary) bool {
	if p == nil {
		return false
	}
	if p.Dice != "" {
		return true
	}
	if p.Subexpression != nil {
		return addHasDice(p.Subexpression)
	}
	return false
}

func funcValueHasDice(f *FuncCall) bool {
	if f == nil {
		return false
	}
	for _, v := range f.Values {
		if v.Dice != "" {
			return true
		}
		if v.Subexpression != nil && addHasDice(v.Subexpression) {
			return true
		}
	}
	return false
}

func evalTop(ast *TopExpr, vars map[string]float64, rng RandSource) (float64, error) {
	if ast.If != nil {
		cond, err := evalOr(ast.If.Condition, vars, rng)
		if err != nil {
			return 0, err
		}
		if cond != 0 {
			return evalOr(ast.If.Then, vars, rng)
		}
		return evalOr(ast.If.Else, vars, rng)
	}
	return evalOr(ast.Or, vars, rng)
}

func evalOr(o *OrExpr, vars map[string]float64, rng RandSource) (float64, error) {
	if o == nil || o.Chain == nil {
		return 0, nil
	}
	if len(o.Chain.Right) == 0 {
		return evalAndValue(o.Chain.Left, vars, rng)
	}
	v, err := evalAnd(o.Chain.Left, vars, rng)
	if err != nil {
		return 0, err
	}
	if v != 0 {
		return 1, nil
	}
	for _, op := range o.Chain.Right {
		v, err = evalAnd(op.Right, vars, rng)
		if err != nil {
			return 0, err
		}
		if v != 0 {
			return 1, nil
		}
	}
	return 0, nil
}

func evalAndValue(a *AndExpr, vars map[string]float64, rng RandSource) (float64, error) {
	if a == nil {
		return 0, nil
	}
	if len(a.Right) == 0 {
		return evalNotValue(a.Left, vars, rng)
	}
	return evalAnd(a, vars, rng)
}

func evalNotValue(n *NotExpr, vars map[string]float64, rng RandSource) (float64, error) {
	if n == nil || n.Cmp == nil {
		return 0, nil
	}
	if n.Prefix != nil || n.Cmp.Tail != nil {
		return evalNot(n, vars, rng)
	}
	return evalAdd(n.Cmp.Left, vars, rng)
}

func evalAnd(a *AndExpr, vars map[string]float64, rng RandSource) (float64, error) {
	if a == nil {
		return 0, nil
	}
	v, err := evalNot(a.Left, vars, rng)
	if err != nil {
		return 0, err
	}
	if v == 0 {
		return 0, nil
	}
	for _, op := range a.Right {
		v, err = evalNot(op.Right, vars, rng)
		if err != nil {
			return 0, err
		}
		if v == 0 {
			return 0, nil
		}
	}
	return 1, nil
}

func evalNot(n *NotExpr, vars map[string]float64, rng RandSource) (float64, error) {
	if n == nil || n.Cmp == nil {
		return 0, nil
	}
	v, err := evalCmp(n.Cmp, vars, rng)
	if err != nil {
		return 0, err
	}
	if n.Prefix != nil {
		if v != 0 {
			return 0, nil
		}
		return 1, nil
	}
	return v, nil
}

func evalCmp(c *CmpExpr, vars map[string]float64, rng RandSource) (float64, error) {
	left, err := evalAdd(c.Left, vars, rng)
	if err != nil {
		return 0, err
	}
	if c.Tail == nil {
		return left, nil
	}
	right, err := evalAdd(c.Tail.Right, vars, rng)
	if err != nil {
		return 0, err
	}
	ok := false
	switch c.Tail.Op {
	case "==":
		ok = left == right
	case "!=":
		ok = left != right
	case "<":
		ok = left < right
	case "<=":
		ok = left <= right
	case ">":
		ok = left > right
	case ">=":
		ok = left >= right
	default:
		return 0, fmt.Errorf("unknown operator %q", c.Tail.Op)
	}
	if ok {
		return 1, nil
	}
	return 0, nil
}

func evalAdd(a *AddExpr, vars map[string]float64, rng RandSource) (float64, error) {
	if a == nil {
		return 0, nil
	}
	sum, err := evalMul(a.Left, vars, rng)
	if err != nil {
		return 0, err
	}
	for _, op := range a.Right {
		v, err := evalMul(op.Right, vars, rng)
		if err != nil {
			return 0, err
		}
		if op.Op == "-" {
			sum -= v
		} else {
			sum += v
		}
	}
	return sum, nil
}

func evalMul(m *MulExpr, vars map[string]float64, rng RandSource) (float64, error) {
	if m == nil {
		return 0, nil
	}
	prod, err := evalMulOperand(m.Left, vars, rng)
	if err != nil {
		return 0, err
	}
	for _, op := range m.Right {
		v, err := evalMulOperand(op.Right, vars, rng)
		if err != nil {
			return 0, err
		}
		switch op.Op {
		case "*":
			prod *= v
		case "/":
			if v == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			prod /= v
		case "%":
			if v == 0 {
				return 0, fmt.Errorf("modulo by zero")
			}
			prod = math.Mod(prod, v)
		case "^":
			prod = math.Pow(prod, v)
		default:
			return 0, fmt.Errorf("unknown operator %q", op.Op)
		}
	}
	return prod, nil
}

func evalMulOperand(o *MulOperand, vars map[string]float64, rng RandSource) (float64, error) {
	if o == nil {
		return 0, nil
	}
	if o.Func != nil {
		return evalFunc(o.Func, vars, rng)
	}
	return evalPrimary(o.Primary, vars, rng)
}

func evalPrimary(p *Primary, vars map[string]float64, rng RandSource) (float64, error) {
	if p == nil {
		return 0, nil
	}
	if p.Number != nil {
		return *p.Number, nil
	}
	if p.Dice != "" {
		d, err := ParseDiceToken(p.Dice)
		if err != nil {
			return 0, err
		}
		return rollDice(d, rng)
	}
	if p.Variable != "" {
		key := strings.TrimSpace(p.Variable)
		if len(key) >= 2 && key[0] == '{' && key[len(key)-1] == '}' {
			key = key[1 : len(key)-1]
		}
		v, ok := vars[key]
		if !ok {
			return 0, fmt.Errorf("unresolved variable %q", key)
		}
		return v, nil
	}
	if p.Subexpression != nil {
		return evalAdd(p.Subexpression, vars, rng)
	}
	return 0, nil
}

func evalFuncValue(v *FuncValue, vars map[string]float64, rng RandSource) (float64, error) {
	if v == nil {
		return 0, nil
	}
	if v.Number != nil {
		return *v.Number, nil
	}
	if v.Dice != "" {
		d, err := ParseDiceToken(v.Dice)
		if err != nil {
			return 0, err
		}
		return rollDice(d, rng)
	}
	if v.Variable != "" {
		key := strings.TrimSpace(v.Variable)
		if len(key) >= 2 && key[0] == '{' && key[len(key)-1] == '}' {
			key = key[1 : len(key)-1]
		}
		val, ok := vars[key]
		if !ok {
			return 0, fmt.Errorf("unresolved variable %q", key)
		}
		return val, nil
	}
	if v.Subexpression != nil {
		return evalAdd(v.Subexpression, vars, rng)
	}
	return 0, nil
}

func evalFunc(fc *FuncCall, vars map[string]float64, rng RandSource) (float64, error) {
	args := make([]float64, 0, len(fc.Values))
	for _, v := range fc.Values {
		n, err := evalFuncValue(v, vars, rng)
		if err != nil {
			return 0, err
		}
		args = append(args, n)
	}
	switch fc.Name {
	case "min":
		return evalMin(args), nil
	case "max":
		return evalMax(args), nil
	case "clamp":
		if len(args) != 3 {
			return 0, fmt.Errorf("clamp expects 3 arguments")
		}
		return evalClamp(args[0], args[1], args[2]), nil
	case "floor":
		if len(args) != 1 {
			return 0, fmt.Errorf("floor expects 1 argument")
		}
		return math.Floor(args[0]), nil
	case "ceil":
		if len(args) != 1 {
			return 0, fmt.Errorf("ceil expects 1 argument")
		}
		return math.Ceil(args[0]), nil
	case "round":
		if len(args) != 1 {
			return 0, fmt.Errorf("round expects 1 argument")
		}
		return math.Round(args[0]), nil
	case "abs":
		if len(args) != 1 {
			return 0, fmt.Errorf("abs expects 1 argument")
		}
		return math.Abs(args[0]), nil
	default:
		return 0, fmt.Errorf("unknown function %q", fc.Name)
	}
}

func toIX(v float64) int {
	return int(math.Round(v))
}

func evalMin(args []float64) float64 {
	if len(args) == 0 {
		return 0
	}
	m := toIX(args[0])
	for _, a := range args[1:] {
		m = ix.Min(m, toIX(a))
	}
	return float64(m)
}

func evalMax(args []float64) float64 {
	if len(args) == 0 {
		return 0
	}
	m := toIX(args[0])
	for _, a := range args[1:] {
		m = ix.Max(m, toIX(a))
	}
	return float64(m)
}

func evalClamp(lo, hi, v float64) float64 {
	loI, hiI, vI := toIX(lo), toIX(hi), toIX(v)
	return float64(ix.Max(loI, ix.Min(hiI, vI)))
}

func rollDice(d *DiceRoll, rng RandSource) (float64, error) {
	if d.Count < 1 || d.Sides < 1 {
		return 0, fmt.Errorf("invalid dice %dd%d", d.Count, d.Sides)
	}
	rolls := make([]int, d.Count)
	sum := 0
	for i := 0; i < d.Count; i++ {
		rolls[i] = rng.Intn(d.Sides) + 1
		sum += rolls[i]
	}
	if d.Mod == nil {
		return float64(sum), nil
	}
	sort.Ints(rolls)
	switch d.Mod.Kind {
	case "kh":
		keep := d.Mod.Arg
		if keep > len(rolls) {
			keep = len(rolls)
		}
		s := 0
		for i := len(rolls) - keep; i < len(rolls); i++ {
			s += rolls[i]
		}
		return float64(s), nil
	case "kl":
		keep := d.Mod.Arg
		if keep > len(rolls) {
			keep = len(rolls)
		}
		s := 0
		for i := 0; i < keep; i++ {
			s += rolls[i]
		}
		return float64(s), nil
	case "dhl":
		if len(rolls) == 0 {
			return 0, nil
		}
		return float64(rolls[len(rolls)-1] - rolls[0]), nil
	default:
		return float64(sum), nil
	}
}
