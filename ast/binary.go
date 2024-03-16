package ast

import (
	"github.com/meyegui/golox/scanner"
)

type Binary struct {
	Left     Expr
	Operator *scanner.Token
	Right    Expr
}

func (b Binary) isExpr() {}

func (b *Binary) Accept(ev ExprVisitor) {
	ev.VisitBinary(b)
}
