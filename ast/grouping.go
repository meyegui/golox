package ast

type Grouping struct {
	Expr Expr
}

func (g Grouping) isExpr() {}

func (g *Grouping) Accept(ev ExprVisitor) {
	ev.VisitGrouping(g)
}
