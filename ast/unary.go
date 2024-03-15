package ast

import (
	"github.com/meyegui/golox/scanner"
)

type Unary struct {
	Operator scanner.Token
	Right    Expr
}

func (u Unary) isExpr() {}
