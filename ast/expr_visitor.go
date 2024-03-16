package ast

type ExprVisitor interface {
	VisitBinary(b *Binary)
	VisitGrouping(g *Grouping)
	VisitLiteral(l *Literal)
	VisitUnary(u *Unary)
}
