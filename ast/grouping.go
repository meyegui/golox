package ast

type Grouping struct {
	Expr Expr
}

func (g Grouping) isExpr() {}
