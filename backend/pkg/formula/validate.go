package formula

import "fmt"

var allowedFuncs = map[string]struct {
	minArgs int
	maxArgs int // 0 means unlimited
}{
	"max":   {2, 0},
	"min":   {2, 0},
	"floor": {1, 1},
	"ceil":  {1, 1},
	"round": {1, 1},
	"abs":   {1, 1},
	"clamp": {3, 3},
}

// ValidateSemantics checks function names and arity on a parsed formula.
func ValidateSemantics(f *Formula) []string {
	var errs []string
	walkTop(f.Expr, func(fc *FuncCall) {
		spec, ok := allowedFuncs[fc.Name]
		if !ok {
			errs = append(errs, fmt.Sprintf("unknown function %q", fc.Name))
			return
		}
		n := len(fc.Values)
		if n < spec.minArgs {
			errs = append(errs, fmt.Sprintf("function %q expects at least %d argument(s), got %d", fc.Name, spec.minArgs, n))
			return
		}
		if spec.maxArgs > 0 && n > spec.maxArgs {
			errs = append(errs, fmt.Sprintf("function %q expects at most %d argument(s), got %d", fc.Name, spec.maxArgs, n))
		}
	})
	return errs
}

func walkTop(e *TopExpr, onFunc func(*FuncCall)) {
	if e == nil {
		return
	}
	if e.If != nil {
		walkOr(e.If.Condition, onFunc)
		walkOr(e.If.Then, onFunc)
		walkOr(e.If.Else, onFunc)
		return
	}
	walkOr(e.Or, onFunc)
}

func walkOr(o *OrExpr, onFunc func(*FuncCall)) {
	if o == nil {
		return
	}
	if o.Chain == nil {
		return
	}
	walkAnd(o.Chain.Left, onFunc)
	for _, op := range o.Chain.Right {
		walkAnd(op.Right, onFunc)
	}
}

func walkAnd(a *AndExpr, onFunc func(*FuncCall)) {
	if a == nil {
		return
	}
	walkNot(a.Left, onFunc)
	for _, op := range a.Right {
		walkNot(op.Right, onFunc)
	}
}

func walkNot(n *NotExpr, onFunc func(*FuncCall)) {
	if n == nil || n.Cmp == nil {
		return
	}
	walkAdd(n.Cmp.Left, onFunc)
	if n.Cmp.Tail != nil {
		walkAdd(n.Cmp.Tail.Right, onFunc)
	}
}

func walkAdd(a *AddExpr, onFunc func(*FuncCall)) {
	if a == nil {
		return
	}
	walkMul(a.Left, onFunc)
	for _, op := range a.Right {
		walkMul(op.Right, onFunc)
	}
}

func walkMul(m *MulExpr, onFunc func(*FuncCall)) {
	if m == nil {
		return
	}
	walkMulOperand(m.Left, onFunc)
	for _, op := range m.Right {
		walkMulOperand(op.Right, onFunc)
	}
}

func walkMulOperand(op *MulOperand, onFunc func(*FuncCall)) {
	if op == nil {
		return
	}
	if op.Func != nil {
		onFunc(op.Func)
		for _, arg := range op.Func.Values {
			walkFuncValue(arg, onFunc)
		}
		return
	}
	walkPrimary(op.Primary, onFunc)
}

func walkFuncValue(v *FuncValue, onFunc func(*FuncCall)) {
	if v == nil {
		return
	}
	if v.Subexpression != nil {
		walkAdd(v.Subexpression, onFunc)
	}
}

func walkPrimary(pr *Primary, onFunc func(*FuncCall)) {
	if pr == nil {
		return
	}
	if pr.Subexpression != nil {
		walkAdd(pr.Subexpression, onFunc)
	}
}
