package formula

// Formula is the root of a parsed expression.
type Formula struct {
	Expr *TopExpr `@@`
}

// TopExpr is a conditional or logical-or expression.
type TopExpr struct {
	Or *OrExpr `@@`
	If *IfExpr `| @@`
}

// IfExpr is a conditional: if cond then a else b.
type IfExpr struct {
	Condition *OrExpr `"if" @@`
	Then      *OrExpr `"then" @@`
	Else      *OrExpr `"else" @@`
}

// OrChain is a left-associative or-chain.
type OrChain struct {
	Left  *AndExpr `@@`
	Right []*OrOp  `@@*`
}

// OrExpr is a logical-or expression.
type OrExpr struct {
	Chain *OrChain `@@`
}

// OrOp is the right-hand side of an or expression.
type OrOp struct {
	Op    string   `"or"`
	Right *AndExpr `@@`
}

// AndExpr is a left-associative and-chain.
type AndExpr struct {
	Left  *NotExpr `@@`
	Right []*AndOp `@@*`
}

// AndOp is the right-hand side of an and expression.
type AndOp struct {
	Op    string   `"and"`
	Right *NotExpr `@@`
}

// NotPrefix is an optional unary not.
type NotPrefix struct {
	Mark string `"not"`
}

// NotExpr is an optional unary not followed by a comparison.
type NotExpr struct {
	Prefix *NotPrefix `[ "not" ]`
	Cmp    *CmpExpr   `@@`
}

// CmpTail is the operator and right side of a comparison.
type CmpTail struct {
	Op    string   `@("==" | "!=" | "<=" | ">=" | "<" | ">")`
	Right *AddExpr `@@`
}

// CmpExpr is an optional comparison between additive expressions.
type CmpExpr struct {
	Left *AddExpr `@@`
	Tail *CmpTail `[ @@ ]`
}

// AddExpr is a left-associative +/- chain.
type AddExpr struct {
	Left  *MulExpr `@@`
	Right []*AddOp `@@*`
}

// AddOp is the right-hand side of an add/subtract expression.
type AddOp struct {
	Op    string   `@("+" | "-")`
	Right *MulExpr `@@`
}

// MulExpr is a left-associative */, %, and ^ chain.
type MulExpr struct {
	Left  *MulOperand `@@`
	Right []*MulOp    `@@*`
}

// MulOperand is a primary value or function call.
type MulOperand struct {
	Primary *Primary   `@@ |`
	Func    *FuncCall  `@@`
}

// MulOp is the right-hand side of a multiply/divide/power expression.
type MulOp struct {
	Op    string      `@("*" | "/" | "%" | "^")`
	Right *MulOperand `@@`
}

// Primary is a literal, dice roll, variable, or parenthesized expr.
type Primary struct {
	Number        *float64 `@(Float|Int)`
	Dice          string   `| @Dice`
	Variable      string   `| @Var`
	Subexpression *AddExpr `| "(" @@ ")"`
}

// FuncValue is a function argument (arithmetic only; no nested calls).
type FuncValue struct {
	Number        *float64 `@(Float|Int)`
	Dice          string   `| @Dice`
	Variable      string   `| @Var`
	Subexpression *AddExpr `| "(" @@ ")"`
}

// DiceRoll is NdS with an optional modifier suffix (filled after lexing).
type DiceRoll struct {
	Count int
	Sides int
	Mod   *DiceMod
}

// DiceMod is a dice modifier (kh, kl, dhl, r, x, t, c) with numeric argument.
type DiceMod struct {
	Kind string
	Arg  int
}

// FuncCall is a named function with zero or more arguments.
type FuncCall struct {
	Name   string       `@Ident`
	Values []*FuncValue `"(" [ @@ ( "," @@ )* ] ")"`
}
