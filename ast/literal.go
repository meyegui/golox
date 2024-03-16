package ast

type Literal struct {
	Value any
}

func (l Literal) isExpr() {}

func (l *Literal) Accept(ev ExprVisitor) {
	ev.VisitLiteral(l)
}
