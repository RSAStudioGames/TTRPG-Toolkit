package formula

import (
	"github.com/alecthomas/participle/v2/lexer"
)

var lexerDef = lexer.MustSimple([]lexer.SimpleRule{
	{"Whitespace", `[ \t\n\r]+`},
	{"Keyword", `\b(if|then|else|and|or|not)\b`},
	{"Dice", `\d+d\d+((kh|kl|dhl)\d+|[rxct]\d+)?`},
	{"Var", `\{[A-Za-z_][A-Za-z0-9_]*\}`},
	{"Float", `\d+\.\d+`},
	{"Int", `\d+`},
	{"Ident", `[A-Za-z_][A-Za-z0-9_]*`},
	{"Operator", `==|!=|<=|>=|[+\-*/%^<>,()]`},
})
