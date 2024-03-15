package ast

type Literal struct {
	Value any
}

func (l Literal) isExpr() {}
