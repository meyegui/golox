package ast

type Expr interface {
	isExpr()
	Accept(v ExprVisitor)
}
